from pyGui import Window, Options, Mux
import random as rand, time, os
from Text import Text
from Define import Define
from WordList import WordList
from DropFile import DropFile
from Paging import Paging
from Audio import Audio
from const import *
from pyGui.utils import Box, Point
from google_auth_oauthlib import flow
from google.cloud import texttospeech as tts


options: Options
mux: Mux
text: Text
win: Window

rand.seed(time.time())
dispBox = None
textBox = None
defBox = None

# https://cloud.google.com/docs/authentication/end-user
def oauthFlow():
    launch_browser = True
    appflow = flow.InstalledAppFlow.from_client_secrets_file(
    "tts_client_secret.json", scopes=["https://www.googleapis.com/auth/cloud-platform"]
    )
    if launch_browser:
        appflow.run_local_server()
    else:
        appflow.run_console()

    creds = appflow.credentials
    client = tts.TextToSpeechClient(credentials=creds)
    return client 


def setDimensions(window: Window):
    assert window.image != None, 'window and associated drawing image are not initialized'
    x0, y0, x1, y1 = window.image.getbbox()
    textBox = Box(x0, y0, int(x1*.75), int(y1*.90))
    defBox = Box(int(x1*.75), 0, x1, y1)
    display = Box(x0, y0, x1, y1)
    return display, textBox, defBox

def start():
    googleClient = oauthFlow()
    options = Options("PyTextAide", WINDOW_WIDTH, WINDOW_HEIGHT, False, None)
    win = Window(options)
    dispBox, textBox, defBox = setDimensions(win)
    mux = Mux(win)

    mux.addEnv(Text(textBox, id=2, name='TextThread'))
    mux.addEnv(Define(defBox, id=3, name='DefinitionThread'))
    mux.addEnv(WordList(None, id=4, name='WordListThread'))
    mux.addEnv(DropFile(dispBox, TTF_BOLD, id=5, name='PathDropThread'))
    mux.addEnv(Paging(Point(100, 25), Box(0, textBox.y1, textBox.x1, dispBox.y1), TTF_BOLD, id=6, name='PagingThread'))
    mux.addEnv(Audio(Box(textBox.x0, textBox.y1, textBox.x1,  dispBox.y1), id=7, name='AudioThread', googleApiClient=googleClient))

    # mux.run starts up all envs that have been added to it
    mux.run()
    win.run()
    

def main():
    start()
    

if __name__ == '__main__':
    main()