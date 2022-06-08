from PIL import ImageColor

'''
--celadon-blue: #0081a7ff;
--maximum-blue-green: #00afb9ff;
--light-yellow: #fdfcdcff;
--peach-puff: #fed9b7ff;
--bittersweet: #f07167ff;
'''

class Colors():
    def __init__(self):
        self.ultra = ImageColor.getrgb('#00072dff')
        self.navy = ImageColor.getrgb('#001c55ff')
        self.celadon = ImageColor.getrgb('#0081a7ff')
        self.lightBlue = ImageColor.getrgb('#0e6ba8ff')
        self.paleBlue = ImageColor.getrgb('#00afb9ff')
        self.paleBlueTransp = ImageColor.getrgb('#a6e1fa88')
        self.black = ImageColor.getrgb('#000000ff')
        self.light = ImageColor.getrgb('#fdfcefff')

MARGIN = 5
MAXLINES = 40
WINDOW_WIDTH = 1200
WINDOW_HEIGHT = 800
TTF = '../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf'
TTF_BOLD = '../../fonts/Anonymous_Pro/AnonymousPro-Bold.ttf'
AUDIO_DIR = './audio'
TEXT_SZ = 28
TITLE_SZ = 36