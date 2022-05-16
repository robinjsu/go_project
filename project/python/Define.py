from cgitb import text
from typing import List
from PIL import Image, ImageDraw, ImageFont, ImageOps
import os, requests, random as rand

from pyGui import *
from pyGui.utils import *
from const import *

fontSize = 16
lineHeight = 5
color = Colors()
class Define(Env):
    word: str
    bounds: Box
    font: ImageFont.ImageFont
    headerFont: ImageFont.ImageFont
    padding: int
    width: int
    height: int
    padW: int
    padH: int
    anchor: Point
    pixelsPerChar: Point
    charsPerWidth: int
    plainText: List[str]

    def __init__(self, box: Box, id=rand.randint(0,100)):
        super().__init__(id=id)
        self.padding = MARGIN
        self.bounds = box
        self.width = abs(self.bounds.x1 - self.bounds.x0)
        self.height = abs(self.bounds.y1 - self.bounds.y0)
        self.padW = self.width - (self.padding * 2)
        self.padH = self.height - (self.padding * 2)
        self.anchor = Point(self.bounds.x0+MARGIN, self.bounds.y0+MARGIN)
    
    def init(self):
        self.setFont('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf', fontSize)

    def setFont(self, ttf, sz):
        self.font = loadFont(ttf, sz)
        self.pixelsPerChar = getFontSize(self.font)
        self.charsPerWidth = self.padW // math.ceil(self.pixelsPerChar.x)

    def setWordHeader(self):
        w = self.word.rstrip(trailing_chars).lstrip(trailing_chars)
        chars = len(w)
        spacingW = (self.charsPerWidth - chars) // 2
        anchorX = spacingW * self.pixelsPerChar.x
        anchorY = self.pixelsPerChar.y
        anchor = int(self.anchor.x + anchorX), int(self.anchor.y + anchorY)
        def drawHeader(baseImg: Image.Image) -> Image.Image:
            bg = Image.new("RGBA", (self.padW, self.pixelsPerChar.y * 3), color.paleBlue)
            textImg = Image.new("RGBA", self.font.getsize(w), color.paleBlue)
            context = ImageDraw.ImageDraw(textImg)
            context.text(
                (0, 0),
                w,
                color.ultra,
                self.font,
                anchor='la'
            )
            baseImg.alpha_composite(bg, (self.bounds.x0+MARGIN, self.bounds.y0+MARGIN))
            baseImg.alpha_composite(textImg, anchor)
            return baseImg
        return drawHeader

    def setDefSection(self, definitions):
        def drawSection(base: Image.Image) -> Image.Image:
            anchor = Point(MARGIN,MARGIN)
            bg = Image.new("RGBA", (self.padW, int(self.padH - (self.pixelsPerChar.y * 3))), color.lightBlue)
            drawCtx = ImageDraw.ImageDraw(bg)
            for defn in definitions:
                formatted = [f'[{defn[0]}]'] + formatText(defn[1], self.charsPerWidth)
                joinedStr = '\n'.join(formatted)
                drawCtx.multiline_text(
                    (anchor.x, anchor.y), 
                    joinedStr, 
                    color.black, 
                    self.font, 
                    anchor='la', 
                    spacing=lineHeight
                )
                bbx = drawCtx.multiline_textbbox((anchor.x, anchor.y), joinedStr, self.font, anchor='la', spacing=lineHeight)
                anchor = Point(MARGIN, bbx[3]+ lineHeight)
            base.alpha_composite(bg, (self.anchor.x, int(self.anchor.y + (self.pixelsPerChar.y * 3))))
            return base
        return drawSection

    def onBroadcast(self, event: Broadcast):
        if event.event == "DEFINE":
            self.word = event.obj.text
            defs = event.obj.getDefinitions()
            self.drawImg(self.setWordHeader())
            if defs == [] or defs is None:
                print('no definitions retrieved')
            else:
                self.drawImg(self.setDefSection(defs))