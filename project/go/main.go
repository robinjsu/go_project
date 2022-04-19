package main

import (
	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

// const (
// 	WIDTH     = 1200
// 	HEIGHT    = 900
// 	FONTSZ    = 14
// 	FONT_REG  = "../../fonts/Karma/Karma-Regular.ttf"
// 	FONT_BOLD = "../../fonts/Karma/Karma-Bold.ttf"
// 	FONT_H    = 14
// 	NEWLINE   = byte('\n')
// 	maxLineW  = 125
// )

func run() {
	// create GUI window (not resizable for now)
	window, err := win.New(win.Title("Text Reader Aide"), win.Size(MAXWIDTH, HEIGHT))
	if err != nil {
		panic(err)
	}

	//  multiplex main window env
	mux, mainEnv := gui.NewMux(window)

	go Display(mux.MakeEnv())
	go Text(mux.MakeEnv(), "./alice.txt")

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
