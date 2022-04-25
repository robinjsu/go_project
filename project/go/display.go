package main

import (
	"image"
	"image/draw"

	"github.com/faiface/gui"
)

func Display(env gui.Env) {

	loadPage := func(drw draw.Image) image.Rectangle {
		// TODO: standardize points
		// TODO: refactor drawing funcs!
		mainPage := image.Rect(0, 0, MAXWIDTH, MAXHEIGHT)
		textBar := image.Rect(0, 0, MAX_X_TEXT, MAXHEIGHT)
		sideBar := image.Rect(750, 0, MAXWIDTH, MAXHEIGHT)
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
