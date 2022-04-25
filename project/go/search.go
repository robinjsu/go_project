package main

import (
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
)

const (
	MIN_X = 800
)

// TODO: a bit more string clean-up to do
func displayWords(wordList []string, face font.Face) []imageObj {
	var images []imageObj
	y := face.Metrics().Height.Ceil() * 2

	for i, w := range wordList {
		img, format := drawText(w, face)
		x1 := img.Bounds().Dx()
		fontR := image.Rect(MIN_X, (y * i), (x1 + MIN_X), (y * (i + 1)))
		images = append(images, imageObj{text: format, img: img, placement: fontR})
	}
	return images
}

func drawSearchBar(images []imageObj, bounds *image.Rectangle) func(draw.Image) image.Rectangle {
	searchBar := func(drw draw.Image) image.Rectangle {
		newR := *bounds
		draw.Draw(drw, newR, &image.Uniform{TEAL}, image.Pt(0, 0), draw.Over)
		for _, obj := range images {
			draw.Draw(drw, obj.placement, obj.img, image.Pt(0, 0), draw.Over)
		}
		return newR
	}
	return searchBar
}

func highlightWord(images []imageObj, p image.Point, drawDst image.Rectangle, define chan<- string) (func(draw.Image) image.Rectangle, string) {
	var target image.Rectangle
	var lookup string
	for _, img := range images {
		if p.In(img.placement) {
			lookup = img.text.txt
			target = img.placement
		}
	}
	highlight := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, drawDst, image.Transparent, image.ZP, draw.Over)
		draw.Draw(drw, target, &image.Uniform{LIGHT_GRAY}, image.ZP, draw.Over)
		return drawDst
	}
	return highlight, lookup
}

func Search(env gui.Env, fontFaces map[string]font.Face, words <-chan string, define chan<- string) {
	var list []string
	var display []imageObj
	for {
		select {
		case lookup := <-words:
			list = splitStr(lookup)
			display = displayWords(list, fontFaces["regular"])
			env.Draw() <- drawSearchBar(display, &wordCorner)
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				if image.Pt(e.X, e.Y).In(wordCorner) {
					env.Draw() <- drawSearchBar(display, &wordCorner)
					highlight, target := highlightWord(display, image.Pt(e.X, e.Y), wordCorner, define)
					define <- target
					env.Draw() <- highlight
				}
				// case win.MoUp:
			}
		}
	}
}
