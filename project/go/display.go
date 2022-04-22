package main

import (
	"image"
	"image/draw"

	"github.com/faiface/gui"
)

func Display(env gui.Env) {

	loadPage := func(drw draw.Image) image.Rectangle {
		mainPage := image.Rect(0, 0, MAXWIDTH, HEIGHT)
		textBar := image.Rect(0, 0, .75*MAXWIDTH, HEIGHT)
		sideBar := image.Rect(0.75*MAXWIDTH, 0, MAXWIDTH, HEIGHT)
		draw.Draw(drw, mainPage, image.White, image.ZP, draw.Src)
		draw.Draw(drw, textBar, image.White, textBar.Min, draw.Src)
		draw.Draw(drw, sideBar, image.NewUniform(DEEP_BLUE), sideBar.Min, draw.Src)

		return mainPage
	}

	env.Draw() <- loadPage

	for {
		select {
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}
}
