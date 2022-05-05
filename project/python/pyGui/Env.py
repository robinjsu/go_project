from queue import Queue
from typing import List
import threading

from .Channel import DrawChan, EventChan

class Env:
    events: EventChan
    draw: DrawChan

    def __init__(self, main: bool):
        self.events = EventChan(Queue(), Queue())
        self.draw = DrawChan()
    
    def eventChan(self) -> EventChan:
        return self.events
    
    def drawChan(self) -> DrawChan:
        return self.draw

    def setEventChan(self, eventsIn, eventsOut):
        self.events = EventChan(eventsIn, eventsOut)
    
    def setDrawChan(self, drawChan):
        self.draw = drawChan
    
    def run(self) -> None:
        '''
        run should initialize separate threads for drawing and events streams
        '''
        pass


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
        # super().__init__(False)
    
        self.main = mainEnv
        self.envs = []
  
    def addEnv(self) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        newEnv = Env(False)
        newEnv.setEventChan(self.main.eventChan().getEventsOut(),self.main.eventChan().getEventsIn())
        newEnv.setDrawChan(self.main.drawChan())
        self.envs.append(newEnv)

        return newEnv

    def run(self) -> None:
        '''
        Begin drawing and event threads
        '''
        self.mainEventStream.start()
        self.muxEventStream.start()
        self.drawStream.start()