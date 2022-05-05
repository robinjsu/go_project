from queue import Queue
from typing import List, Callable
import threading
from PIL import Image

from .Channel import DrawChan, EventChan
from .Event import MouseEvent, KeyEvent

# TODO: how to poll for events in a non-blocking way, such that the env can write custom callbacks for when an event is received?
class Env:
    events: EventChan
    draw: DrawChan
    window: bool
    ready: threading.Condition

    def __init__(self, main=False):
        self.window = main
        self.events = EventChan(Queue(), Queue())
        self.draw = DrawChan()
    
    def addCond(self, condLock):
        self.ready = condLock
    
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
    
    def drawImg(self, drawCommand: Callable[...,Image.Image]):
        self.draw.send(drawCommand)

    def init(self):
        pass
    
    def run(self) -> None:
        '''
        should start running it's logic and callback listeners on a thread separate from the main Python interpreter thread
        '''
        def startThreads():
            with self.ready:
                self.ready.wait()
            self.init()
            while True:
                event = self.eventChan().receive()
                if type(event) == MouseEvent:
                    self.onMouseClick(event.action)
                elif type(event) == KeyEvent:
                    self.onKeyPress(event.key)
        threading.Thread(target=startThreads, name="DisplayThread", daemon=True).start() 


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

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'missing Main Env: Mux must be created from an existing Env'
    
        self.main = mainEnv
        self.envs = []

  
    def addEnv(self, newEnv: Env, lock) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        newEnv.setEventChan(self.main.eventChan().getEventsOut(),self.main.eventChan().getEventsIn())
        newEnv.setDrawChan(self.main.drawChan())
        newEnv.addCond(lock)
        self.envs.append(newEnv)

        return newEnv

    # def run(self) -> None:
    #     '''
    #     Begin drawing and event threads
    #     '''
    #     self.mainEventStream.start()
    #     self.muxEventStream.start()
    #     self.drawStream.start()