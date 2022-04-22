package main

import (
	"image"
	"image/draw"

	"golang.org/x/image/math/fixed"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
)

// Content contains a buffer with text content; each line of text corresponds to text up until the next newline character is found
// TODO: changing approach to parsing - first take in each paragraph, then have separate function to split each p into smaller parts (maybe even as words) in order to
// figure out a way to track location of each word and how it should be tracked so that each word can be clicked on to search
// TODO: take paragraph breaks into account, and split at newline if exists within the lineMaxW

func loadTxt(face font.Face, cont *Content) func(drw draw.Image) image.Rectangle {
	load := func(drw draw.Image) image.Rectangle {
		// coordinates refer to the destination image's coordinate space
		page := image.Rect(0, 0, 900, 900)
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for i, lns := range cont.format {
			text := &font.Drawer{
				Dst:  drw,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(FONTSZ, (i+1)*face.Metrics().Height.Ceil()),
			}
			cont.format[i].bounds, cont.format[i].span = text.BoundString(lns.txt)
			text.DrawString(lns.txt)
		}
		return page
	}
	return load
}

func highlightLine(face font.Face, cont *Content, p image.Point, words chan<- string) func(draw.Image) image.Rectangle {
	var line image.Rectangle
	load := func(drw draw.Image) image.Rectangle {
		var txt string
		for _, ln := range cont.format {
			rct := ln.bounds
			if p.Y >= (rct.Min.Y).Floor() && p.Y <= (rct.Max.Y).Ceil() && p.X <= TEXTWIDTH {
				txt = ln.txt
				line = image.Rect((rct.Min.X).Floor(), (rct.Min.Y).Floor(), (rct.Max.X).Floor(), (rct.Max.Y).Floor())
				draw.Draw(drw, line, &image.Uniform{HIGHLIGHT_GRAY}, image.ZP, draw.Over)
			}
		}
		// send words to Search component
		words <- txt
		return line
	}
	return load
}

func Text(env gui.Env, textFile string, fontFaces map[string]font.Face, words chan<- string) {
	textBounds := image.Rect(0, 0, 900, 900)
	// fontFaces := loadFonts(FONT_REG, FONT_BOLD)
	cont := NewContent()
	_, err := cont.parseText(textFile, fontFaces["regular"])
	if err != nil {
		error.Error(err)
		// panic("panic! text file not properly loaded")
	}
	loadText := loadTxt(fontFaces["regular"], cont)
	env.Draw() <- loadText

	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
			case win.MoUp:
				p := image.Pt(e.X, e.Y)
				if p.In(textBounds) {
					env.Draw() <- loadTxt(fontFaces["regular"], cont)
					load := highlightLine(fontFaces["bold"], cont, image.Pt(e.X, e.Y), words)
					env.Draw() <- load
				}
			}
		}
	}
}
