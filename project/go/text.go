package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
)

var (
	currentPage = 0
)

func drawTextLines(images []imageObj, face font.Face, bounds *image.Rectangle) func(drw draw.Image) image.Rectangle {
	load := func(drw draw.Image) image.Rectangle {
		// coordinates refer to the destination image's coordinate space
		// TODO: standardize points
		page := *bounds
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for _, obj := range images {
			draw.Draw(drw, obj.placement, obj.img, image.Pt(0, 0), draw.Over)
		}
		return page
	}
	return load
}

func highlightLine(face font.Face, images []imageObj, p image.Point, words chan<- string) func(draw.Image) image.Rectangle {
	var line image.Rectangle
	load := func(drw draw.Image) image.Rectangle {
		var txt string
		for _, ln := range images {
			rct := ln.placement.Bounds()
			fmt.Println(rct)
			if p.In(rct) {
				txt = ln.text.txt
				draw.Draw(drw, rct, &image.Uniform{HIGHLIGHT_GRAY}, image.ZP, draw.Over)
			}
		}
		// send words to Search component
		words <- txt
		return line
	}
	return load
}

func Text(env gui.Env, textFile string, fontFaces map[string]font.Face, words chan<- string) {
	cont := NewContent()
	_, err := cont.parseText(textFile, fontFaces["regular"])
	if err != nil {
		error.Error(err)
		// panic("panic! text file not properly loaded")
	}
	pages := makePages(cont.wrapped, LINES_PER_PAGE)
	// TODO: some sort of concurrency issue happening here? should I use locks...
	textLines := formatTextImages(pages[currentPage], fontFaces["regular"], MIN_X_TEXT)
	loadText := drawTextLines(textLines, fontFaces["regular"], &textBounds)
	env.Draw() <- loadText
	var p image.Point
	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				p = image.Pt(e.X, e.Y)

			case win.MoUp:
				p = image.Pt(e.X, e.Y)
				if p.In(textBounds) {
					textLines := formatTextImages(pages[currentPage], fontFaces["regular"], MIN_X_TEXT)
					env.Draw() <- drawTextLines(textLines, fontFaces["regular"], &textBounds)
					load := highlightLine(fontFaces["bold"], textLines, image.Pt(e.X, e.Y), words)
					env.Draw() <- load
				}
			case win.KbDown:
				if e.Key == win.KeyDown && currentPage < len(pages)-1 {
					currentPage += 1
					textLines := formatTextImages(pages[currentPage], fontFaces["regular"], MIN_X_TEXT)
					loadText := drawTextLines(textLines, fontFaces["regular"], &textBounds)
					env.Draw() <- loadText
				}
				if e.Key == win.KeyUp && currentPage > 0 {
					currentPage -= 1
					textLines := formatTextImages(pages[currentPage], fontFaces["regular"], MIN_X_TEXT)
					loadText := drawTextLines(textLines, fontFaces["regular"], &textBounds)
					env.Draw() <- loadText
				}
			case win.KbUp:
			}
		}
	}
}
