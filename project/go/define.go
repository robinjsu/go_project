package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
)

// TODO: load different size fonts

const (
	DEF_MIN_X = 900
	DEF_MIN_Y = 425
)

func getWord(s string) (Word, error) {
	definitions, err := getDef(s)
	if err != nil {
		return Word{}, err
	}
	return definitions, nil
}

func displayDefs(word Word, face font.Face) [][]imageObj {
	y := face.Metrics().Height.Ceil() * 2
	var definitions [][]imageObj
	for _, d := range word.Def {
		var defImages []imageObj
		for j, txt := range d.Wrapped {
			img, format := drawText(txt, face)
			x1 := img.Bounds().Dx()
			lineR := image.Rect(DEF_MIN_X, DEF_MIN_Y+(y*j), x1+DEF_MIN_X, DEF_MIN_Y+(y*(j+1)))
			defImages = append(defImages, imageObj{text: format, img: img, placement: lineR})
		}
		definitions = append(definitions, defImages)
	}
	return definitions
}

func drawDefs(word string, face font.Face, images [][]imageObj) func(draw.Image) image.Rectangle {
	newR := makeInsetRect(defCorner, MARGIN)
	headerImg, _ := drawText(word, face)
	headerR := makeHeaderR(newR, headerImg, MARGIN)
	y := 0
	display := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, newR, image.White, newR.Min, draw.Src)
		draw.Draw(drw, headerR, headerImg, image.Pt(0, 0), draw.Over)
		for _, def := range images {
			for j, line := range def {
				newPlacement := line.placement.Add(image.Pt(0, y))
				draw.Draw(drw, newPlacement, line.img, image.Pt(0, y), draw.Over)
				y = line.text.bounds.Max.Y.Ceil() * (j + 1)
			}
		}

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
				// fmt.Println("definition unavailable")
				defs = Word{Word: "definition unavailable"}
			} else {
				defs.formatDefs()
			}
			images := displayDefs(defs, fontFaces["regular"])
			env.Draw() <- drawDefs(defs.Word, fontFaces["regular"], images)
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}

}
