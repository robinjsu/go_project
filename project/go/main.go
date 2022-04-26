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
	fontFaces := loadFonts(FONTSZ, FONT_REG, FONT_BOLD)

	// create channels for comms between goroutines
	words := make(chan string)
	define := make(chan string)
	save := make(chan Word)
	// next := make(chan string)
	// prev := make(chan string)

	// each component is muxed from main, occupying its own thread
	go Display(mux.MakeEnv())
	go Text(mux.MakeEnv(), "./alice.txt", fontFaces, words)
	// go PagingBtns(mux.MakeEnv(), next, prev, fontFaces)
	go Search(mux.MakeEnv(), fontFaces, words, define)
	go Define(mux.MakeEnv(), fontFaces, define, save)
	go WordList(mux.MakeEnv(), save, "Edith Wharton")

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
