from cgitb import text
from typing import List
from PIL import Image, ImageDraw, ImageFont, ImageOps
import os, requests, random as rand

from pyGui import *
from pyGui.utils import *
from const import *

fontSize = 20
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
    pixelsPerChar: float
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
        # self.headerFont = loadFont(ttf, int(sz*1.5))
        self.pixelsPerChar = self.font.getlength('A')
        self.charsPerWidth = self.padW // math.ceil(self.pixelsPerChar)

    def setWordHeader(self):
        w = self.word.rstrip(trailing_chars).lstrip(trailing_chars)
        chars = len(w)
        textSz = self.font.getsize(w)
        spacingW = (self.charsPerWidth - chars) // 2
        anchorX = spacingW * self.pixelsPerChar
        anchorY = (50 - textSz[1]) // 2
        anchor = int(self.anchor.x + anchorX), int(self.anchor.y + anchorY)
        textSz = self.font.getsize(w)
        def drawHeader(baseImg: Image.Image) -> Image.Image:
            bg = Image.new("RGBA", (self.padW, 50), color.paleBlue)
            textImg = Image.new("RGBA", textSz, color.paleBlue)
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
    
    def setWordDef(self, partOfSpeech: str, defn: str):
        partFormat = f'[{partOfSpeech}]'
        fmtDef = formatText(defn, self.charsPerWidth)
        anchor = Point(self.anchor.x, int(self.anchor.y + (self.padH * .2)))
        textHt = self.font.getsize(partOfSpeech)[1] + lineHeight
        sectionHt = textHt * len(fmtDef)
        def drawDef(base: Image.Image) -> Image.Image:
            # bg = Image.new("RGBA", (self.padW, self.padH - 50), color.white)
            bg = Image.new("RGBA", (self.padW, sectionHt), color.white)
            drawDef = ImageDraw.ImageDraw(bg)
            lineAnchor = Point(0,0)
            drawDef.text((lineAnchor.x, lineAnchor.y), partFormat, color.black, self.font, anchor='la')
            lineAnchor.add(0, textHt)
            for l in range(len(fmtDef)):
                print(fmtDef[l])
                drawDef.text(
                    (lineAnchor.x, lineAnchor.y),
                    fmtDef[l],
                    color.black,
                    self.font,
                    anchor='la'
                )
                lineAnchor.add(0, textHt)
            base.alpha_composite(bg, (anchor.x, anchor.y))
            return base
        return drawDef

    def setDefSection(self, definitions):
        def drawSection(base: Image.Image) -> Image.Image:
            bg = Image.new("RGBA", (self.padW, int(self.padH * 0.8)))

    def onBroadcast(self, event: Broadcast):
        if event.event == "DEFINE":
            self.word = event.obj.text
            defs = event.obj.getDefinitions()
            self.drawImg(self.setWordHeader())
            if defs == []:
                print('no definitions retrieved')
            else:
                self.drawImg(self.setWordDef(defs[0][0], defs[0][1]))

            
    
