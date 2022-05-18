package main

import (
	"image"
	"image/draw"

	gui "github.com/faiface/gui"
)

func Display(env gui.Env, load <-chan bool) {

	loadPage := func(drw draw.Image) image.Rectangle {
		sideBar := image.Rect(MIN_X_SEARCH, 0, MAXWIDTH, MAXHEIGHT)
		draw.Draw(drw, sideBar, image.NewUniform(DEEP_BLUE), sideBar.Min, draw.Over)

		return sideBar
	}

	for {
		select {
		case ready := <-load:
			if ready {
				env.Draw() <- loadPage
			}
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
		}
	}
}
