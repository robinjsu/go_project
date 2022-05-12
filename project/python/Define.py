from typing import List
from PIL import Image, ImageDraw, ImageFont, ImageOps
import os, requests

from pyGui import *
from pyGui.utils import *
from const import *

DICT_API_KEY = os.getenv('DICT_API_KEY')
assert DICT_API_KEY is not None, 'no api key provided'
WORDS_URL = 'https://wordsapiv1.p.rapidapi.com/words'

class Define(Env):
    wordList: List

    def onBroadcast(self, event: Broadcast):
        if event.event == "DEFINE":
            word = (event.message).rstrip('.,;:!? ()"\'')
            print(word)
            headers = {
                'X-RapidAPI-Key': DICT_API_KEY,
                'Accept': '*/*',
                'Connection': 'keep-alive'
            }
            r = requests.get(f'{WORDS_URL}/{word}/definitions', headers=headers)
            print(r.json())
    
    
