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

from Channel import DrawChan, EventChan
from Event import *
from Env import Env

# from Env import Env


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
    # events: EventChan
    # draw: DrawChan

    drawLock: threading.RLock
    drawStream: threading.Thread

    win: glfw._GLFWwindow
    image: Image.Image
    options: Options
    mouseX: float
    mouseY: float

    def __init__(self, options: Options):
        super().__init__(True)
        self.win = glfw._GLFWwindow()
        self.options = options
        self.mouseX = 0
        self.mouseY = 0
        self.image = Image.new("RGBA", (self.options.width, self.options.height), (255,255,255,255))
        self.initGLFW()

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
        # set window properties using window_hint()
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

    def poll_events(self) -> None:
        while not glfw.window_should_close(self.win):
            glfw.wait_events()

        self.events.close()
        
        # must be called from main thread
        glfw.terminate()
        return 

    def startOpenGLThread(self, lock: threading.RLock) -> None:
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
        print('run window loops')
        self.initDrawThread(self.startOpenGLThread)
        self.events.open()
        print('start draw stream')
        self.drawStream.start()
        print('poll events')
        self.poll_events()
        

def drawSomething(baseImg: Image.Image) -> Image.Image:
    im = baseImg.copy()
    drwCtx = ImageDraw(im)
    drwCtx.rectangle((0,0,500,500), fill=(0,0,255,255))
    fnt = ImageFont.truetype("../../fonts/Karma/Karma-Regular.ttf", 36)
    drwCtx.text((150,200), "Hello, Python PIL App!", font=fnt, fill=(0,0,0,255))
    out = Image.alpha_composite(baseImg, im)
    return out

def drawCommand(q: DrawChan) -> None:
    dfunc = drawSomething
    q.send(dfunc)


# def main():
    
#     # simulating drawing event coming from a component
#     dThread = threading.Thread(target=drawCommand, args=(win.draw,))
#     dThread.start()
#     dThread.join()
#     win.run()
    

    

# if __name__ == '__main__':
#     main()
        



