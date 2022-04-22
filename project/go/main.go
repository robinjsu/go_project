package main

import (
	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

const (
	MAXWIDTH  = 1200
	TEXTWIDTH = 900
	HEIGHT    = 900
	FONTSZ    = 16
	FONT_REG  = "../../fonts/Karma/Karma-Regular.ttf"
	FONT_BOLD = "../../fonts/Karma/Karma-Bold.ttf"
	FONT_H    = 20
	NEWLINE   = byte('\n')
	maxLineW  = 125
)

func run() {
	// create GUI window (not resizable for now)
	window, err := win.New(win.Title("GoTextAide"), win.Size(MAXWIDTH, HEIGHT))
	if err != nil {
		panic(err)
	}

	//  multiplex main window env
	mux, mainEnv := gui.NewMux(window)
	fontFaces := loadFonts(FONT_REG, FONT_BOLD)
	words := make(chan string)
	define := make(chan string)
	// each component is muxed from main, occupying its own thread
	go Display(mux.MakeEnv())
	go Text(mux.MakeEnv(), "./alice.txt", fontFaces, words)
	go Search(mux.MakeEnv(), fontFaces, words, define)
	go Define(mux.MakeEnv(), define)

	for e := range mainEnv.Events() {
		switch e.(type) {
		case win.WiClose:
			close(mainEnv.Draw())
		}
	}

}

func main() {
	mainthread.Run(run)
}
