from pyGui import Window, Options, Mux
import random as rand, time
from Display import Display
from Text import Text
from const import * 

options: Options
mux: Mux
display: Display
text: Text
win: Window

# TODO: where to put this?
rand.seed(time.time())

def start():
    options = Options("PyTextAide", 1200, 900, False, None)
    win = Window(options)
    mux = Mux(win)
    display = mux.addEnv(Display(win.image.getbbox(), id=1))
    text = mux.addEnv(Text(TEXTBOX, id=2))
    display.run("DisplayThread")
    text.run("TextThread")
    mux.run()
    win.run()

def main():
    start()
    

if __name__ == '__main__':
    main()