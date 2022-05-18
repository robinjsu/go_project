package main

import (
	"image"
	"image/draw"

	gui "github.com/faiface/gui"
	"golang.org/x/image/font"
)

func drawHeader(images []imageObj, bounds *image.Rectangle) func(draw.Image) image.Rectangle {
	searchBar := func(drw draw.Image) image.Rectangle {
		newR := makeInsetR(*bounds, MARGIN)
		draw.Draw(drw, newR, &image.Uniform{TEAL}, image.Pt(0, 0), draw.Over)
		for _, obj := range images {
			draw.Draw(drw, obj.placement, obj.img, image.Pt(0, 0), draw.Over)
		}
		return newR
	}
	return searchBar
}

func Header(env gui.Env, fontFaces map[string]font.Face, words <-chan string, define chan<- string) {
	var list []string
	var display []imageObj
	for {
		select {
		case lookup := <-words:
			list = splitStr(lookup)
			lineHeight := fontFaces["regular"].Metrics().Height.Floor() * 2
			display = formatTextImages(list, 0, MIN_X_SEARCH, MARGIN*2, lineHeight, fontFaces["regular"])
			env.Draw() <- drawHeader(display, &wordCorner)
			define <- lookup
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}
}
