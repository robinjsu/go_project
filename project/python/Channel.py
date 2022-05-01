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
    close: bool

    def __init__(self) -> None:
        self.In = Queue()
        self.Out = Queue()
    
    def open(self):
        '''
        Start separate thread to poll for incoming and outgoing events
        '''
        eventThread = threading.Thread(target=self.poll_events, name="eventsThread")
        eventThread.start()
    
    def poll_events(self) -> None:
        while not self.close:
            event = self.In.get(block=True)
            print(f'event received: {event}')
            self.In.task_done()
            self.Out.put(event, block=True)
        while True:
            try:
                event = self.In.get()
            except Empty:
                return
            self.In.task_done()
            self.Out.put(event)
            
    def close(self) -> None:
        self.close = True
        
    # def send(self, event: Any):
    #     pass

    # def receive(self) -> Any:
    #     pass


class DrawChan(Queue):
    def __init__(self) -> None:
        return super().__init__()

    def send(self, img: Callable[..., Image.Image]):
        self.put(img)

    def receive(self) -> Callable[..., Image.Image]:
        try:
            drawCommand = self.get(block=True, timeout=0.1)
        except Empty:
            print("no draw command issued")
            return None
        self.task_done()
        return drawCommand  