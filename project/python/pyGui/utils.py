from typing import NamedTuple, Tuple
from PIL import ImageFont
import io, os

class Point(NamedTuple):
    x: int
    y: int

class Box(NamedTuple):
    x0: int
    y0: int
    x1: int
    y1: int

    def contains(self, p: Point) -> bool:
        '''
        Returns whether the given point is within the bounds of the box. Exclusive of the Box.x1 and Box.y1 values
        :param p: Point object to test
        '''
        return p.x >= self.x0 and p.x < self.x1 and p.y >= self.y0 and p.y < self.y1

def loadFont(filepath) -> ImageFont.ImageFont:
    return ImageFont.truetype(filepath, size=12)

def loadFile(filepath) -> Tuple[io.FileIO, int]:
    stats = os.stat(filepath)
    textObj = open(filepath, 'rb')
    return textObj, stats.st_size