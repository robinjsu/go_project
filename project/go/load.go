package main

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	anchor = fixed.Point26_6{
		X: fixed.I(50),
		Y: fixed.I(50),
	}
	fontBounds = fixed.Rectangle26_6{}
	fontAdv    = fixed.I(0)
	bg         = image.Rect(0, 0, MAXWIDTH, MAXHEIGHT)
	fileString = ""
)

func setFont(fontFaces map[string]font.Face) {
	regFont := fontFaces["regular"]
	fontBounds, fontAdv = font.BoundString(regFont, "A")
	anchor = anchor.Add(fixed.P(0, fontBounds.Max.Y.Ceil()))
}

func makeSplash() func(draw.Image) image.Rectangle {
	drawSplash := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, bg, image.White, image.ZP, draw.Src)
		return bg
	}
	return drawSplash
}

func printLetter(r rune, face font.Face) draw.Image {
	fileString = fmt.Sprintf("%s%s", fileString, string(r))
	bgImage := image.NewRGBA(bg)
	text := &font.Drawer{
		Dst:  bgImage,
		Src:  image.Black,
		Face: face,
		Dot:  anchor,
	}
	text.DrawString(string(r))
	anchor = text.Dot
	return bgImage
}

func drawLetters(letterImg draw.Image) func(draw.Image) image.Rectangle {
	typing := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, letterImg.Bounds(), letterImg, image.ZP, draw.Over)
		return drw.Bounds()
	}
	return typing
}

func printKey(k win.Key) {
	fmt.Print(k)
}

func Load(env gui.Env, fontFaces map[string]font.Face, filepath chan<- string) {
	env.Draw() <- makeSplash()
	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.KbType:
				textImg := printLetter(e.Rune, fontFaces["regular"])
				env.Draw() <- drawLetters(textImg)
			case win.KbDown:
				if e.Key == win.KeyEnter {
					filepath <- fileString
				}
				// case win.MoDown:
				// case win.MoUp:
			}
		}
	}
}
