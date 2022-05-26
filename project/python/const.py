from PIL import ImageColor
from typing import NamedTuple

class Colors():
    def __init__(self):
        self.ultra = ImageColor.getrgb('#00072dff')
        self.navy = ImageColor.getrgb('#001c55ff')
        self.blue = ImageColor.getrgb('#0a2472ff')
        self.lightBlue = ImageColor.getrgb('#0e6ba8ff')
        self.paleBlue = ImageColor.getrgb('#a6e1faff')
        self.paleBlueTransp = ImageColor.getrgb('#a6e1fa88')
        self.black = ImageColor.getrgb('#000000ff')
        self.white = ImageColor.getrgb('#ffffffff')

MARGIN = 5
MAXLINES = 40
WINDOW_WIDTH = 1200
WINDOW_HEIGHT = 900
TTF = '../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf'
TTF_BOLD = '../../fonts/Anonymous_Pro/AnonymousPro-Bold.ttf'
