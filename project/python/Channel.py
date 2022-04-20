import logging
import queue
from queue import Queue, Empty
import typing, threading



# base class for event "channels"
class Channel:
    qIn: Queue
    qOut: Queue
    close: threading.Event

    def __init__(self):
        self.qIn = Queue()
        self.qOut = Queue()
        self.close = threading.Event()


    def receive(self) -> None:
        while self.close.is_set() != True:
            item = self.qIn.get()
            print(f'received {item}')
            self.qIn.task_done()