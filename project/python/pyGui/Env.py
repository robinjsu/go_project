from queue import Queue
from typing import List, Callable
import threading, random as rand
from PIL import Image

from .Channel import DrawChan, EventChan
from .Event import MouseEvent, KeyEvent, Broadcast
# from .Window import Window

# TODO: how to poll for events in a non-blocking way, such that the env can write custom callbacks for when an event is received?
class Env:
    events: EventChan
    draw: DrawChan
    window: bool
    _ready: threading.Condition
    id: int

    def __init__(self, main=False, id=0):
        self.window = main
        self.events = EventChan(Queue(), Queue())
        self.draw = DrawChan()
        self.setId(id) 
        if main == True:
            self.createLock()

    
    def createLock(self):
        assert self.window == True, 'must be Main Env to create lock'
        self._ready = threading.Condition()
    
    def getLock(self):
        return self._ready
    
    def addCond(self, lock):
        self._ready = lock
    
    def setId(self, newId):
        self.id = newId
    
    def eventChan(self) -> EventChan:
        return self.events
    
    def drawChan(self) -> DrawChan:
        return self.draw

    def setEventChan(self, eventsIn, eventsOut):
        self.events = EventChan(eventsIn, eventsOut)
    
    def setDrawChan(self, drawChan):
        self.draw = drawChan
    
    def onStartUp(self):
        pass
    
    def onMouseClick(self, action):
        pass

    def onKeyPress(self, keyPressed):
        pass
    
    def onBroadcast(self, event):
        pass

    def drawImg(self, drawCommand: Callable[...,Image.Image]):
        self.draw.send(drawCommand)

    def init(self):
        pass
    
    def run(self) -> None:
        '''
        should start running it's logic and callback listeners on a thread separate from the main Python interpreter thread
        '''
        def startThread():
            with self._ready:
                self._ready.wait()
            self.init()
            while True:
                event = self.eventChan().receive()
                print(event)
                if type(event) == Broadcast:
                  self.onBroadcast()
                if type(event) == MouseEvent:
                    self.onMouseClick(event.action)
                elif type(event) == KeyEvent:
                    self.onKeyPress(event.key)
        threading.Thread(target=startThread, name="DisplayThread", daemon=True).start() 


# mux should receive events from the main env and pass on to each of the sub envs
# mux will receive events from the sub envs and pass on to the main env
# the two situations above can be handled by separate EventChan objects
# TODO: NEW IDEA - mux is not an env, but simply acts as the component that connects the window to its components
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
    # listenMainEvents: threading.Thread
    # broadcastEvents: threading.Thread

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'missing Main Env: Mux must be created from an existing Env'
    
        self.main = mainEnv
        self.envs = []
        self.mainEvents = EventChan(mainEnv.eventChan().getEventsOut(), mainEnv.eventChan().getEventsIn())
        self.muxEvents = []
        # self.listenMainEvents = self.forwardMainEvents()
        # self.broadcastEvents = self.broadcast()
  
    def addEnv(self, newEnv: Env) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        # newEnv.setEventChan(self.main.eventChan().getEventsOut(),self.main.eventChan().getEventsIn())
        # self.muxEvents.append(newEnv.eventChan())
        newEnv.setDrawChan(self.main.drawChan())
        newEnv.addCond(self.main.getLock())
        envChan = EventChan(newEnv.eventChan().getEventsOut(), newEnv.eventChan().getEventsIn())
        self.muxEvents.append(envChan)
        self.envs.append(newEnv)
        
        return newEnv

    def forwardMainEvents(self):
        def propagate():
            while True:
                event = self.mainEvents.receive()
                print(event)
                for env in self.muxEvents:
                    env.send(event)
        threading.Thread(target=propagate, daemon=True).start()
    
    def broadcast(self):
        def brdcstEvents():
            ready = self.main.getLock()
            with ready:
                ready.wait() 
            while True:
                for env in range(len(self.envs)):
                    event = self.muxEvents[env].receive()
                    self.mainEvents.send(event)
                    for e in range(len(self.envs)):
                        if self.envs[e].id != self.envs[env].id and event is not None:
                            self.muxEvents[e].send(event)
        
        threading.Thread(target=(brdcstEvents), daemon=True).start() 

    def run(self) -> None:
        '''
        Begin drawing and event threads
        '''
        # self.listenMainEvents.start()
        # self.broadcastEvents.start()
        self.forwardMainEvents()
        self.broadcast()