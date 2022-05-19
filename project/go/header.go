package main

import (
	"image"
	"image/draw"
	"strings"

	gui "github.com/faiface/gui"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func setHeader(word string, bounds *image.Rectangle, lineHeight int, face font.Face) func(draw.Image) image.Rectangle {

	trimmed := strings.Trim(word, " .,:;?!'\"“”()[]")
	wordHt := face.Metrics().Height.Floor()
	wordWd := font.MeasureString(face, word)
	dotX := bounds.Min.X + ((bounds.Dx() - wordWd.Ceil()) / 2)
	dotY := bounds.Min.Y + ((bounds.Dy() - wordHt) / 2) + wordHt

	wrd := &font.Drawer{
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(dotX, dotY),
	}
	wrdBnds, _ := wrd.BoundString(trimmed)

	wrd.Dst = image.NewRGBA(image.Rect(
		wrdBnds.Min.X.Floor(),
		wrdBnds.Min.Y.Floor(),
		wrdBnds.Max.X.Ceil(),
		wrdBnds.Max.Y.Ceil(),
	))
	header := func(drw draw.Image) image.Rectangle {
		newR := makeInsetR(*bounds, MARGIN/2)
		draw.Draw(drw, newR, &image.Uniform{TEAL}, image.Pt(0, 0), draw.Over)
		wrd.DrawString(word)
		draw.Draw(drw, wrd.Dst.Bounds(), wrd.Dst, wrd.Dst.Bounds().Min, draw.Over)
		return newR
	}
	return header

}

func Header(env gui.Env, fontFaces map[string]font.Face, words <-chan string, define chan<- string) {
	for {
		select {
		case lookup := <-words:
			lineHeight := fontFaces["regular"].Metrics().Height.Floor()
			env.Draw() <- setHeader(lookup, &wordCorner, lineHeight, fontFaces["bold"])
			define <- lookup
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}
}
