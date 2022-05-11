from typing import List
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand

from pyGui import *
from pyGui.utils import *
from const import *


''' 
TODO: 
 - add hint for paging as way to read
 - how to best set an universal anchor for where the text begins?
'''

lineSpacing = 4
input = InputEvent
color = Colors()
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
    plainText: List
    numPages: int
    anchor = tuple

    def __init__(self, box: tuple, id=rand.randint(0,100)):
        super().__init__(id=id)
        self.padding = MARGIN
        self.bounds = Box(box[0], box[1], box[2], box[3])
        self.width = abs(self.bounds.x1 - self.bounds.x0)
        self.height = abs(self.bounds.y1 - self.bounds.y0)
        self.padW = self.width - (self.padding * 2)
        self.padH = self.height - (self.padding * 2)
        self.page = None
        self.anchor = tuple([self.bounds.x0+(self.padding*2), self.bounds.y0+(self.padding*2)])
    

    def setFont(self, ttf):
        self.font = ttf
        self.pixelsPerLetter = self.font.getlength('A')
        self.charsPerWidth = self.padW // math.ceil(self.pixelsPerLetter)


    # assume monospaced font for now
    def formatText(self, text):
        lines = []
        self.plainText = []
        idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
        line = text[:idx].rstrip(' \n')
        # sz = self.font.getsize(line)
        # lines.append(
        #     Line(
        #         line, 
        #         sz,
        #         None
        #     )
        # )
        self.plainText.append(line)
        if len(text) > idx - 1:
            text = text[idx:]

        while line != '':
            idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
            line = text[:idx].rstrip(' \n')
            # sz = self.font.getsize(line)
            # lines.append(
            #     Line(
            #         line, 
            #         sz,
            #         None
            #     )
            # )
            self.plainText.append(line)
            if len(text) > idx - 1:
                text = text[idx:]

        # return lines


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


    def setText(self, lines: List[str]):
        txtLines = []
        anchor = Point(self.anchor[0], self.anchor[1])
        for l in range(len(lines)):
            txtLine = Line(
                line=lines[l],
                size=self.font.getsize(lines[l]),
                words=self.setTextPos(lines[l], anchor)
            )
            txtLines.append(txtLine)
            anchor.add(0, ((self.font.getsize(lines[l]))[1] + lineSpacing))
        self.lines = txtLines

        def drawText(baseImg: Image.Image) -> Image.Image:
            paddedBox = ImageOps.pad(
                Image.new("RGBA", (self.width, self.height), color.white), (self.padW, self.padH)
            )
            bg = Image.new("RGBA", (self.width, self.height), Colors().ultra)
            drawCtx = ImageDraw.ImageDraw(paddedBox)
            for l in self.lines:
                for w in l.words:
                    drawCtx.text(
                        (w.box.x0, w.box.y0), 
                        w.text, 
                        color.black, 
                        self.font,
                        anchor='la'
                    )
                # anchor.add(0, ((self.font.getsize(lines[l].line))[1] + lineSpacing))
            # bg.alpha_composite(paddedBox)
            baseImg.alpha_composite(paddedBox, (self.bounds.x0, self.bounds.y0))
            return baseImg

        return drawText
    
    def findWord(self, p: Point):
        # lines = self.lines[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]
        for l in self.lines:
            for w in l.words:
                if w.box.contains(p):
                    print(p, w.box.x0, w.box.y0, w.box.x1, w.box.y1)
                    def highlightWord(base: Image.Image) -> Image.Image:
                        textbox = Image.new("RGBA", w.box.size())
                        drawWord = ImageDraw.ImageDraw(textbox)
                        drawWord.rectangle(((0,0), w.box.size()), color.paleBlueTransp)
                        drawWord.text((0,0), w.text, color.black, self.font)
                        base.alpha_composite(textbox, (w.box.x0, w.box.y0), (0,0))
                        return base
                    return highlightWord


    def onMouseClick(self, event: MouseEvent):
        pt = Point(event.xpos, event.ypos)
        if event.action == input.MouseDown:
            if self.page == None:
                self.page = 0
                self.draw.send(self.setText(self.plainText[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))
            else:
                self.draw.send(self.setText(self.plainText[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))
                self.draw.send(self.findWord(pt))

    def onKeyPress(self, keyEvent: KeyEvent):
        if keyEvent.key == input.ArrowDown and keyEvent.action == 1:
            if self.page < self.numPages-1:
                self.page += 1
            self.draw.send(self.setText(self.plainText[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))
        elif keyEvent.key == input.ArrowUp and keyEvent.action == 1:
            if self.page > 0:
                self.page -= 1
            self.draw.send(self.setText(self.plainText[self.page*MAXLINES:((self.page*MAXLINES)+MAXLINES)]))


    def resize(self):
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
        
    
    def init(self):
        self.setFont(loadFont('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf'))
        text, _ = loadFile('alice.txt')
        self.lines = self.formatText(text)
        self.numPages = int(math.ceil(len(self.plainText) / MAXLINES))