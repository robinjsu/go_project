package main

import (
	"fmt"
	"image"
	"image/draw"
	"strings"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
)

func getWord(s string) (Word, error) {
	definitions, err := getDef(s)
	if err != nil {
		return Word{}, err
	}
	return definitions, nil
}

func displayDefs(word Word, face font.Face) [][]imageObj {
	ht := face.Metrics().Height.Floor()
	var definitions [][]imageObj
	y0 := MIN_Y_DEF
	for _, d := range word.Def {
		var defImages []imageObj
		defImages = formatTextImages(d.Wrapped, y0, MIN_X_DEF, MARGIN*3, ht, face)
		y0 += len(d.Wrapped) * ht
		definitions = append(definitions, defImages)
	}
	return definitions
}

func drawDefs(word string, faces map[string]font.Face, images [][]imageObj, bounds *image.Rectangle) func(draw.Image) image.Rectangle {
	newR := makeInsetR(*bounds, MARGIN)
	y := 0
	display := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, newR, image.White, newR.Min, draw.Src)
		for _, def := range images {
			for _, line := range def {
				newPlacement := line.placement
				draw.Draw(drw, newPlacement, line.img, line.img.Bounds().Min, draw.Over)
				y += line.text.bounds.Max.Y.Ceil()
			}
		}

		return newR
	}
	return display
}

func Define(env gui.Env, fontFaces map[string]font.Face, define <-chan string, save chan<- Word) {
	for {
		select {
		case word := <-define:
			word = strings.Trim(word, " ,.!?';:“”’\"()")
			defs, err := getWord(word)
			if err != nil {
				fmt.Println(err)
			}
			if len(defs.Def) == 0 {
				defs = Word{Word: "definition unavailable"}
			} else {
				defs.formatDefs(fontFaces["regular"], &defCorner)
				save <- defs
			}
			images := displayDefs(defs, fontFaces["regular"])
			env.Draw() <- drawDefs(defs.Word, fontFaces, images, &defCorner)
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}

}
