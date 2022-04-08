package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func parseFont(file string) *tt.Font {
	ttfFile, err := os.ReadFile(file)
	if err != nil {
		error.Error(err)
	}
	ttf, err := tt.Parse(ttfFile)
	if err != nil {
		error.Error(err)
	}
	fmt.Printf("Parsed Font Family: %s - %s\n", ttf.Name(tt.NameIDFontFamily), ttf.Name(tt.NameIDFontSubfamily))

	return ttf
}

func run() {

	FONTFAMILY := "../fonts/Karma/Karma-Regular.ttf"
	fontstyle := parseFont(FONTFAMILY)
	face := tt.NewFace(fontstyle, &tt.Options{
		Size: 24,
	})

	window, err := win.New(win.Title("Hello, Window!"), win.Size(640, 480))
	if err != nil {
		panic(err)
	}

	whitescreen := func(drw draw.Image) image.Rectangle {
		r := image.Rect(0, 0, 640, 480)
		hello := &font.Drawer{
			Dst:  drw,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(24, 24),
		}

		draw.Draw(drw, r, image.White, image.ZP, draw.Src)
		hello.DrawString("Hello, Window!")

		return r
	}
	window.Draw() <- whitescreen

	for event := range window.Events() {
		switch event.(type) {
		case win.WiClose:
			close(window.Draw())
		}
	}
}
func main() {
	mainthread.Run(run)
}
