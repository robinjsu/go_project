from typing import List
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math

from pyGui import *
from pyGui.utils import loadFile

class Text(Env):
    bounds: Box
    font: ImageFont.ImageFont
    padding: int
    width = int
    height = int
    pixelsPerLetter = float
    charsPerWidth = int

    def __init__(self, box: tuple):
        self.padding = 5
        self.bounds = Box(box[0], box[1], box[2], box[3])
        self.width = self.bounds.x1 - self.bounds.x0
        self.height = self.bounds.y1 - self.bounds.y0
    
    def setFont(self, ttf):
        self.font = ttf
        self.pixelsPerLetter = self.font.getlength('A')
        self.charsPerWidth = self.width // math.ceil(self.pixelsPerLetter)


    # assume monospaced font for now
    def formatText(self, fileObj, fileSz):
        lines: List[str]

        paddingW = self.width - (self.padding * 2)
        paddingH = self.width - (self.padding * 2)
        paddedBox = ImageOps.pad(Image.new("RGBA", (self.width, self.height)), (paddingW, paddingH))
        buffer = io.BufferedReader(fileObj, fileSz)
        text = buffer.peek()
        text = text.decode('utf-8')
        idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
        line = text[:idx].rstrip(' \n')
        # break line into words each with its own bounding box here
        lines.append(line)
        if len(text) > idx - 1:
            text = text[idx:]
        while line != '':
            idx = self.makeLineBreak(text[:self.charsPerWidth+1]) + 1
            line = text[:idx].rstrip(' \n')
            # break line into words each with its own bounding box here
            lines.append(line)
            if len(text) > idx - 1:
                text = text[idx:]
        return lines

    def makeLineBreak(self, line: str) -> int:
        if '\n' in line:
            return line.find('\n')
        else:
            return line.rfind(' ')

    def renderText(self):
        pass


def makeLineBreak(line: str) -> int:
    if '\n' in line:
        return line.find('\n')
    else:
        return line.rfind(' ')

if __name__== '__main__':
    lines = []
    font = ImageFont.truetype('../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf')
    # fileObj, fileSz = loadFile('alice.txt')
    # buffer = io.BufferedReader(fileObj, fileSz)
    # text = buffer.read()
    # text = text.decode('utf-8')
    # ln = 100
    # idx = makeLineBreak(text[:ln]) + 1
    # # break line into words each with its own bounding box here
    # lines.append(text[:idx])
    # if len(text) > idx - 1:
    #     text = text[idx:]
    # while text != '':
    #     idx = makeLineBreak(text[:ln]) + 1
    #     # break line into words each with its own bounding box here
    #     line = text[:idx]
    #     lines.append(line.rstrip(' \n'))
    #     if len(text) > idx - 1:
    #         text = text[idx:]
        
    # print(lines)
