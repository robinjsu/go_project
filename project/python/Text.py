from typing import List, NamedTuple, Tuple
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand, threading

from pyGui import *
from pyGui.utils import *
from const import *

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

    def __init__(self, box: tuple, id=rand.randint(0,100)):
        super().__init__(id=id)
        self.padding = 5
        self.bounds = Box(box[0], box[1], box[2], box[3])
        self.width = self.bounds.x1 - self.bounds.x0
        self.height = self.bounds.y1 - self.bounds.y0
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
    
    def setFont(self, ttf):
        self.font = ttf
        self.pixelsPerLetter = self.font.getlength('A')
        self.charsPerWidth = self.padW // math.ceil(self.pixelsPerLetter)


    # assume monospaced font for now
    def formatText(self, fileObj, fileSz):
        lines = []
        anchor = Point(0,0)
        buffer = io.BufferedReader(fileObj, fileSz)
        text = buffer.peek()
        text = text.decode('utf-8')

        idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
        line = text[:idx].rstrip(' \n')
        sz = self.font.getsize(line)
        lines.append(
            Line(
                line, 
                sz, 
                self.setTextPos(line, anchor)
            )
        )
        anchor.add(0, (sz[1] + lineSpacing))
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
                    self.setTextPos(line, anchor)
                )
            )
            anchor.add(0, (sz[1] + lineSpacing))
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
            paddedBox = ImageOps.pad(
                Image.new("RGBA", (self.width, self.height)), (self.padW, self.padH)
            )
            bg = Image.new("RGBA", (self.width, self.height), Colors().white)
            c = Colors()
            drawCtx = ImageDraw.ImageDraw(paddedBox)
            for l in lines:
                for w in l.words:
                    drawCtx.text(
                        (w.box.x0, w.box.y0), 
                        w.text, 
                        c.black, 
                        self.font
                    )

            bg.alpha_composite(paddedBox, (MARGIN,MARGIN))
            baseImg.alpha_composite(bg, (MARGIN, MARGIN))
            return baseImg
        return drawText
    
    

    def onMouseClick(self, action):
        if action == event.MouseDown():
            self.setFont(ImageFont.truetype('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf'))
            fileObj, fileSz = loadFile('alice.txt')
            lines = self.formatText(fileObj, fileSz)
            self.draw.send(self.setText(lines))

    def resize(self):
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
        paddedBox = ImageOps.pad(
            Image.new("RGBA", (self.width, self.height)), (self.padW, self.padH)
        )

    # def run(self) -> None:
    #     def startThread():
    #         with self.ready:
    #             self.ready.wait()
    #         self.init()
    #         while True:
    #             event = self.eventChan().receive()
    #             if type(event) == MouseEvent:
    #                 self.onMouseClick(event.action)
    #             elif type(event) == KeyEvent:
    #                 self.onKeyPress(event.key)
    #     threading.Thread(target=startThread, name="DisplayThread", daemon=True).start() 

# if __name__== '__main__':
    # lines = []
    # font = ImageFont.truetype('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf')
    # fileObj, fileSz = loadFile('alice.txt')
    # buffer = io.BufferedReader(fileObj, fileSz)
    # text = buffer.read()
    # text = text.decode('utf-8')
    
    
    
