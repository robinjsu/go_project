package main

import (
	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

func run() {
	// create GUI window (not resizable for now)
	window, err := win.New(win.Title("Text Reader Aide"), win.Size(900, 600))
	if err != nil {
		panic(err)
	}

	//  multiplex main window env
	mux, mainEnv := gui.NewMux(window)

	go Display(mux.MakeEnv())
	go Text(mux.MakeEnv(), "./test.txt")

	for evnt := range mainEnv.Events() {
		switch evnt.(type) {
		case win.WiClose:
			close(window.Draw())
		}
	}
}

func main() {
	mainthread.Run(run)
}
