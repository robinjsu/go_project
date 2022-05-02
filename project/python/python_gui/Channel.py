import queue
from queue import Queue, Empty
from PIL import Image
from typing import Any, Callable
import threading

# base class for event "channels"
# TODO: all logic for handling Queues should happen here, and encapsulated
class EventChan:
    In: Queue
    Out: Queue
    closed: bool

    def __init__(self, chanIn, chanOut) -> None:
        self.In = chanIn
        self.Out = chanOut
        self.closed = False
    
    def open(self):
        '''
        Start separate thread to poll for incoming and outgoing events
        '''
        eventThread = threading.Thread(target=self.poll_events, name=f'EventsThread', daemon=True)
        print('start events thread')
        eventThread.start()
    
    def poll_events(self) -> None:
        while not self.closed:
            event = self.In.get(block=True)
            self.In.task_done()
            self.Out.put(event, block=True)
        # TODO: is this part necessary? it is meant to clear out the queue
        while True:
            try:
                event = self.In.get()
            except Empty:
                return
            self.In.task_done()
            self.Out.put(event)
            
    def close(self) -> None:
        self.closed = True
    
    def getEventsOut(self) -> Queue:
        return self.Out
    
    def eventError(self) -> Exception:
        pass

class DrawChan(Queue):
    def __init__(self) -> None:
        return super().__init__()

    def send(self, img: Callable[..., Image.Image]):
        self.put(img)

    def receive(self) -> Callable[..., Image.Image]:
        try:
            drawCommand = self.get(block=True, timeout=0.1)
        except Empty:
            return None
        self.task_done()
        return drawCommand  