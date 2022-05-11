import glfw
from typing import Any, NamedTuple, Callable
from PIL.ImageDraw import *
    
class Event(NamedTuple):
    MouseDown:  int
    MouseUp:    int
    ArrowUp:    int
    ArrowDown:  int
    ArrowLeft:  int
    ArrowRight: int

InputEvent = Event(
    int(glfw.PRESS),
    int(glfw.RELEASE),
    int(glfw.KEY_UP),
    int(glfw.KEY_DOWN),
    int(glfw.KEY_LEFT),
    int(glfw.KEY_RIGHT)
)

class MouseEvent(NamedTuple):
    button: int
    xpos:   int
    ypos:   int
    action: int

class KeyEvent(NamedTuple):
    key:    int
    action: int

class Broadcast(NamedTuple):
    event: Any
    message: str
## can maybe define a few subtypes of Broadcast? or leave open to implementation
