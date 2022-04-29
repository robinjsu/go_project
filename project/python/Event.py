import glfw
import OpenGL.GL as gl
import OpenGL.GLU as glu
import PIL.Image as pil
from PIL.ImageDraw import *
from typing import NamedTuple, Any, Callable
from Window import Window

class Event:
    def setCallback(func: Callable[...,Any]):
        pass


def kbCallback(win: glfw._GLFWwindow, key: int, scancode: int, action: int, mods: int) -> None:
    if key == int(glfw.KEY_DOWN):
        print("kbdown")
    elif key == int(glfw.KEY_UP):
        print("kbup")
    elif key == int(glfw.KEY_LEFT):
        print("kbleft")
    elif key == int(glfw.KEY_RIGHT):
        print("kbright")
    else:
        print(glfw.get_key_name(key, scancode))


def mCallback(win: glfw._GLFWwindow, button: int, action: int, mods: int):
    if action == int(glfw.PRESS):
        print(f'cursor position: {glfw.get_cursor_pos(win)}')
    elif action == int(glfw.RELEASE):
        print("mouse release")


    