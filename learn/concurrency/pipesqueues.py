from multiprocessing import Process, Pipe, JoinableQueue
from pickletools import read_unicodestring4
from queue import Queue
import threading
import time
import asyncio

# this is equivalent to basic_chan() in ./channels.go/ 
# like an unbuffered channel in Go
# https://docs.python.org/3/library/multiprocessing.html#exchanging-objects-between-processes
def sender(conn, msg):
        # send object to other end of connection
        conn.send(msg)
        conn.close()

def pipe(msg):
    # Pipe returns two connection objects connected by pipe - pipe is bidirectional
    rconn, sconn = Pipe()
    # instantiate Process object (subprocess) with target function that is started with run()
    # default to start a process is "spawn": resources necessary to run object, but nothing else from parent process
    # args to the function
    p = Process(target=sender, args=(sconn,msg))
    # spin new process
    p.start()
    # receive object from other end of the Pipe connection. blocks until it receives something,
    # or raises error if other end is closed
    print(rconn.recv())
    # join() blocks until the process whose join method is terminated
    # optional argument for a timeout
    p.join()


# https://docs.python.org/3/library/queue.html
def exampleQueue():
    q = Queue()

    def worker():
        while True:
            item = q.get()
            print(f'Working on {item}')
            print(f'Finished {item}')
            q.task_done()

    # Turn-on the worker thread.
    threading.Thread(target=worker, daemon=True).start()

    # Send thirty task requests to the worker.
    for item in range(30):
        q.put(item)

    # Block until all tasks are done.
    q.join()
    print('All work completed')  


# if it's unknown how many items will be put in a queue, how to structure the while loop on #42? 
# or is there a different primitive that is better suited for this task?
def send(q):
        for i in range(10):
            print(f'sending {i}')
            q.put(i)
        print(f'producer done')
    

def receive(q):
    while True:
        item = q.get()
        print(f'received: {item}')
        q.task_done()
        # if q.empty == True:
        #     break

def queue():
    q = JoinableQueue(4)
    p1 = Process(target=send, args=(q,))
    p2 = Process(target=receive, args=(q,))
    p1.start()
    p2.start() 
    # q.join()
    # print('queue joined')
    q.close()
    # p2.join()
    p2.join()
    p1.join()
    # q.close()


def send_t(q: Queue):
        for i in range(10):
            print(f'sending {i}')
            q.put(i)
        q.join()
        print(f'producer done')
        return
    
def receive_t(q: Queue, e: threading.Event):
    while True:
        item = q.get()
        print(f'received: {item}')
        q.task_done()
        if e.is_set():
            return    

def queue_threading():
    q = Queue(4)
    done = threading.Event()
    t1 = threading.Thread(target=send_t, args=(q,))
    t2 = threading.Thread(target=receive_t, args=(q,done))
    t1.start()
    t2.start()
    t1.join()
    # print('put -1')
    for i in range(10,20):
        q.put(i)
    done.set()
    

if __name__ == "__main__":
    # queue()
    # pipe("hello, robin!")
    # exampleQueue()
    queue_threading()