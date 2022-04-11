package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/faiface/gui"
	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func parseFont(file string) (*tt.Font, error) {
	ttfFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ttf, err := tt.Parse(ttfFile)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Parsed Font Family: %s - %s\n", ttf.Name(tt.NameIDFontFamily), ttf.Name(tt.NameIDFontSubfamily))

	return ttf, nil
}

func Text(env gui.Env, textFile string) {
	// load ttf file as bytes
	FONTFAMILY := "../fonts/Karma/Karma-Regular.ttf"
	// parse bytes and return a pointer to a Font type object
	fontstyle, err := parseFont(FONTFAMILY)
	if err != nil {
		error.Error(err)
		panic("panic! TTF file not properly loaded")
	}
	// create face, which provides the `glyph mask images``
	face := tt.NewFace(fontstyle, &tt.Options{
		// options... here just font size (0 is 12-point default)
		Size: 0,
	})

	content, err := os.ReadFile(textFile)
	if err != nil {
		error.Error(err)
		panic("panic! text file not properly loaded")
	}

	loadText := func(drw draw.Image) image.Rectangle {
		page := image.Rect(0, 0, 900, 600)
		text := &font.Drawer{
			Dst:  drw,
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(12, 12),
		}

		fmt.Printf("dot location: %v", text.Dot)

		draw.Draw(drw, page, image.White, image.ZP, draw.Src)
		text.DrawBytes(content)
		fmt.Printf("dot location after draw: %v", text.Dot)

		return page
	}

	env.Draw() <- loadText

	for {
		select {
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}

		}
	}
}
