from python_gui import Window as w, Env, Event
from python_gui.utils import Box
from PIL import Image, ImageDraw, ImageFont
import threading

class Display:
    env: Env
    box: Box

    def __init__(self, env: Env.Env, mainBox: tuple):
        self.env = env
        self.box = Box(x0=mainBox[0], y0=mainBox[1], x1=mainBox[2], y1=mainBox[3])
        
    def setBg(self):
        def drawBg(baseImg: Image.Image) -> Image.Image:
            mainPage = Image.new("RGBA", (self.box.x1, self.box.y1))
            drw = ImageDraw.ImageDraw(mainPage)
            drw.rectangle((0,0,800,800), (255,255,255,255), (0,0,255,255), 5)
            drw.rectangle((800,0,1200,900), (0,0,255,255))
            return mainPage
        return drawBg

    def run(self):
        def startThreads():
            self.env.drawChan().send(self.setBg())
            while True:
                event = self.env.eventChan().receive()
                if event == Event.MouseEvent:
                    pass
                elif event == Event.KeyEvent:
                    pass
                print(event)
        threading.Thread(target=startThreads, name="DisplayThread", daemon=True).start()


def drawSomething(baseImg: Image.Image) -> Image.Image:
    im = baseImg.copy()
    drwCtx = ImageDraw.ImageDraw(im)
    drwCtx.rectangle((0,0,500,500), fill=(0,0,255,255))
    fnt = ImageFont.truetype("../../fonts/Karma/Karma-Regular.ttf", 36)
    drwCtx.text((150,200), "Hello, Python PIL App!", font=fnt, fill=(0,0,0,255))
    out = Image.alpha_composite(baseImg, im)
    return out

options = w.Options("Hello Mux!", 1200, 900, False, None)
win = w.Window(options)
mux = Env.Mux(win)
display = Display(mux.addEnv(), win.image.getbbox())

mux.run()
display.run()
win.run()

