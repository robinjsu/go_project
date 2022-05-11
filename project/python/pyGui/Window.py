import glfw
import OpenGL.GL as gl
from PIL import Image
from typing import NamedTuple
import threading

from .Event import *
from .Env import Env
from .utils import Box, Point

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

    def __init__(self, options: Options):
        super().__init__(True)
        assert (self.events is not None) and (self.draw is not None), f'events and draw channels not properly initialized'
        self.win = glfw._GLFWwindow()
        self.options = options
        self.image = Image.new("RGBA", (self.options.width, self.options.height), (255,255,255,255))
        self.setMousePos(0,0)
        self.initGLFW()
        self.handleDrawCommands()
    
    def createLock(self):
        self._ready = threading.Condition()

    def setMousePos(self, x, y):
        '''
        Save mouse position as cursor moves across the window context
        '''
        self.mouseX = x
        self.mouseY = y

    def setCallbacks(self):
        '''
        Set input callbacks for main context. The callback functions put events onto the main event queue,
        which gets propagated to all sub-Envs created with the Mux.
        '''
        def cursorCallback(win: glfw._GLFWwindow, x: float, y: float):
            self.setMousePos(x,y)
            mainBox = Box(0,0,800,800)
            point = Point(glfw.get_cursor_pos(win)[0], glfw.get_cursor_pos(win)[1])
            if mainBox.contains(point):
                cursor = glfw.create_standard_cursor(glfw.IBEAM_CURSOR)
                
            else:
                cursor = glfw.create_standard_cursor(glfw.ARROW_CURSOR)
            glfw.set_cursor(win, cursor)
    
        
        def mCallback(win: glfw._GLFWwindow, button: int, action: int, mods: int):
            mouseEvent = MouseEvent(button, glfw.get_cursor_pos(win)[0], glfw.get_cursor_pos(win)[1], action)
            self.events.send(mouseEvent)

        def kbCallback(win: glfw._GLFWwindow, key: int, scancode: int, action: int, mods: int) -> None:
            keyEvent = KeyEvent(key, action)
            self.events.send(keyEvent)

        glfw.set_cursor_pos_callback(self.win, cursorCallback)
        glfw.set_mouse_button_callback(self.win, mCallback)
        glfw.set_key_callback(self.win, kbCallback)

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
        
        glfw.set_input_mode(self.win, glfw.CURSOR, glfw.CURSOR_NORMAL)

        self.setCallbacks()

        return

    def pollWinEvents(self) -> None:
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
        '''
        Create thread for listening for drawing commands send from the Mux or another Env.
        '''
        drawLock = threading.RLock()
        def startOpenGLThread(lock: threading.RLock) -> None:
            self.drawLock = threading.RLock()
            lock.acquire()
            glfw.make_context_current(self.win)
            while self.events.closed == False:
                # Render here
                drawFunc = self.draw.receive()
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
        gl.glViewport(0,0,width, height)
        gl.glRasterPos2d(-1,1)
        gl.glPixelZoom(1, -1)
        gl.glDrawPixels(
            width, 
            height, 
            gl.GL_RGBA, 
            gl.GL_UNSIGNED_BYTE, img.tobytes()
        )
        gl.glFlush()
        return
    
    def run(self) -> None:   
        '''
        Begin drawing and event threads
        ''' 
        with self._ready:
            self._ready.notify_all()
        self.drawStream.start()
        self.pollWinEvents()
        return