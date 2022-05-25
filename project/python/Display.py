from typing import Callable
from pyGui import Env, MouseEvent, KeyEvent, Box, Point
import random as rand, threading
from PIL import Image, ImageDraw
from const import Colors
import glfw

color = Colors()
class Display(Env):
    env: Env
    box: Box

    def __init__(self, box: Box, id=rand.randint(0,100), name=''):
        super().__init__(id=id, threadName=name)
        self.box = box
        
    def setBg(self) -> Callable[..., Image.Image]:
        def drawBg(baseImg: Image.Image) -> Image.Image:
            mainPage = Image.new("RGBA", (self.box.x1, self.box.y1))
            drw = ImageDraw.ImageDraw(mainPage)
            drw.rectangle((0,0,self.box.x1*.75,self.box.y1*.75), color.white, color.ultra)
            baseImg.alpha_composite(mainPage)
            return baseImg
        return drawBg
    

    def init(self):
        self.drawImg(self.setBg())