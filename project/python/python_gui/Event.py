import glfw
from typing import NamedTuple, Callable
from PIL.ImageDraw import *

class Event():
    MouseDown:  Callable
    MouseUp:    Callable
    ArrowUp:    Callable
    ArrowDown:  Callable
    ArrowLeft:  Callable
    ArrowRight: Callable

    def MouseDown(self) -> int:
        return int(glfw.PRESS)
    
    def MouseUp(self) -> int:
        return int(glfw.RELEASE)
    
    def ArrowUp(self) -> int: 
        return int(glfw.KEY_UP)

    def ArrowDown(self) -> int:
        return int(glfw.KEY_DOWN)

    def ArrowLeft(self) -> int:
        return int(glfw.KEY_LEFT)
    
    def ArrowRight(self) -> int:
        return int(glfw.KEY_RIGHT)
    
    
class MouseEvent(NamedTuple):
    button: int
    xpos:   int
    ypos:   int
    action: int

class KeyEvent(NamedTuple):
    key:    int
    action: int
