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

    def __init__(self, mainBox: tuple, id=rand.randint(0,100)):
        super().__init__(id=id)
        self.box = Box(x0=mainBox[0], y0=mainBox[1], x1=mainBox[2], y1=mainBox[3])
        
    def setBg(self) -> Callable[..., Image.Image]:
        def drawBg(baseImg: Image.Image) -> Image.Image:
            mainPage = Image.new("RGBA", (self.box.x1, self.box.y1))
            drw = ImageDraw.ImageDraw(mainPage)
            drw.rectangle((0,0,800,800), (255,255,255,255), color.ultra)
            drw.rectangle((800,0,1200,900), color.navy)

            baseImg.alpha_composite(mainPage)
            return baseImg
        return drawBg
    
    def onMouseClick(self, action):
        pass
    
    def onKeyPress(self, keyPressed):
        pass

    def init(self):
        self.draw.send(self.setBg())