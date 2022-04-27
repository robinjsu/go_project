package main

import (
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
)

// TODO: a bit more string clean-up to do

func drawSearchBar(images []imageObj, bounds *image.Rectangle) func(draw.Image) image.Rectangle {
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
			lineHeight := fontFaces["regular"].Metrics().Height.Floor() * 2
			display = formatTextImages(list, 0, MIN_X_SEARCH, MARGIN*2, lineHeight, fontFaces["regular"])
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
