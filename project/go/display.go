package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/faiface/gui"
)

func Display(env gui.Env) {

	loadPage := func(drw draw.Image) image.Rectangle {
		mainPage := image.Rect(0, 0, WIDTH, HEIGHT)
		textBar := image.Rect(0, 0, .75*WIDTH, HEIGHT)
		sideBar := image.Rect(0.75*WIDTH, 0, WIDTH, HEIGHT)
		draw.Draw(drw, mainPage, image.White, image.ZP, draw.Src)
		draw.Draw(drw, textBar, image.White, textBar.Min, draw.Src)
		draw.Draw(drw, sideBar, image.NewUniform(color.RGBA{0, 0, 255, 100}), sideBar.Min, draw.Src)

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
