import errno
# import pyglfw.pyglfw as glfw
import glfw
import OpenGL.GL as gl
import OpenGL.GLU as glu
import PIL.Image as pil
from PIL.ImageDraw import *
from typing import NamedTuple, Any, Callable
from queue import Queue

# from Env import Env


class Options(NamedTuple):  
    title: str
    width: int
    height: int
    resizable: bool
    borderless: bool
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

    def __init__(self, options: Options):
        self.options = options


    # def Events(self) -> Queue:
    #     return self.eventsOut

    # def Draw(self) -> Callable[..., pil.Image]:
    #     return self.draw

    def createOpenGLThread(self) -> None:
        if not glfw.init():
            return
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
        
        glfw.make_context_current(self.win)
        

        while not glfw.window_should_close(self.win):
            # Render here, e.g. using pyOpenGL
            png = pil.open("test_app.png")
            
            self.renderWindow(png)

            # Swap front and back buffers
            glfw.swap_buffers(self.win)

            # Poll for and process events
            glfw.poll_events()

        glfw.terminate()
    
    def renderWindow(self, img: pil.Image):
        # assert type(img) == pil.Image, f'provided image is not pil.Image object type, actual type: {type(img)}'
        
        if not self.win:
            print("glfw context not created")
            return
        cpy = img.copy()
        gl.glRasterPos2d(-1,1)
        gl.glPixelZoom(1, -1)
        gl.glDrawPixels(img.width, img.height, gl.GL_RGBA, gl.GL_UNSIGNED_BYTE, cpy.tobytes())



def main():
    options = Options("Hello Python!", 1200, 900, False, None, None)
    win = Window(options)
    win.createOpenGLThread()

if __name__ == '__main__':
    main()
        



