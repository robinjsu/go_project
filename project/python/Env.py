from queue import Queue
from typing import Dict, Any
from Channel import Channel

class Env:
    event_chan: Channel
    draw: Queue
    parent: Any

    def __init__(self, parent_env, relay=None) -> None:
        if relay == None:
            self.event_chan = Channel()
        else:
            self.event_chan = relay
        self.parent = parent_env
    
    def poll_events(self):
        while True:
            self.event_chan.receive()
    
    def send(self, item):
        self.event_chan.qIn.put(item)

class Mux(Env):
    envs: Dict[Any, Env]

    def __init__(self):
        super().__init__(None)
        self.envs = dict()
    
    def makeNewEnv(self, title) -> Env:
        newEnv = Env(self, Channel())
        self.envs[title] = newEnv

        return newEnv

    def relayMsg(self, subEnvTitle, msg):
        self.envs[subEnvTitle].send(msg)