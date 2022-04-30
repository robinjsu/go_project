from queue import Queue
from typing import Any, Callable
import PIL as pil
from Channel import DrawChan, EventChan
class Env:
    Events: EventChan
    Draw: DrawChan

    def __init__(self):
        events = EventChan()
        events.open()

    def draw(self, drawFunc: Callable[...,pil.Image.Image]) -> None:
        self.Draw.send(drawFunc)

# TODO: the mux will likely have to handling thread scheduling, in some kind of round-robin fashion probably?
class Mux(Env):
    def __init__(self):
        super().__init__()

    
