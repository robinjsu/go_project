package main

import (
	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

func run() {
	// create GUI window (not resizable for now)
	window, err := win.New(win.Title("GoTextAide"), win.Size(MAXWIDTH, MAXHEIGHT))
	if err != nil {
		panic(err)
	}

	//  multiplex main window env
	mux, mainEnv := gui.NewMux(window)
	fontFaces, err := loadFonts(FONTSZ, FONT_REG, FONT_BOLD)
	if err != nil {
		panic(err)
	}

	// create channels for comms between goroutines
	words := make(chan string)
	define := make(chan string)
	save := make(chan Word)
	filepath := make(chan string)
	load := make(chan bool)

	// each component is muxed from main, occupying its own thread
	go Display(mux.MakeEnv(), load)
	go Text(mux.MakeEnv(), "./alice.txt", fontFaces, words, filepath, load)
	go Header(mux.MakeEnv(), fontFaces, words, define)
	go Define(mux.MakeEnv(), fontFaces, define, save)
	go WordList(mux.MakeEnv(), save, "test")
	go Load(mux.MakeEnv(), fontFaces, filepath)

	// main application loop
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
