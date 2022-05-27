import glfw
import OpenGL.GL as gl
from PIL import Image
from typing import NamedTuple, List
import threading

from .Event import *
from .Env import Env
from .utils import Box, Point

'''
TODO: figure out pixels vs screen coordinates issue ==> what is the best way to resize for the given screen/monitor DPI scale?
'''

class Options(NamedTuple):  
    title: str
    width: int
    height: int
    resizable: bool
    maximized: bool

class Window(Env):
    '''
    Window manages the window context and drawing to the interface. It also manages the mouse and keyboard events within the context.
    '''
    options: Options
    win: glfw._GLFWwindow
    win2: glfw._GLFWwindow
    image: Image.Image
    mouseX: float
    mouseY: float
    drawStream: threading.Thread
    xscale: float
    yscale: float

    def __init__(self, options: Options):
        super().__init__(True)
        assert (self.eventChan() is not None) and (self.drawChan() is not None), f'events and draw channels not properly initialized'
        # self.win = glfw.create_window()
        self.options = options
        self.initGLFW()
        x, y = glfw.get_window_content_scale(self.win)
        self.xscale, self.yscale = x, y
        # get image size in pixels
        width, height = glfw.get_framebuffer_size(self.win)
        self.image = Image.new("RGBA", (width, height), (255,255,255,255))
        self.setMousePos(0,0)
        self.handleDrawCommands()

    def setMousePos(self, x: float, y: float):
        '''
        Save mouse position as cursor moves across the window context
        :param x: 
        :param y: 
        '''
        self.mouseX = x
        self.mouseY = y

    def setGLFWCallbacks(self):
        '''
        Set input callbacks for main context. The callback functions put events onto the main event queue,
        which gets broadcast to all sub-Envs created with the Mux.
        '''
        def cursorCallback(win: glfw._GLFWwindow, x: float, y: float):
            self.setMousePos(x* self.xscale ,y*self.yscale)
         
        def mCallback(win: glfw._GLFWwindow, button: int, action: int, mods: int):
            posX, posY = glfw.get_cursor_pos(win)
            mouseEvent = MouseEvent(button, posX*self.xscale, posY*self.yscale, action)
            self.sendEvent(mouseEvent)

        def kbCallback(win: glfw._GLFWwindow, key: int, scancode: int, action: int, mods: int) -> None:
            keyEvent = KeyEvent(key, action)
            self.sendEvent(keyEvent)

        def contentScaleCallback(win: glfw._GLFWwindow, xscale: float, yscale: float):
            self.xscale = xscale
            self.yscale = yscale

        # def framebufferCallback(win: glfw._GLFWwindow, xscale, yscale):
        #     xscale, yscale = glfw.get_window_content_scale(self.win)
        #     print('rescale')
        #     xratio, yratio = xscale // 1, yscale // 1
        #     resized = self.image.resize((int(self.image.width * xratio), int(self.image.height * yratio)))
            # pass

        def dropCallback(win: glfw._GLFWwindow, paths):
            pathDrop = PathDropEvent(paths[0])
            self.sendEvent(pathDrop)


        glfw.set_cursor_pos_callback(self.win, cursorCallback)
        glfw.set_mouse_button_callback(self.win, mCallback)
        glfw.set_key_callback(self.win, kbCallback)
        # glfw.set_framebuffer_size_callback(self.win, framebufferCallback)
        glfw.set_window_content_scale_callback(self.win, contentScaleCallback)
        glfw.set_drop_callback(self.win, dropCallback)
# 
    def initGLFW(self) -> None:
        '''
        Initialze GLFW library and window context. Callback functions for user input are also initialized.
        '''
        if not glfw.init():
            return
        if self.options.resizable == True:
            glfw.window_hint(glfw.RESIZABLE, glfw.TRUE)
        else:
            glfw.window_hint(glfw.RESIZABLE, glfw.FALSE)

        # glfw.window_hint(glfw.MAXIMIZED, glfw.TRUE)
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
        
        glfw.set_input_mode(self.win, glfw.CURSOR, glfw.CURSOR_NORMAL)

        self.setGLFWCallbacks()

        return

    def pollWinEvents(self) -> None:
        '''
        Main loop that listens for events in window
        All callbacks are triggered from here, and should be running in the main thread
        '''
        while not glfw.window_should_close(self.win):
            glfw.wait_events()
        
        self.sendEvent(BroadcastEvent(BroadcastType().CLOSE, None))
        self.eventChan().close()

        # must be called from main thread
        glfw.terminate()
        return 

    def handleDrawCommands(self) -> None:
        '''
        Create thread for listening for drawing commands send from the Mux or another Env.
        '''
        drawLock = threading.RLock()
        def startOpenGLThread(lock: threading.RLock) -> None:
            self.drawLock = threading.RLock()
            lock.acquire()
            glfw.make_context_current(self.win)
            while self.eventChan().closed == False:
                # Render here
                drawFunc = self.drawChan().receive()
                if drawFunc != None:
                    self.renderWindow(drawFunc(self.image))
                    # Swap front and back buffers
                    glfw.swap_buffers(self.win)

            lock.release()
        self.drawStream = threading.Thread(target=startOpenGLThread, args=(drawLock,), name='WindowDrawThread', daemon=True)
        return
    
  
    def renderWindow(self, img: Image.Image) -> None:
        '''
        Main rendering operations with OpenGL, ensures proper display to the user.
        img: the main OpenGL context within the Window object that gets drawn to.
        '''
        if not self.win:
            print("glfw context not created")
            return
        self.image = img
        width, height = glfw.get_framebuffer_size(self.win)
        gl.glViewport(0, 0, width, height)
        gl.glRasterPos2d(-1,1)
        gl.glPixelZoom(1, -1)
        gl.glDrawPixels(
            img.width, 
            img.height, 
            gl.GL_RGBA, 
            gl.GL_UNSIGNED_BYTE, 
            img.tobytes()
        )
        gl.glFlush()
        return
    
    def run(self) -> None:   
        '''
        Begin drawing and event threads
        ''' 
        self.drawStream.start()
        with self.getLock():
            self.getLock().notify_all()
        self.pollWinEvents()
        return