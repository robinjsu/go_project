from typing import List
from PIL import Image, ImageDraw, ImageFont
import random as rand

from pyGui import *
from pyGui.utils import *
from const import *


broadcast = Event.BroadcastType()
class WordList(Env):
    bounds: Box
    wordList: List

    def __init__(self, box: Box, id=rand.randint(0,100), name=''):
        super().__init__(id=id, threadName=name)
        self.bounds = box
        self.wordList = []
    
    def addToList(self, word: Word):
        self.wordList.append(word)
    
    def saveWordList(self, filename='new_word_list.txt'):
        with open(filename, 'w') as f:
            for wd in self.wordList:
                defs = ''
                for d in wd.definitions:
                    defs += f'[{d[0]}] {d[1]}; '
                wordLn = f'{wd.text.rstrip(trailing_chars).lstrip(trailing_chars)}: {defs[:-2]}\n'
                f.write(wordLn)

    def onBroadcast(self, event: BroadcastEvent):
        if event.event == broadcast.SAVE:
            self.addToList(event.obj)
        if event.event == broadcast.CLOSE:
            self.saveWordList()
            
