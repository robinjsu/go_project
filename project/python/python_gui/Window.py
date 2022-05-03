from email.mime import base
import glfw
import OpenGL.GL as gl
import OpenGL.GLU as glu
# import PIL.Image as Image
# from PIL.ImageDraw import ImageDraw
from PIL import Image, ImageDraw, ImageFont
from typing import NamedTuple, Any, Callable
import queue as q
import threading

from .Channel import DrawChan, EventChan
from .Event import *
from .Env import Env

class Options(NamedTuple):  
    title: str
    width: int
    height: int
    resizable: bool
    maximized: bool

class Window(Env):
    '''
    Window manages the window context and drawing to the interface
    It also manages the mouse and keyboard events within the context
    '''
    win: glfw._GLFWwindow
    image: Image.Image
    options: Options
    mouseX: float
    mouseY: float
    drawStream: threading.Thread

    def __init__(self, options: Options):
        super().__init__(True)
        assert (self.events is not None) and (self.draw is not None), f'events and draw channels not properly initialized'
        self.win = glfw._GLFWwindow()
        self.options = options
        self.image = Image.new("RGBA", (self.options.width, self.options.height), (255,255,255,255))
        self.setMousePos(0,0)
        self.initGLFW()
        self.handleDrawCommands()

    def setMousePos(self, x, y):
        self.mouseX = x
        self.mouseY = y

    def setCallbacks(self):
        mouseEvent, kbEvent = MouseEvent(), KbEvent()
        mouseEvent.setCallback(self)
        kbEvent.setCallback(self)

    def initGLFW(self) -> None:
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

        return

    def pollEvents(self) -> None:
        '''
        Main loop that listens for events in window
        All callbacks are triggered from here, and should be running in the main thread
        '''
        while not glfw.window_should_close(self.win):
            glfw.wait_events()
        self.events.close()

        # must be called from main thread
        glfw.terminate()
        return 

    def handleDrawCommands(self) -> None:
        drawLock = threading.RLock()
        def startOpenGLThread(lock: threading.RLock) -> None:
            self.drawLock = threading.RLock()
            lock.acquire()
            glfw.make_context_current(self.win)
            while self.events.closed == False:
                # Render here
                drawFunc = self.draw.receive()
                # print(f'draw event received: {drawFunc}')
                if drawFunc != None:
                    self.renderWindow(drawFunc(self.image))
                    # Swap front and back buffers
                    glfw.swap_buffers(self.win)

            lock.release()
        self.drawStream = threading.Thread(target=startOpenGLThread, args=(drawLock,), name='WindowDrawThread')
        return
  
    def renderWindow(self, img: Image.Image) -> None:
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
        return
    
    def run(self) -> None:         
        self.drawStream.start()
        self.pollEvents()
        return