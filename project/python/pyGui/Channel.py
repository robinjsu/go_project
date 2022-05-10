from queue import Queue, Empty
from PIL import Image
from typing import Any, Callable
import threading

class EventChan:
    In: Queue
    Out: Queue
    closed: bool
    eventThread: threading.Thread

    def __init__(self, chanIn, chanOut) -> None:
        self.In = chanIn
        self.Out = chanOut
        self.closed = False
    
    def receive(self) -> Any:
        try:
            event = self.In.get(block=True)
        except Empty:
            return None
        self.In.task_done()
        return event 
    
    def receiveTimeout(self, timeout: float) -> Any:
        try:
            event = self.In.get(block=True, timeout=timeout)
        except Empty:
            return None
        self.In.task_done()
        return event
    
    def send(self, event: Any) -> None:
        self.Out.put(event, block=True)
            
    def close(self) -> None:
        self.closed = True

    def getEventsIn(self) -> Queue:
        return self.In
    
    def getEventsOut(self) -> Queue:
        return self.Out
    
    def eventError(self) -> Exception:
        pass

class DrawChan(Queue):
    def __init__(self) -> None:
        return super().__init__()

    def send(self, img: Callable[..., Image.Image]):
        self.put(img, block=True)

    def receive(self) -> Callable[..., Image.Image]:
        try:
            drawCommand = self.get(block=True)
        except Empty:
            return None
        self.task_done()
        return drawCommand  
