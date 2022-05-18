package main

import (
	"image"
	"image/draw"

	"github.com/faiface/gui"
)

func Display(env gui.Env, load <-chan string) {

	loadPage := func(drw draw.Image) image.Rectangle {
		sideBar := image.Rect(MIN_X_SEARCH, 0, MAXWIDTH, MAXHEIGHT)
		draw.Draw(drw, sideBar, image.NewUniform(DEEP_BLUE), sideBar.Min, draw.Over)

		return sideBar
	}

	for {
		select {
		case f := <-load:
			if f == "ok" {
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
