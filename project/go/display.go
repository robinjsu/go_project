package main

import (
	"image"
	"image/draw"

	gui "github.com/faiface/gui"
)

func Display(env gui.Env) {
	loadPage := func(drw draw.Image) image.Rectangle {
		sideBar := image.Rect(MIN_X_SEARCH, 0, MAXWIDTH, MAXHEIGHT)
		draw.Draw(drw, sideBar, image.NewUniform(DEEP_BLUE), sideBar.Min, draw.Src)

		return sideBar
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
