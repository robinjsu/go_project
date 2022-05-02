import glfw
import OpenGL.GL as gl
import OpenGL.GLU as glu
import PIL.Image as pil
from PIL.ImageDraw import *
from typing import NamedTuple, Any, Callable
# from Window import Window

class Event:
    def setCallback(self, window):
        pass

class MouseEvent(Event):
    def setCallback(self, window):
        def cursorCallback(win: glfw._GLFWwindow, x: float, y: float):
            window.setMousePos(x,y)
        glfw.set_cursor_pos_callback(window.win, cursorCallback)

        def mCallback(win: glfw._GLFWwindow, button: int, action: int, mods: int):
            if action == int(glfw.PRESS):
                print(f'cursor position: {glfw.get_cursor_pos(win)}')
            elif action == int(glfw.RELEASE):
                print("mouse release")
        glfw.set_mouse_button_callback(window.win, mCallback)

class KbEvent(Event):
    def setCallback(self, window):
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

        glfw.set_key_callback(window.win, kbCallback)






    