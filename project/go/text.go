package main

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"strings"

	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"

	"golang.org/x/image/font"
)

var (
	currentPage = 0
	fileLoaded  = false
)

func drawTextLines(images []imageObj, face font.Face, bounds *image.Rectangle) func(drw draw.Image) image.Rectangle {
	load := func(drw draw.Image) image.Rectangle {
		// coordinates refer to the destination image's coordinate space
		page := *bounds
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for _, obj := range images {
			draw.Draw(drw, obj.placement, obj.img, image.Pt(0, 0), draw.Src)
		}
		return page
	}
	return load
}

func findWord(face font.Face, line imageObj, p image.Point) (image.Rectangle, string) {
	anchor := image.Pt(line.placement.Min.X+line.text.bounds.Min.X.Floor(), line.placement.Min.Y+line.text.bounds.Min.Y.Floor())
	words := strings.Split(line.text.txt, " ")
	for _, w := range words {
		_, adv := font.BoundString(face, w)
		wordR := image.Rect(anchor.X, anchor.Y, anchor.X+adv.Floor(), anchor.Y+line.text.bounds.Max.Y.Ceil())
		if p.In(wordR) {
			// fmt.Printf("found word: %v", wordR)
			return wordR, w
		}
		ln := fmt.Sprintf("%s ", w)
		anchor = anchor.Add(image.Pt(font.MeasureString(face, ln).Round(), 0))
	}
	return image.Rect(0, 0, 0, 0), ""
}

func highlightWord(face font.Face, images []imageObj, p image.Point, words chan<- string) func(draw.Image) image.Rectangle {
	var line image.Rectangle
	load := func(drw draw.Image) image.Rectangle {
		// var txt string
		var wrd string
		var wordBounds image.Rectangle
		for _, ln := range images {
			rct := ln.placement.Bounds()
			if p.In(rct) {
				wordBounds, wrd = findWord(face, ln, p)
				draw.Draw(drw, wordBounds, &image.Uniform{LIGHT_GRAY}, image.ZP, draw.Over)
			}
		}
		// send words to Search component
		words <- wrd
		return line
	}
	return load
}

func Text(env gui.Env, fontFaces map[string]font.Face, words chan<- string, filepath chan string, load chan bool, page chan string, content chan<- [][]string) {
	var cont *Content = NewContent()
	var textLines []imageObj
	var pages [][]string
	var p image.Point
	var loadText func(draw.Image) image.Rectangle
	var lineHeight int

	for {
		select {
		case file := <-filepath:
			_, err := cont.parseText(file, fontFaces["regular"], &textBounds)
			if err != nil {
				fmt.Println(errors.Unwrap(err))
			}
			// content <- cont
			pages = makePages(cont.wrapped, LINES_PER_PAGE)
			content <- pages
			lineHeight = fontFaces["regular"].Metrics().Height.Ceil() * 2
			textLines = formatTextImages(pages[currentPage], 0, MIN_X_TEXT, MARGIN, lineHeight, fontFaces["regular"])
			loadText = drawTextLines(textLines, fontFaces["regular"], &textBounds)
			env.Draw() <- loadText
			fileLoaded = true
			load <- fileLoaded
		case pg := <-page:
			if fileLoaded {
				switch pg {
				case "prev":
					if currentPage > 0 {
						currentPage -= 1
					}
					textLines = formatTextImages(pages[currentPage], 0, MIN_X_TEXT, MARGIN, lineHeight, fontFaces["regular"])
					loadText = drawTextLines(textLines, fontFaces["regular"], &textBounds)
					env.Draw() <- loadText
				case "next":
					if currentPage < len(pages)-1 {
						currentPage += 1
					}
					textLines = formatTextImages(pages[currentPage], 0, MIN_X_TEXT, MARGIN, lineHeight, fontFaces["regular"])
					loadText = drawTextLines(textLines, fontFaces["regular"], &textBounds)
					env.Draw() <- loadText
				}
			}
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			if fileLoaded {
				switch e := e.(type) {
				case win.MoDown:
				case win.MoUp:
					p = image.Pt(e.X, e.Y)
					if p.In(textBounds) {
						env.Draw() <- drawTextLines(textLines, fontFaces["regular"], &textBounds)
						loadText = highlightWord(fontFaces["bold"], textLines, p, words)
						env.Draw() <- loadText
					}
				case win.KbDown:
					if e.Key == win.KeyUp && currentPage > 0 {
						currentPage -= 1
					} else if e.Key == win.KeyDown && currentPage < len(pages)-1 {
						currentPage += 1
					}
					textLines = formatTextImages(pages[currentPage], 0, MIN_X_TEXT, MARGIN, lineHeight, fontFaces["regular"])
					loadText = drawTextLines(textLines, fontFaces["regular"], &textBounds)
					env.Draw() <- loadText
				}
			}
		}
	}
}
