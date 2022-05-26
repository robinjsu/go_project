from typing import Tuple, List
from PIL import ImageFont
import io, os, math, requests


DICT_API_KEY = os.getenv('DICT_API_KEY')
assert DICT_API_KEY is not None, 'no api key provided'
WORDS_URL = 'https://wordsapiv1.p.rapidapi.com/words'
trailing_chars = ' ,.;:!?()"\'-“”‘’'

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
    definitions: List[Tuple[str, str]]

    def __init__(self, text, box, defn=None):
        self.text = text
        self.box = box
        self.definitions = defn

    def getDefinitions(self):
        assert self.text != None or self.text != '', f'word has not been initialized properly - Word.text: {self.text}'
        word = (self.text).rstrip(trailing_chars)
        word = word.lstrip(trailing_chars)
        headers = {
            'X-RapidAPI-Key': DICT_API_KEY,
            'Accept': '*/*',
            'Connection': 'keep-alive'
        }
        resp = requests.get(f'{WORDS_URL}/{word}/definitions', headers=headers)
        if resp.status_code == 200:
            defs = []
            definitions = resp.json()['definitions']
            for df in definitions:
                defs.append(tuple([df['partOfSpeech'], df['definition']]))
            self.definitions = defs
            return defs

    
    def copy(self):
        return Word(self.text, self.box, self.definitions)

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

def getFontSize(font) -> Point:
    wd = font.getlength('A')
    ht = font.getsize('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz')
    return Point(wd, ht[1])

def loadFile(filepath) -> Tuple[io.FileIO, int]:
    stats = os.stat(filepath)
    textObj = open(filepath, 'r')
    text = textObj.read()
    return text, stats.st_size

 # assume monospaced font for now
def formatText(text: str, charsPerWidth):
    '''
    Format a body of text into lines that won't exceed the calculated width of the component.
    :param text: a string representing entire body of text to display
    '''
    plainText = []
    if len(text) > charsPerWidth:
        idx = makeLineBreak(text[:charsPerWidth]) + 1
        line = text[:idx].rstrip(' \n')
        plainText.append(line)
        while text != '':
            text = text[idx:]   
            if len(text) > charsPerWidth:
                idx = makeLineBreak(text[:charsPerWidth]) + 1
            else:
                idx = len(text)
            line = text[:idx].rstrip(' \n')
            plainText.append(line)
    return plainText

def makeLineBreak(line: str) -> int:
    if '\n' in line:
        return line.find('\n')
    else:
        return line.rfind(' ')