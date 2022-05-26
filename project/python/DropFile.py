from typing import List, Callable
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand

from pyGui import *
from pyGui.utils import *
from const import *

color = Colors()
class DropFile(Env):
    env: Env
    box: Box
    font: ImageFont.ImageFont
    fontSize: int

    def __init__(self, box, fontFile, id=rand.randint(0,100), name=''):
        super().__init__(id=id, threadName=name)
        self.box = box
        self.fontSize = 20
        self.font = loadFont(fontFile, self.fontSize)
    
    def setSplash(self) -> Callable[..., Image.Image]:
        def drawSplash(baseImg: Image.Image) -> Image.Image:
            bg = Image.new("RGBA", self.box.size())
            drw = ImageDraw.ImageDraw(bg)
            drw.rectangle((0,0,self.box.x1, self.box.y1), color.paleBlue)
            drw.text((50,50), "DROP .TXT FILE OVER WINDOW TO START...", color.black, self.font)
            baseImg.alpha_composite(bg)
            return baseImg

        return drawSplash

    def init(self):
        self.drawImg(self.setSplash())
    

    
    

    
    
