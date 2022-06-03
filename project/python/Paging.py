from typing import Callable
from PIL import Image, ImageDraw, ImageFont
import time, threading, random as rand

from pyGui import *
from pyGui.Event import PathDropEvent, BroadcastType, BroadcastEvent, InputType
from pyGui.utils import *
from const import *

colors = Colors()
broadcast = BroadcastType()
input = InputType()
class Paging(Env):
    btnSize: Point
    btnPrev: Box
    btnNext: Box
    anchorX: int
    anchorY: int
    font: ImageFont.ImageFont
    page: int

    def __init__(self, sz: Point, btnArea: Box, fontFile: str, id=rand.randint(0,100), name=''):
        super().__init__(id, name)
        self.btnSize = sz
        self.anchorX = int((btnArea.size()[0] - (self.btnSize.x * 2)) // 2) + btnArea.x0
        self.anchorY = int((btnArea.size()[1] - (self.btnSize.y * 2)) // 2) + btnArea.y0
        self.btnPrev = Box(self.anchorX, self.anchorY, (self.anchorX+self.btnSize.x), (self.anchorY + self.btnSize.y))
        self.btnNext = Box((self.anchorX + self.btnSize.x), self.anchorY, (self.anchorX+(self.btnSize.x * 2)), (self.anchorY+self.btnSize.y))
        self.font = loadFont(fontFile, 14)
        self.page = 0

    def drawButtons(self) -> Callable[..., Image.Image]:
        pSz = self.font.getsize('PREV')
        ppadX = (self.btnSize.x - pSz[0]) // 2
        ppadY = (self.btnSize.y - pSz[1]) // 2
        nSz = self.font.getsize('NEXT')
        npadX = (self.btnSize.x - nSz[0]) // 2
        npadY = (self.btnSize.y - nSz[1]) // 2 

        colorP = colors.blue
        pColor = colors.white
        colorN = colors.white
        nColor = colors.black

        while True:
            def draw(baseImg: Image.Image) -> Image.Image:
                prev = Image.new("RGBA", (self.btnSize.x, self.btnSize.y), colorP)
                next = Image.new("RGBA", (self.btnSize.x, self.btnSize.y), colorN)
                pDrw = ImageDraw.ImageDraw(prev)
                nDrw = ImageDraw.ImageDraw(next)

                pDrw.text((ppadX, ppadY), 'PREV', pColor)
                nDrw.text((npadX, npadY), 'NEXT', nColor)

                baseImg.alpha_composite(prev, (self.btnPrev.x0, self.btnPrev.y0))
                baseImg.alpha_composite(next, (self.btnNext.x0, self.btnNext.y0))

                return baseImg
            
            self.drawImg(draw)

            colorP, colorN = colorN, colorP
            pColor, nColor = nColor, pColor

            time.sleep(0.5)
    

    def onMouseClick(self, event: MouseEvent):
        pt = Point(event.xpos, event.ypos)
        if event.action == input.Press:
            if self.btnPrev.contains(pt):
                self.sendEvent(BroadcastEvent(broadcast.PREV, None))
            elif self.btnNext.contains(pt):
                self.sendEvent(BroadcastEvent(broadcast.NEXT, None))
        
    def init(self):
        threading.Thread(target=self.drawButtons).start()

       
            
           
            




