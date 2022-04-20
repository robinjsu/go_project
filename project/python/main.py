import threading
from Env import Env, Mux
        

def main():
    mainEnv = Mux()
    subEnv = mainEnv.makeNewEnv("sub")
    open = True
    mainThread = threading.Thread(target=mainEnv.poll_events)
    subThread = threading.Thread(target=subEnv.poll_events)
    mainThread.start()
    subThread.start()
    mainEnv.send("hello main Env!")
    mainEnv.relayMsg("sub", "hello Sub Env!!")


if __name__ == '__main__':
    main()