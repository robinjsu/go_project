package main

import (
	"log"
	"os"

	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"

	"github.com/faiface/mainthread"
)

func init() {
	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		log.Println("no google credentials supplied")
	}
	_, ok = os.LookupEnv("DICT_API_KEY")
	if !ok {
		log.Println("no dictionary api key supplied")
	}
}

func run() {
	// create GUI window (not resizable for now)
	window, err := win.New(win.Title("GoTextAide"), win.Size(MAXWIDTH, MAXHEIGHT))
	if err != nil {
		log.Fatal(err)
	}

	//  multiplex main window env
	mux, mainEnv := gui.NewMux(window)
	fontFaces, err := loadFonts(FONTSZ, FONT_REG, FONT_BOLD)
	if err != nil {
		log.Fatal(err)
	}
	largeFont, err := loadFonts(FONTSZ*2, FONT_BOLD)
	if err != nil {
		log.Fatal(err)
	}

	// create channels for comms between goroutines
	words := make(chan string)
	define := make(chan string)
	save := make(chan Word)
	filepath := make(chan string)
	load := make(chan bool)
	page := make(chan string)
	text := make(chan [][]string)

	// each component is muxed from main, occupying its own thread
	go Display(mux.MakeEnv())
	go Text(mux.MakeEnv(), copyFonts(fontFaces), words, filepath, load, page, text)
	go Header(mux.MakeEnv(), copyFonts(fontFaces), words, define)
	go Define(mux.MakeEnv(), copyFonts(fontFaces), define, save)
	go WordList(mux.MakeEnv(), save)
	go Load(mux.MakeEnv(), largeFont["bold"], filepath)
	go TextToSpeech(mux.MakeEnv(), load, text)
	go PagingBtns(mux.MakeEnv(), page, fontFaces, load)

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
