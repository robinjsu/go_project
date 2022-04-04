from multiprocessing import Process, Pipe, JoinableQueue
import time



# this is equivalent to basic_chan() in ./channels.go/ 
# like an unbuffered channel in Go
# https://docs.python.org/3/library/multiprocessing.html#exchanging-objects-between-processes
def pipe(msg):
    def f(conn):
        # send object to other end of connection
        conn.send(msg)
        conn.close()
    # Pipe returns two connection objects connected by pipe - pipe is bidirectional
    sconn, rconn = Pipe()
    # instantiate Process object (subprocess) with target function that is started with run()
    # default to start a process is "spawn": resources necessary to run object, but nothing else from parent process
    # args to the function
    p = Process(target=f, args=(sconn,))
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
def queue():
    def send(q):
        for i in range(10):
            print(f'sending {i}')
            q.put(i)
        print(f'producer done')
        q.close()
    q = JoinableQueue(4)
    p = Process(target=send, args=(q,))
    p.start()
    # putting parent proc to sleep ensures that the send() function fills the queue
    # before checking that the queue is empty
    time.sleep(1)
    while q.empty() == False:
        print(f'received: {q.get()}')
        # again, ensures that queue is not empty before trying the loop condition
        time.sleep(1)
    p.join()
        

if __name__ == "__main__":
    queue()
