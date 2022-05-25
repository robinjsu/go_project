package main

import (
	"fmt"
	"image"
	"image/draw"

	gui "github.com/faiface/gui"
)

func drawAudio(img image.Image) func(draw.Image) image.Rectangle {
	imgSz := img.Bounds().Size()
	footerSz := footer.Bounds().Size()
	padH := (footerSz.Y - imgSz.Y) / 2
	pos := footer.Min.Add(image.Pt(10, padH))
	drawPNG := func(drw draw.Image) image.Rectangle {
		pngR := image.Rectangle{pos, pos.Add(imgSz)}
		draw.Draw(drw, pngR, img, image.ZP, draw.Over)

		return pngR.Bounds()
	}
	return drawPNG
}

func TextToSpeech(env gui.Env, load chan bool) {
	audioPng, err := getPNG("images/audio_icon.png")
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case ready := <-load:
			if ready == true {
				fmt.Println("tts received msg")
				env.Draw() <- drawAudio(audioPng)
				load <- true
			}
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			// switch e := e.(type) {

			// }
		}
	}
}
