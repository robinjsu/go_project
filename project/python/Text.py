from typing import List, Callable
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand

from pyGui import *
# from pyGui.Event import BroadcastType
from pyGui.utils import *
from const import *


''' 
TODO: 
 - add hint for paging as way to read
'''

lineSpacing = 10
fontSize = 24
input = Event.InputType()
broadcast = Event.BroadcastType()
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
    linesPerPage: int
    page: int
    lines: List
    plainText: List
    numPages: int
    anchor = tuple

    def __init__(self, box: Box, id=rand.randint(0,100), name=''):
        super().__init__(id=id, threadName=name)
        self.padding = MARGIN
        self.bounds = box
        self.width = abs(self.bounds.x1 - self.bounds.x0)
        self.height = abs(self.bounds.y1 - self.bounds.y0)
        self.padW = self.width - (self.padding * 2)
        self.padH = self.height - (self.padding * 2)
        self.page = None
        self.anchor = tuple([self.bounds.x0+(self.padding*2), self.bounds.y0+(self.padding*2)])
    

    def setFont(self, ttf, sz):
        '''
        Set font style for this component. Accepts TrueType standard font styles.
        :param ttf: an ImageFont object (from the Python Pillow library)
        '''
        self.font = loadFont(ttf, sz)
        self.pixelsPerLetter, letterHeight = self.font.getsize('A')
        self.charsPerWidth = self.padW // self.pixelsPerLetter
        self.linesPerPage = int(math.floor(self.padH / (letterHeight + lineSpacing)))


    def setTextPos(self, line: str, anchor: Point):
        '''
        Set word positions for a single line of text, given a text anchor as the starting position.
        :param line: a line of text as a string.
        :param anchor: a Point object representing the starting position for the text line, as the top left corner of the text.
        '''
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
        '''
        Set coordinate positions for given lines of text. Calculates position for each individual word.
        :param lines: a list of strings represented text. Text should have been previously split into lines using
            `formatText` to ensure that text will not overflow beyond component borders.
        '''
        txtLines = []
        anchor = Point(self.anchor[0], self.anchor[1])
        for line in lines:
            txtLine = Line(
                line=line,
                size=self.font.getsize(line),
                words=self.setTextPos(line, anchor)
            )
            txtLines.append(txtLine)
            anchor.add(0, ((self.font.getsize(line))[1] + lineSpacing))
        self.lines = txtLines

        def drawText(baseImg: Image.Image) -> Image.Image:
            paddedBox = ImageOps.pad(
                Image.new("RGBA", (self.width, self.height), color.white), (self.padW, self.padH)
            )
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
            baseImg.alpha_composite(paddedBox, (self.bounds.x0, self.bounds.y0))
            return baseImg

        return drawText
    

    def findWord(self, p: Point) -> tuple:
        '''
        Find and highlight the word clicked on with the mouse.
        :param p: Point object that represents where the user clicked within the text area.
        '''
        wrd = None
        for l in self.lines:
            for w in l.words:
                if w.box.contains(p):
                    wrd = w.copy()
                    def highlightWord(base: Image.Image) -> Image.Image:
                        textbox = Image.new("RGBA", w.box.size())
                        drawWord = ImageDraw.ImageDraw(textbox)
                        drawWord.rectangle(((0,0), w.box.size()), color.paleBlueTransp)
                        drawWord.text((0,0), w.text, color.black, self.font)
                        base.alpha_composite(textbox, (w.box.x0, w.box.y0), (0,0))
                        return base
                    return (highlightWord, wrd)
        return None, None


    def onMouseClick(self, event: MouseEvent):
        '''
        Callback function that responds to a mouse button being pressed or released.
        :param event: a MouseEvent object that represents the mouse button and action that occurred.
        '''
        pt = Point(event.xpos, event.ypos)
        if event.action == input.Press:
            if self.bounds.contains(pt):
                highlightFunc, word = self.findWord(pt)
                if highlightFunc != None:
                    self.sendEvent(BroadcastEvent(broadcast.DEFINE, word))
                    self.drawImg(self.setText(self.plainText[self.page*self.linesPerPage:((self.page*self.linesPerPage)+self.linesPerPage)]))
                    self.drawImg(highlightFunc)


    def onKeyPress(self, keyEvent: KeyEvent):
        '''
        Callback function that responds to a key being pressed or released.
        :param event: a KeyEvent object that represents the key pressed and action that occurred.
        '''
        if keyEvent.key == input.ArrowDown and keyEvent.action == input.Press:
            if self.page < self.numPages-1:
                self.page += 1
            self.drawImg(self.setText(self.plainText[self.page*self.linesPerPage:((self.page*self.linesPerPage)+self.linesPerPage)]))
        elif keyEvent.key == input.ArrowUp and keyEvent.action == input.Press:
            if self.page > 0:
                self.page -= 1
            self.drawImg(self.setText(self.plainText[self.page*self.linesPerPage:((self.page*self.linesPerPage)+self.linesPerPage)]))


    def resize(self):
        self.padW = self.width - (self.padding * 2)
        self.padH = self.width - (self.padding * 2)
        
    
    def init(self):
        self.page = 0
        self.setFont('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf', fontSize)
        text, _ = loadFile('alice.txt')
        self.plainText = formatText(text, self.charsPerWidth)
        self.numPages = int(math.ceil(len(self.plainText) / self.linesPerPage))
        self.drawImg(self.setText(self.plainText[self.page*self.linesPerPage:((self.page*self.linesPerPage)+self.linesPerPage)]))