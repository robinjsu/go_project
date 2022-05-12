from typing import Tuple, List
from PIL import ImageFont
import io, os, math

class Point:
    x: int
    y: int
     
    def __init__(self, x=0, y=0):
        self.x = x
        self.y = y
    
    def add(self, x, y):
        self.x += x
        self.y += y
    
    def copy(self):
        return Point(self.x, self.y)
class Box:
    x0: int
    y0: int
    x1: int
    y1: int

    def __init__(self, x0=0, y0=0, x1=0, y1=0):
        self.x0 = x0
        self.y0 = y0
        self.x1 = x1
        self.y1 = y1

    def setBoxDims(self, p=Point(0,0)):
        self.x0 = 0
        self.y0 = 0
        self.x1 = p.x
        self.y1 = p.y

    def contains(self, p: Point) -> bool:
        '''
        Returns whether the given point is within the bounds of the box. Exclusive of the Box.x1 and Box.y1 values
        :param p: Point object to test
        '''
        return p.x >= math.floor(self.x0) and p.x < math.ceil(self.x1) and p.y >= math.floor(self.y0) and p.y < math.ceil(self.y1)

    def add(self, p: Point):
        self.x1 += p.x
        self.y1 += p.y
    
    def move(self, p: Point):
        self.x0 += p.x
        self.y0 += p.y
        self.x1 += p.x
        self.y1 += p.y
    
    def copy(self):
        return Box(self.x0, self.y0, self.x1, self.y1)
    
    def size(self):
        w = math.ceil(abs(self.x1 - self.x0))
        h = math.ceil(abs(self.y1 - self.y0))
        return w,h

class Word:
    text: str
    box: Box

    def __init__(self, text, box):
        self.text = text
        self.box = box
class Line:
    line:str
    size: Tuple[int]
    words: List[Word]

    def __init__(self, line, size, words):
        self.line = line
        self.size = size
        self.words = words

    def position(self):
        pass
    def highlight(self):
        pass

def loadFont(filepath, fontSize) -> ImageFont.ImageFont:
    return ImageFont.truetype(filepath, size=fontSize)

def loadFile(filepath) -> Tuple[io.FileIO, int]:
    stats = os.stat(filepath)
    textObj = open(filepath, 'r')
    text = textObj.read()
    return text, stats.st_size

def makeLineBreak(line: str) -> int:
        if '\n' in line:
            return line.find('\n')
        else:
            return line.rfind(' ')

