package main

import (
	"fmt"
	"os"
	"bytes"
)

func basic_chan() {
	// instantiate new channel
	stream := make(chan interface{})
	go func() {
		// ensure channel is closed after sending object
		defer close(stream)
		// syntax for sending object through channel
		stream<- "hello, world!"
	}()
	// syntax for receiving object from channel;
	// automatically blocks until value is received
	val, ok := <-stream
	// use %t for bool format string type
	fmt.Printf("%v : (%t)\n", val, ok)
}

func buffered_chan() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)
	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 10; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}



func main() {
	// basic_chan()
	buffered_chan()
}