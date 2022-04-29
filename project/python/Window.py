import errno
# import pyglfw.pyglfw as glfw
import glfw
import OpenGL.GL as gl
import OpenGL.GLU as glu
import PIL.Image as pil
from PIL.ImageDraw import *
from typing import NamedTuple, Any, Callable
from queue import Queue

from Event import *

# from Env import Env


class Options(NamedTuple):  
    title: str
    width: int
    height: int
    resizable: bool
    maximized: bool
    
    # def __init__(self, title="", width=640, height=480, resizable=False, brdlss=False, maximzd=False):
    #     self.title = title
    #     self.width = width
    #     self.height = height
    #     self.resizable = resizable
    #     self.borderless = brdlss
    #     self.maximized = maximzd


class Window():
    '''
    Window manages the window context and drawing to the interface
    It also manages the mouse and keyboard events within the context
    '''
    # eventsOut: Queue
    # eventsIn: Any
    draw: Callable[..., pil.Image]

    win: glfw._GLFWwindow
    image: pil.Image
    options: Options
    mouseX: float
    mouseY: float

    def __init__(self, options: Options):
        self.options = options
        self.mouseX = 0
        self.mouseY = 0
        self.initGLFW()


    # def Events(self) -> Queue:
    #     return self.eventsOut

    # def Draw(self) -> Callable[..., pil.Image]:
    #     return self.draw

    def setCallbacks(self):
        glfw.set_key_callback(self.win, kbCallback)
        glfw.set_mouse_button_callback(self.win, mCallback)
        
        def cursorCallback(win, x, y):
            self.mouseX = x
            self.mouseY = y
        glfw.set_cursor_pos_callback(self.win, cursorCallback)

    def initGLFW(self):
        if not glfw.init():
            return
        
        if self.options.resizable == True:
            glfw.window_hint(glfw.RESIZABLE, glfw.TRUE)
        else:
            glfw.window_hint(glfw.RESIZABLE, glfw.FALSE)
        glfw.window_hint(glfw.MAXIMIZED, glfw.TRUE)

        self.win = glfw.create_window(
            self.options.width, 
            self.options.height, 
            self.options.title, 
            None, 
            None
        )
        if not self.win:
            glfw.terminate()
            return
        
        self.setCallbacks()
        

    def createOpenGLThread(self) -> None:
        glfw.make_context_current(self.win)
        while not glfw.window_should_close(self.win):
            # Render here, e.g. using pyOpenGL
            png = pil.open("test_app.png")
            self.renderWindow(png)
            # Swap front and back buffers
            glfw.swap_buffers(self.win)
            # puts thread to sleep, wakes upon receipt of new event
            glfw.wait_events()
            print(self.mouseX, self.mouseY)

        glfw.terminate()
    
    def renderWindow(self, img: pil.Image):
        if not self.win:
            print("glfw context not created")
            return
        cpy = img.copy()
        x, y, width, height = cpy.getbbox()
        gl.glViewport(x, y, width, height)
        gl.glRasterPos2d(-1,1)
        gl.glPixelZoom(1, -1)
        gl.glDrawPixels(
            img.width, 
            img.height, 
            gl.GL_RGBA, 
            gl.GL_UNSIGNED_BYTE, cpy.tobytes()
        )

        gl.glFlush()



def main():
    options = Options("Hello Python!", 1200, 900, False, None)
    win = Window(options)
    win.createOpenGLThread()

if __name__ == '__main__':
    main()
        



