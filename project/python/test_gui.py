from pyGui import Window, Options, Env, Mux
# from pyGui import Box
from PIL import Image, ImageDraw, ImageFont
from Display import Display

# def drawSomething(baseImg: Image.Image) -> Image.Image:
#     im = baseImg.copy()
#     drwCtx = ImageDraw.ImageDraw(im)
#     drwCtx.rectangle((0,0,500,500), fill=(0,0,255,255))
#     fnt = ImageFont.truetype("../../fonts/Karma/Karma-Regular.ttf", 36)
#     drwCtx.text((150,200), "Hello, Python PIL App!", font=fnt, fill=(0,0,0,255))
#     out = Image.alpha_composite(baseImg, im)
#     return out

options: Options
mux: Mux
display: Display
win: Window

def init():
    options = Options("Hello Mux!", 1200, 900, False, None)
    win = Window(options)
    mux = Mux(win)
    display = Display(mux.addEnv(), win.image.getbbox())

    return win, display

def main():
    win, display = init()
    display.run()
    win.run()


if __name__ == '__main__':
    main()
