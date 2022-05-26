from queue import Queue
from typing import List, Callable, Any
import threading, random as rand
from PIL import Image

from .Channel import DrawChan, EventChan
from .Event import InputType, MouseEvent, KeyEvent, BroadcastEvent, PathDropEvent

# TODO: how to poll for events in a non-blocking way, such that the env can write custom callbacks for when an event is received?
class Env:
    __events: EventChan
    __draw: DrawChan
    __ready: threading.Condition
    window: bool
    id: int
    name: str

    def __init__(self, main=False, id=0, threadName=''):
        self.window = main
        self.__events = EventChan(Queue(), Queue())
        self.__draw = DrawChan()
        self.setId(id) 
        if main == True:
            self.createLock()
        self.name = threadName

    def createLock(self):
        assert self.window == True, 'must be Main Env to create lock'
        self.__ready = threading.Condition()
    
    def getLock(self):
        return self.__ready
    
    def addCond(self, lock):
        self.__ready = lock
    
    def setId(self, newId):
        self.id = newId
    
    def eventChan(self) -> EventChan:
        return self.__events
    
    def drawChan(self) -> DrawChan:
        return self.__draw

    def setEventChan(self, eventsIn, eventsOut):
        '''
        Link the Env's event channels to the correct channel objects
        :param eventsIn: Queue to receive events
        :param eventsOut: Queue to send events
        '''
        self.__events = EventChan(eventsIn, eventsOut)
    
    def setDrawChan(self, drawChan):
        '''
        Link the Env's draw channel to the Window (main) components draw channel, where it sends commands.
        :param drawChan: DrawChan object where drawing functions are sent
        '''
        self.__draw = drawChan
    
    def onMouseClick(self, event):
        '''
        Callback function that responds to a mouse button being pressed or released.
        :param event: a MouseEvent object that represents the mouse button and action that occurred.
        '''
        pass

    def onKeyPress(self, keyEvent: KeyEvent):
        '''
        Callback function that responds to a key being pressed or released.
        :param event: a KeyEvent object that represents the key pressed and action that occurred.
        '''
        pass
    
    def onBroadcast(self, event):
        '''
        Callback function that responds to a Broadcast event, sent by a separate component and propagated by the Mux.
        :param event: an arbitrary Broadcast object, containing an event type and a message.
        '''
        pass

    def onPathDrop(self, event):
        '''
        Callback function that responds to an item's path being dragged and dropped over the window
        '''
        pass

    def drawImg(self, drawCommand: Callable[...,Image.Image]):
        '''
        Create drawing function that gets sent to the Window (main Env) for rendering.
        :param drawCommand: the drawing function to send. The function must match the function signature.
        '''
        self.__draw.send(drawCommand)

    def sendEvent(self, event: Any):
        '''
        Broadcast event to all other components. 
        :param event: a BroadcastEvent object that contains the broadcast type and object to send, if necessary.
        '''
        self.__events.send(event)

    def init(self):
        '''
        A function to contain any component setup before the event loop starts in its own thread.
        '''
        pass
    
    def run(self, name='') -> None:
        '''
        should start running it's logic and callback listeners on a thread separate from the main Python interpreter thread
        '''
        def startThread():
            with self.__ready:
                self.__ready.wait()
            self.init()
            while True:
                event = self.eventChan().receive()
                if type(event) == BroadcastEvent:
                  self.onBroadcast(event)
                if type(event) == MouseEvent:
                    self.onMouseClick(event)
                elif type(event) == KeyEvent:
                    self.onKeyPress(event)
                elif type(event) == PathDropEvent:
                    self.onPathDrop(event)
        threading.Thread(target=startThread, name=f'{name}', daemon=True).start() 


class Mux():
    '''
    The Mux class acts as a multiplexer for the main environment. It can create new sub-environments that communicate with the Mux via Channel objects.
    The Mux receives events from the main Env, which it passes to the multiplex Envs. It receives drawing commands from the sub-Envs, which it passes to the main Env.
    Will throw an error if it is initialized without a main Env.
    '''

    main: Env
    envs: List[Env]
    muxEvents: List[EventChan]
    mainEvents: EventChan
    drawEvents: DrawChan

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'missing Main Env: Mux must be created from an existing Env'
    
        self.main = mainEnv
        self.envs = []
        self.mainEvents = EventChan(mainEnv.eventChan().getEventsOut(), mainEnv.eventChan().getEventsIn())
        self.muxEvents = []
  
    def addEnv(self, newEnv: Env) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        newEnv.setDrawChan(self.main.drawChan())
        newEnv.addCond(self.main.getLock())
        envChan = EventChan(newEnv.eventChan().getEventsOut(), newEnv.eventChan().getEventsIn())
        self.muxEvents.append(envChan)
        self.envs.append(newEnv)
        
        return newEnv

    def forwardMainEvents(self):
        '''
        Create thread for listening for main events and forwarding them to the other components.
        '''
        def propagate():
            while True:
                event = self.mainEvents.receive()
                for env in self.muxEvents:
                    env.send(event)
        threading.Thread(target=propagate, daemon=True).start()
    
    def beginBroadcast(self):
        '''
        Create thread for listening for events coming from components and 
        broadcasting them to all other components (including the main Window).
        '''
        def brdcstEvents():
            ready = self.main.getLock()
            with ready:
                ready.wait() 
            while True:
                for env in range(len(self.envs)):
                    event = self.muxEvents[env].receiveTimeout(timeout=0.1)
                    self.mainEvents.send(event)
                    for e in range(len(self.envs)):
                        if self.envs[e].id != self.envs[env].id and event is not None:
                            self.muxEvents[e].send(event)
        
        threading.Thread(target=(brdcstEvents), daemon=True).start() 

    def run(self) -> None:
        '''
        Begin drawing and event threads
        '''
        for env in self.envs:
            env.run(env.name)
        self.forwardMainEvents()
        self.beginBroadcast()