package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	gui "github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func setBtnArea(button image.Rectangle, r image.Rectangle) []image.Rectangle {
	newX := ((r.Dx() - (button.Dx() * 2)) / 2) + r.Min.X
	newY := ((r.Dy() - button.Dy()) / 2) + r.Min.Y
	prevBtn := image.Rect(newX, newY, newX+button.Dx(), newY+button.Dy())
	nextBtn := image.Rect(newX+button.Dx(), newY, newX+(button.Dx()*2), newY+button.Dy())

	return []image.Rectangle{prevBtn, nextBtn}
}

func setBtnText(btn image.Rectangle, text string, face font.Face) draw.Image {
	textBounds, _ := font.BoundString(face, text)
	textR := image.Rect(
		textBounds.Min.X.Floor(),
		textBounds.Min.Y.Floor(),
		textBounds.Max.X.Ceil(),
		textBounds.Max.Y.Ceil(),
	)

	xPos := ((btn.Dx() - textR.Dx()) / 2) + btn.Min.X
	yPos := ((btn.Dy() - textR.Dy()) / 2) + btn.Min.Y

	inset := image.Rect(xPos, yPos, xPos+textR.Dx(), yPos+textR.Dy())
	insetImg := image.NewRGBA(inset.Bounds())

	btnTxt := &font.Drawer{
		Dst:  insetImg,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(xPos, yPos+textR.Dy()),
	}

	btnTxt.DrawString(text)

	return insetImg

}

func drawBtns(env gui.Env, prevBtn image.Rectangle, nextBtn image.Rectangle, r image.Rectangle, face font.Face, done chan bool) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	counter := 0.0

	prevClr := PRUSSIAN_BLUE
	nextClr := DARK_BLUE
	gradient := func(color *color.RGBA, tick time.Time) color.RGBA {
		counter += 1
		sineFunc := (math.Sin(counter/500) + 1) * (255 / 2)
		color.B = uint8(sineFunc)
		return *color
	}

	prev := setBtnText(prevBtn, "PREV", face)
	next := setBtnText(nextBtn, "NEXT", face)

	for {
		select {
		case tick := <-ticker.C:
			env.Draw() <- func(drw draw.Image) image.Rectangle {
				draw.Draw(drw, prevBtn, &image.Uniform{gradient(&prevClr, tick)}, prevBtn.Min, draw.Src)
				draw.Draw(drw, prev.Bounds(), prev, prev.Bounds().Min, draw.Over)
				draw.Draw(drw, nextBtn, &image.Uniform{gradient(&nextClr, tick)}, nextBtn.Min, draw.Src)
				draw.Draw(drw, next.Bounds(), next, next.Bounds().Min, draw.Over)
				return r
			}
		case <-done:
			break
		}
	}
}

func PagingBtns(env gui.Env, prev chan<- bool, next chan<- bool, faces map[string]font.Face, ready chan bool) {
	// loaded := false
	done := make(chan bool)
	btns := setBtnArea(buttonR, footer)
	for {
		select {
		case <-ready:
			go drawBtns(env, btns[0], btns[1], footer, faces["bold"], done)
		case e, ok := <-env.Events():
			if !ok {
				done <- true
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				if e.Point.In(btns[0]) {
					prev <- true
				} else if e.Point.In(btns[1]) {
					next <- true
				}
			}
		}
	}
}
