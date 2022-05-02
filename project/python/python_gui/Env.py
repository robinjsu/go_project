from queue import Queue
from typing import Callable, List, Any
import PIL as pil
from Channel import DrawChan, EventChan
import threading
class Env:
    events: EventChan
    draw: DrawChan

    drawLock: threading.RLock
    drawStream: threading.Thread

    def __init__(self, main: bool):
        if main == True:
            self.events = EventChan(Queue(), Queue())
        self.draw = DrawChan()

    def setEventChan(self, eventsIn, eventsOut):
        self.events = EventChan(eventsIn, eventsOut)
 
    def initDrawThread(self, func: Callable[..., Any]):
        self.drawLock = threading.RLock()
        self.drawStream = threading.Thread(target=func, args=(self.drawLock,), name=f'{func.__repr__}')

    def eventChan(self) -> EventChan:
        return self.events
    
    def drawChan(self) -> DrawChan:
        return self.draw
    
    def drawFunc(self, lock):
        pass
    
    def run(self) -> None:
        '''
        run should initialize separate threads for drawing and events
        '''
        self.initDrawThread(self.drawFunc)
        self.events.open()
        self.drawStream.start()



# TODO: the mux will likely have to handling thread scheduling, in some kind of round-robin fashion probably?
# mux should receive events from the main env and pass on to each of the sub envs
# mux will receive events from the sub envs and pass on to the main env
class Mux(Env):
    '''
    The Mux class acts as a multiplexer for the main environment
    Will throw an error if it is initialized without a main.
    '''

    events: EventChan
    draw: DrawChan

    drawLock: threading.RLock
    drawStream: threading.Thread

    main: Env
    envs: List[Env]
    mainEvents: EventChan
    muxEvents: List[EventChan]

    def __init__(self, mainEnv: Env):
        assert mainEnv != None, f'Mux must be connected to a main environment'
        super().__init__(False)
        self.events = EventChan(Queue(), Queue())
        self.main = mainEnv
        self.envs = []
        self.mainEvents = None
        self.muxEvents = []
    
    def receiveMainEvents(self) -> None:
        '''
        Attach main events queue to mux
        '''
        self.mainEvents = EventChan(self.main.events.getEventsOut(), Queue())
  
    def addEnv(self) -> Env:  
        '''
        Add new env to the mux. Associate proper queues between the mux env and new component env
        '''      
        newEnv = Env(False)
        newEnv.setEventChan(self.eventChan().getEventsOut(), Queue())
        self.draw = newEnv.drawChan()
        self.envs.append(newEnv)
        self.muxEvents.append(newEnv.eventChan().getEventsOut())

        return newEnv

    def eventChan(self) -> EventChan:
        return self.events

    def drawChan(self) -> DrawChan:
        return self.draw

    
    def run(self):
        self.initDrawThread(self.drawFunc)
        print('open events stream')
        self.events.open()
        print('open draw stream')
        self.drawStream.start()
        
        print('start envs')
        for e in self.envs:
            e.run()
        print('run main')
        self.main.run()
        
        
        
