package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
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

func drawInsetRect(r image.Rectangle, margin int) image.Rectangle {
	x0 := defCorner.Min.X + MARGIN
	y0 := defCorner.Min.Y + MARGIN
	sz := defCorner.Size().Sub(image.Pt(MARGIN, MARGIN))

	return image.Rect(x0, y0, x0+sz.X, y0+sz.Y)
}

func displayDefs(words string, face font.Face) func(draw.Image) image.Rectangle {
	display := func(drw draw.Image) image.Rectangle {
		newR := drawInsetRect(defCorner, MARGIN)
		draw.Draw(drw, newR, image.White, newR.Min, draw.Src)
		textImg, _ := drawText(words, face)
		newBounds := textImg.Bounds().Add(image.Pt(defCorner.Min.X, defCorner.Min.Y))
		draw.Draw(drw, newBounds, textImg, textImg.Bounds().Min, draw.Over)
		return newR
	}
	return display
}

func Define(env gui.Env, fontFaces map[string]font.Face, word <-chan string) {
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
			env.Draw() <- displayDefs(defs.Def[0].Definition, fontFaces["regular"])
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}

}
