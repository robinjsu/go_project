from queue import Queue
import threading
from PyQt5 import QtGui as qtgui, QtCore as qtc, QtWidgets as qtw
from Channel import Channel

class Env:
    event_chan: Channel
    draw = Queue

    def __init__(self) -> None:
        self.event_chan = Channel()
        self.draw = Queue()
    
    def poll_events(self):
        while True:
            self.event_chan.receive()
    
    def send(self, item):
        self.event_chan.qIn.put(item)