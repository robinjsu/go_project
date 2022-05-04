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


# TODO: the mux will likely have to handling thread scheduling, in some kind of round-robin fashion probably?
# mux should receive events from the main env and pass on to each of the sub envs
# mux will receive events from the sub envs and pass on to the main env
class Mux(Env):
    '''
    The Mux class acts as a multiplexer for the main environment. It can create new sub-environments that communicate with the Mux via Channel objects.
    The Mux receives events from the main Env, which it passes to the multiplex Envs. It receives drawing commands from the sub-Envs, which it passes to the main Env.
    Will throw an error if it is initialized without a main Env.
    '''

    mainEvents: EventChan
    muxEvents: EventChan
    draw: DrawChan

    drawStream: threading.Thread
    mainEventStream: threading.Thread
    muxEventStream: threading.Thread

    main: Env
    envs: List[Env]

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'missing Main Env: Mux must be created from an existing Env'
        super().__init__(False)
        # assert self.mainEvents is not None and self.draw is not None, f'events and draw channels not properly initialized'
       
        self.main = mainEnv
        self.envs = []
        self.handleDrawCommands()
        self.pollEvents()
        self.pollMuxEvents()
        mainEventChan = self.main.eventChan()
        self.mainEvents = EventChan(mainEventChan.getEventsOut(), Queue())
        self.muxEvents = EventChan(Queue(), self.main.eventChan().getEventsIn())
    
    def mainEventChan(self):
        return self.mainEvents
    
    def muxEventChan(self):
        return self.muxEvents
    
    def pollEvents(self) -> None:
        '''
        Create separate thread to poll for events incoming from main Env.
        '''
        def poll():
            while True:
                mainEvent = self.mainEvents.receive()
                if mainEvent is not None:
                    self.mainEvents.send(mainEvent)
        self.mainEventStream = threading.Thread(target=poll, daemon=True)


    def pollMuxEvents(self) -> None:
        def poll():
            muxEvent = self.muxEvents.receive()
            if muxEvent is not None:
                self.muxEvents.send(muxEvent)
        self.muxEventStream = threading.Thread(target=poll, daemon=True)
        
    def handleDrawCommands(self) -> None:
        '''
        Create thread to handle seinding drawing commands to the main window
        '''
        drawLock = threading.RLock()
        def pollDrawCmds(lock: threading.RLock):
            lock.acquire()
            while True:
                cmd = self.draw.receive()
                if cmd is not None:
                    self.main.drawChan().send(cmd)
        self.drawStream = threading.Thread(target=pollDrawCmds, args=(drawLock,), daemon=True)
  
    def addEnv(self) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        newEnv = Env(False)
        newEnv.setEventChan(self.mainEventChan().getEventsOut(),self.muxEventChan().getEventsIn())
        newEnv.setDrawChan(self.draw)
        self.envs.append(newEnv)

        return newEnv

    def run(self) -> None:
        '''
        Begin drawing and event threads
        '''
        self.mainEventStream.start()
        self.muxEventStream.start()
        self.drawStream.start()