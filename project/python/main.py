import threading
from Env import Env, Mux
        

def main():
    main_env = Mux()
    open = True
    main_thread = threading.Thread(target=main_env.poll_events)
    main_thread.start()
    while open == True:
        for i in range(10):
            main_env.send(i)

if __name__ == '__main__':
    main()