package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	MIN_Y = 300
)

func getWord(s string) (WordDef, error) {
	definitions, err := getDef(s)
	if err != nil {
		return WordDef{}, err
	}
	return definitions, nil
}

func displayDefs(words string, face font.Face, bounds image.Rectangle) func(draw.Image) image.Rectangle {
	display := func(drw draw.Image) image.Rectangle {
		newR := image.Rect(bounds.Min.X+MARGIN, bounds.Min.Y+MARGIN, bounds.Max.X-MARGIN, bounds.Max.Y-MARGIN)
		draw.Draw(drw, newR, image.White, newR.Min, draw.Src)
		def := &font.Drawer{
			Src:  image.Black,
			Face: face,
			Dot:  fixed.P(0, face.Metrics().Height.Ceil()),
		}
		def.Dst = image.NewRGBA(newR)
		def.DrawString(words)
		return newR
	}
	return display
}

func Define(env gui.Env, fontFaces map[string]font.Face, word <-chan string) {
	// var defs WordDef
	defCorner := image.Rect(900, 300, 1200, 1200)
	for {
		select {
		case lookup := <-word:
			defs, err := getWord(lookup)
			if err != nil {
				fmt.Println(err)
			}
			if len(defs.Def) == 0 {
				fmt.Println("definition unavailable")
				continue
			} else {
				fmt.Println(defs)
			}
			env.Draw() <- displayDefs(defs.Def[0].Definition, fontFaces["bold"], defCorner)
			fmt.Println(lookup)
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}

}
