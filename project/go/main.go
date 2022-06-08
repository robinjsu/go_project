package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tts "cloud.google.com/go/texttospeech/apiv1"
	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/faiface/mainthread"
)

// https://pkg.go.dev/golang.org/x/oauth2#example-Config
func getGoogleClient() (*tts.Client, context.Context) {
	ctx := context.Background()
	jsonFile, err := os.ReadFile("tts_client_secret.json")
	if err != nil {
		log.Fatal(FileError{"tts_client_secret.json", err})
	}
	conf, err := google.ConfigFromJSON(jsonFile, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("***Visit the URL for the auth dialog. When redirected to localhost, copy the auth token ('code' parameter the URL)***: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	fmt.Print("\n\nPlease enter the token (code parameter) copied from the redirect URL in the browser: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	ttsClient, err := tts.NewClient(ctx, option.WithTokenSource(oauth2.StaticTokenSource(tok)))
	if err != nil {
		log.Fatal(err)
	}

	return ttsClient, ctx
}

func init() {
	// _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	// if !ok {
	// 	log.Println("no google credentials supplied")
	// }
	_, ok := os.LookupEnv("DICT_API_KEY")
	if !ok {
		log.Println("no dictionary api key supplied")
	}
}

func run() {

	client, ctx := getGoogleClient()
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
	go TextToSpeech(mux.MakeEnv(), load, text, client, ctx)
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
