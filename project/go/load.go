package main

import (
	"image"
	"image/draw"

	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	anchor = fixed.Point26_6{
		X: fixed.I(100),
		Y: fixed.I(100),
	}
	bg      = image.Rectangle{textBounds.Min, textBounds.Max.Add(image.Pt(0, 100))}
	bgImage = image.NewRGBA(bg)
)

func makeSplash() func(draw.Image) image.Rectangle {
	drawSplash := func(drw draw.Image) image.Rectangle {

		draw.Draw(drw, bg, image.NewUniform(SHADOW), image.ZP, draw.Src)
		return bg
	}
	return drawSplash
}

func printMsg(message string, face font.Face) draw.Image {
	text := &font.Drawer{
		Src:  image.Black,
		Face: face,
		Dot:  anchor,
	}

	bounds, _ := text.BoundString(message)
	letterBox := image.Rect(
		bounds.Min.X.Floor(),
		bounds.Min.Y.Floor(),
		bounds.Max.X.Ceil(),
		bounds.Max.Y.Ceil(),
	)
	text.Dst = image.NewRGBA(letterBox)
	text.DrawString(message)
	anchor = text.Dot
	return text.Dst
}

func drawLetters(letterImg draw.Image) func(draw.Image) image.Rectangle {
	typing := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, letterImg.Bounds(), letterImg, letterImg.Bounds().Min, draw.Over)
		return drw.Bounds()
	}
	return typing
}

func Load(env gui.Env, face font.Face, filepath chan<- string) {
	env.Draw() <- makeSplash()
	textImg := printMsg("DRAG AND DROP A .TXT FILE TO START...", face)
	env.Draw() <- drawLetters(textImg)
	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.PathDrop:
				filepath <- e.FilePath
			}
		}
	}
}
