from queue import Queue
from typing import Dict, Any
from Channel import Channel

class Env:
    event_chan: Channel
    draw = Queue

    def __init__(self) -> None:
        self.event_chan = Channel()
        self.draw = Channel()
    
    def poll_events(self):
        while True:
            self.event_chan.receive()
    
    def send(self, item):
        self.event_chan.qIn.put(item)

class Mux(Env):
    envs: Dict[Any, Env]

    def __init__(self):
        super().__init__()
    
    def makeNewEnv(self, title) -> Env:
        newEnv = Env()
        self.envs[title] = newEnv
        return newEnv
