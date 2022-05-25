from collections import namedtuple
from tokenize import Name
import glfw
from typing import Any, NamedTuple
from PIL.ImageDraw import *

class InputType(NamedTuple):
    Press = int(glfw.PRESS)
    Release = int(glfw.RELEASE)
    ArrowUp = int(glfw.KEY_UP)
    ArrowDown = int(glfw.KEY_DOWN)
    ArrowLeft = int(glfw.KEY_LEFT)
    ArrowRight = int(glfw.KEY_RIGHT)

class BroadcastType(NamedTuple):
    CLOSE = -1
    DEFINE = 0
    SAVE = 1

class MouseEvent(NamedTuple):
    button: int
    xpos:   int
    ypos:   int
    action: int

class KeyEvent(NamedTuple):
    key:    int
    action: int

class BroadcastEvent(NamedTuple):
    event: BroadcastType
    obj: Any
## TODO: can maybe define a few subtypes of Broadcast? or leave open to implementation
