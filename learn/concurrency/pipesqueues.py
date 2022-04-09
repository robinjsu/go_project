from multiprocessing import Process, Pipe, JoinableQueue
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

def send_t(q):
        for i in range(10):
            print(f'sending {i}')
            q.put(i)
        print(f'producer done')
        q.put(-1)
        q.join()
        print('queue done')
    
def receive_t(q):
    while True:
        item = q.get()
        print(f'received: {item}')
        q.task_done()
        if item == -1:
            break
        
def queue_threading():
    q = JoinableQueue(4)
    t1 = threading.Thread(target=send_t, args=(q,))
    t2 = threading.Thread(target=receive_t, args=(q,))
    t2.start()
    t1.start()
    # q.join()
    # print('queue finished')
    t2.join()
    t1.join()


if __name__ == "__main__":
    # queue()
    # pipe("hello, robin!")
    queue_threading()