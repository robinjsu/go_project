from typing import List, NamedTuple, Tuple
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand, threading

from pyGui import *
from pyGui.utils import *
from const import *


# TODO: fix anchor starting point, as it needs to adjust to paging!

lineSpacing = 4
event = Event()

class Text(Env):
    bounds: Box
    font: ImageFont.ImageFont
    padding: int
    width: int
    height: int
    padW: int
    padH: int
    pixelsPerLetter: float
    charsPerWidth: int
    page: int
    lines: List
    numPages: int

    def __init__(self, box: tuple, id=rand.randint(0,100)):
        super().__init__(id=id)
        self.padding = 5
        self.bounds = Box(box[0], box[1], box[2], box[3])
        self.width = self.bounds.x1 - self.bounds.x0
        self.height = self.bounds.y1 - self.bounds.y0
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
        self.page = None
    
    def setFont(self, ttf):
        self.font = ttf
        self.pixelsPerLetter = self.font.getlength('A')
        self.charsPerWidth = self.padW // math.ceil(self.pixelsPerLetter)


    # assume monospaced font for now
    def formatText(self, text):
        lines = []
        idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
        line = text[:idx].rstrip(' \n')
        sz = self.font.getsize(line)
        lines.append(
            Line(
                line, 
                sz,
                None
            )
        )
        if len(text) > idx - 1:
            text = text[idx:]

        while line != '':
            idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
            line = text[:idx].rstrip(' \n')
            sz = self.font.getsize(line)
            lines.append(
                Line(
                    line, 
                    sz,
                    None
                )
            )
            if len(text) > idx - 1:
                text = text[idx:]

        return lines

    def makeLineBreak(self, line: str) -> int:
        if '\n' in line:
            return line.find('\n')
        else:
            return line.rfind(' ')

    def setTextPos(self, line: str, anchor: Point):
        words = line.split(' ')
        wrdSzs = []
        for w in words:
            fntSz = self.font.getsize(w)
            wrdSzs.append(Point(fntSz[0], fntSz[1]))
        anchorLoc = anchor.copy()
        wordsPos = []
        for i in range(len(words)):
            box = Box()
            box.setBoxDims(p=Point((wrdSzs[i]).x, (wrdSzs[i]).y))
            box.move(anchorLoc)
            wordsPos.append(Word(words[i], box))
            anchorLoc.add((wrdSzs[i].x + int(math.ceil(self.pixelsPerLetter))), 0)
        return wordsPos


    def setText(self, lines: List[Line]):
        def drawText(baseImg: Image.Image) -> Image.Image:
            anchor = Point(0,0)
            paddedBox = ImageOps.pad(
                Image.new("RGBA", (self.width, self.height)), (self.padW, self.padH)
            )
            bg = Image.new("RGBA", (self.width, self.height), Colors().white)
            c = Colors()
            drawCtx = ImageDraw.ImageDraw(paddedBox)
            for l in lines:
                txtLine = Line(line=l.line, size=l.size, words=self.setTextPos(l.line, anchor))
                for w in txtLine.words:
                    drawCtx.text(
                        (w.box.x0, w.box.y0), 
                        w.text, 
                        c.black, 
                        self.font
                    )
                anchor.add(0, ((self.font.getsize(l.line))[1] + lineSpacing))
            bg.alpha_composite(paddedBox, (MARGIN,MARGIN))
            baseImg.alpha_composite(bg, (MARGIN, MARGIN))
            return baseImg
        return drawText
    
    

    def onMouseClick(self, action):
        if action == event.MouseDown():
            if self.page == None:
                self.page = 0
                self.draw.send(self.setText(self.lines[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))
    
    def onKeyPress(self, keyEvent: KeyEvent):
        if keyEvent.key == event.ArrowDown() and keyEvent.action == 1:
            if self.page < self.numPages-1:
                self.page += 1
            self.draw.send(self.setText(self.lines[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))
        elif keyEvent.key == event.ArrowUp() and keyEvent.action == 1:
            if self.page > 0:
                self.page -= 1
            self.draw.send(self.setText(self.lines[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))


    def resize(self):
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
        
    
    def init(self):
        self.setFont(ImageFont.truetype('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf'))
        text, _ = loadFile('alice.txt')
        self.lines = self.formatText(text)
        self.numPages = int(math.ceil(len(self.lines) / MAXLINES))