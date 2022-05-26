from pyGui import Window, Options, Mux
import random as rand, time
from Display import Display
from Text import Text
from Define import Define
from WordList import WordList
from DropFile import DropFile
from const import *
from pyGui.utils import Box 

options: Options
mux: Mux
display: Display
text: Text
win: Window

# TODO: where to put this?
rand.seed(time.time())


dispBox = None
textBox = None
defBox = None
paging = None

def setDimensions(window: Window):
    assert window.image != None, 'window and associated drawing image are not initialized'
    x0, y0, x1, y1 = window.image.getbbox()
    textBox = Box(x0, y0, int(x1*.75), int(y1*.90))
    defBox = Box(int(x1*.75), 0, x1, y1)
    display = Box(x0, y0, x1, y1)
    return display, textBox, defBox

def start():
    options = Options("PyTextAide", WINDOW_WIDTH, WINDOW_HEIGHT, False, None)
    win = Window(options)
    dispBox, textBox, defBox = setDimensions(win)

    mux = Mux(win)
    # mux.addEnv(Display(textBox, id=1, name="DisplayThread"))
    mux.addEnv(Text(textBox, id=2, name="TextThread"))
    mux.addEnv(Define(defBox, id=3, name="DefinitionThread"))
    mux.addEnv(WordList(None, id=4, name="WordListThread"))
    mux.addEnv(DropFile(dispBox, TTF_BOLD, id=5, name="PathDropThread"))

    # mux.run() will start up all envs that have been added to it
    mux.run()
    win.run()
    

def main():
    start()
    

if __name__ == '__main__':
    main()