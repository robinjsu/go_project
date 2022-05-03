from queue import Queue
from typing import Callable, List, Any
import PIL as pil
from .Channel import DrawChan, EventChan
import threading
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
        run should initialize separate threads for drawing and events
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

    events: EventChan
    draw: DrawChan

    drawStream: threading.Thread
    eventStream: threading.Thread

    main: Env
    envs: List[Env]
    # mainEvents: EventChan
    # muxEvents: List[EventChan]

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'missing Main Env: Mux must be created from an existing Env'
        super().__init__(False)
        # self.events = EventChan(Queue(), Queue())
        self.main = mainEnv
        self.envs = []
        self.handleDrawCommands()
        self.pollEvents()
        # self.mainEvents = None
        # self.muxEvents = []
    
    def receiveMainEvents(self) -> None:
        '''
        Attach main events queue to mux
        '''
        self.events = EventChan(self.main.events.getEventsOut(), Queue())
    
    def pollEvents(self) -> None:
        '''
        Create thread to poll for events from main Env
        '''
        def poll():
            while True:
                event = self.events.receive()
                if event is not None:
                    self.events.send(event)
        self.eventStream = threading.Thread(target=poll, daemon=True)

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
        newEnv.setEventChan(self.eventChan().getEventsOut(), Queue())
        newEnv.setDrawChan(self.draw)
        self.envs.append(newEnv)

        return newEnv

    def run(self) -> None:
        '''
        Start drawing and event threads
        '''
        # TODO: put logic here for polling for drawing events from envs, and to send to the window
        self.eventStream.start()
        self.drawStream.start()
    
        
