from typing import List, Callable
from PIL import Image, ImageDraw, ImageFont, ImageOps
import io, os, math, random as rand
from google.cloud import texttospeech as tts

from pyGui import *
from pyGui.Event import PathDropEvent
from pyGui.Event import BroadcastType
from pyGui.utils import *
from const import *


colors = Colors()
input = Event.InputType()
bdcast = Event.BroadcastType()
audioDir = './audio'
class Audio(Env):
    bounds: Box
    icons: List[Image.Image]
    iconsBounds: List[Box]
    anchor: Point
    ttsClient: tts.TextToSpeechClient
    

    def __init__(self, box: Box, id=0, name=''):
        super().__init__(id=id, threadName=name)
        self.bounds = box.move(Point(MARGIN, 0))


    def loadIcons(self, *icons) -> List[Image.Image]:
        playbackIcons = []
        for icon in icons:
            img = Image.open(icon)
            img = img.convert("RGBA")
            playbackIcons.append(img)
        
        return playbackIcons
    
    def getIconBounds(self, anchor: Point, icon: Image.Image) -> Box:
        return Box(anchor.x, anchor.y, anchor.x + icon.size[0], anchor.y + icon.size[1])

    def setIcons(self, iconSize: tuple) -> List[Box]:
        iconsPos = []
        _, height = self.bounds.size()
        pad = int((height - iconSize[1]) / 2)
        self.anchor = Point(self.bounds.x0, self.bounds.y0 + pad)
        anch = self.anchor.copy()

        for ic in self.icons:
            iconBox = self.getIconBounds(anch, ic)
            iconsPos.append(iconBox)
            anch.add(iconSize[0], 0)
        
        return iconsPos
    
    def drawIcons(self) -> Callable[..., Image.Image]:
        def draw(drw: Image.Image) -> Image.Image:
            for ic in range(len(self.icons)):
                iconBg = Image.new("RGBA", (self.iconsBounds[ic].size()[0], self.iconsBounds[ic].size()[1]), colors.black)
                drw.paste(iconBg, (self.iconsBounds[ic].x0, self.iconsBounds[ic].y0), self.icons[ic])
            return drw
        
        return draw
    
    def loadTTS(self, textBody, name):
        voice = tts.VoiceSelectionParams(
                    language_code="en-US", ssml_gender=tts.SsmlVoiceGender.NEUTRAL)
        ttsRequest = tts.SynthesizeSpeechRequest(
            input=tts.SynthesisInput(text=textBody),
            voice=voice,
            audio_config=tts.AudioConfig(audio_encoding=tts.AudioEncoding.MP3)
        )
        ttsResponse = self.ttsClient.synthesize_speech(ttsRequest)

        with open(name, 'wb') as out:
            out.write(ttsResponse.audio_content)


    # def onBroadcast(self, event: BroadcastEvent):
    #     if event.event == bdcast.TEXT:
    #         assert self.ttsClient is not None
    #         linesPerPage = (event.obj)['maxLines']
    #         textBody = (event.obj)['textBody']
    #         pages = math.ceil(len(textBody) / linesPerPage)
            
    #         for p in range(pages):
    #             textPg = ' '.join(textBody[:linesPerPage])
    #             self.loadTTS(textPg, f'{audioDir}/pg-{p}.mp3')
    #             textBody = textBody[linesPerPage:]


    def onMouseClick(self, event: MouseEvent):
        pt = Point(event.xpos, event.ypos)
        if event.action == input.Press:
            if self.iconsBounds[0].contains(pt):
                print('play')
            elif self.iconsBounds[1].contains(pt):
                print('pause')
            elif self.iconsBounds[2].contains(pt):
                print('prev')
            elif self.iconsBounds[3].contains(pt):
                print('next')

    
    def init(self):
        self.icons = self.loadIcons('./images/play.png', './images/pause.png', './images/prev.png', './images/next.png')
        self.iconsBounds = self.setIcons(self.icons[0].size)
        self.drawImg(self.drawIcons())
        self.ttsClient = tts.TextToSpeechClient()