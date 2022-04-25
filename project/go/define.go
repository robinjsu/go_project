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
	DEF_MIN_X = 800
	DEF_MIN_Y = 300
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
	// TODO: standardize points
	lineR := image.Rect(DEF_MIN_X, DEF_MIN_Y, DEF_MIN_X, DEF_MIN_Y)
	var definitions [][]imageObj
	y0 := DEF_MIN_Y
	for _, d := range word.Def {
		var defImages []imageObj
		for _, txt := range d.Wrapped {
			img, format := drawText(txt, face)
			x1 := img.Bounds().Dx()
			y0 += ht
			// TODO: standardize points
			lineR = image.Rect(DEF_MIN_X+MARGIN, y0, DEF_MIN_X+MARGIN+x1, y0+ht)
			defImages = append(defImages, imageObj{text: format, img: img, placement: lineR})
		}
		definitions = append(definitions, defImages)
	}
	return definitions
}

func drawDefs(word string, faces map[string]font.Face, images [][]imageObj) func(draw.Image) image.Rectangle {
	newR := makeInsetRect(defCorner, MARGIN)
	headerImg, _ := drawText(word, faces["bold"])
	headerR := makeHeaderLeftR(newR, headerImg, MARGIN)
	y := 0
	display := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, newR, image.White, newR.Min, draw.Src)
		draw.Draw(drw, headerR, headerImg, image.ZP, draw.Over)
		for _, def := range images {
			for _, line := range def {
				// TODO: standardize points
				newPlacement := line.placement.Add(image.Pt(MARGIN*2, y))
				draw.Draw(drw, newPlacement, line.img, line.img.Bounds().Min, draw.Over)
				y += line.text.bounds.Max.Y.Ceil()
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
				defs = Word{Word: "definition unavailable"}
			} else {
				defs.formatDefs(MAXLINE_DEF)
			}
			images := displayDefs(defs, fontFaces["regular"])
			env.Draw() <- drawDefs(defs.Word, fontFaces, images)
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}

}
