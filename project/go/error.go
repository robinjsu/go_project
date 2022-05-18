package main

import "fmt"

type DrawError struct {
	component string
}

func (d DrawError) Error() string {
	return fmt.Sprintf("error when drawing %s", d.component)
}

type HttpError struct {
	method string
	url    string
	err    error
}

func (h HttpError) Error() string {
	return fmt.Sprintf("error in http %s request: %s ", h.method, h.err)
}

type FileError struct {
	filename string
	err      error
}

func (fs FileError) Error() string {
	return fmt.Sprintf("error in opening %s: %s", fs.filename, fs.err)
}

type ReadError struct {
	err error
}

func (r ReadError) Error() string {
	return fmt.Sprintf("error in reading buffer or stream: %s", r.err)
}
