package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"time"

	tts "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	gui "github.com/faiface/gui"
	"github.com/faiface/gui/win"
)

func drawAudio(img image.Image) (func(draw.Image) image.Rectangle, image.Rectangle) {
	imgSz := img.Bounds().Size()
	footerSz := footer.Bounds().Size()
	pos := footer.Min.Add(image.Pt(10, (footerSz.Y-imgSz.Y)/2))
	pngR := image.Rectangle{pos, pos.Add(img.Bounds().Size())}

	drawPNG := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, pngR, img, image.ZP, draw.Over)

		return pngR.Bounds()
	}
	return drawPNG, pngR
}

func getSpeech(content string, outFile string) string {
	ctx := context.Background()

	ttsClient, err := tts.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer ttsClient.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: content},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := ttsClient.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(outFile, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return outFile
}

func playBack(audio *os.File, start <-chan bool, pause <-chan bool, restart <-chan bool, done <-chan bool) {
	var controller *beep.Ctrl

	stream, format, err := mp3.Decode(audio)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	controller = &beep.Ctrl{Streamer: stream, Paused: false}
	for {
		select {
		case <-start:
			speaker.Play(controller)
		case <-pause:
			speaker.Lock()
			controller.Paused = !controller.Paused
			speaker.Unlock()
		case <-restart:
			err = stream.Seek(0)
			if err != nil {
				log.Fatal(ReadError{err})
			}
		case <-done:
			speaker.Lock()
			speaker.Clear()
			speaker.Unlock()
			break
		}
	}
}

func TextToSpeech(env gui.Env, load chan bool, content <-chan [][]string) {
	var audioStarted = false
	// var filename = ""
	audioPng, err := getPNG("images/audio_icon.png")
	if err != nil {
		fmt.Println(err)
	}
	audioBtn, pngR := drawAudio(audioPng)
	start := make(chan bool)
	restart := make(chan bool)
	pause := make(chan bool)
	done := make(chan bool)
	// var pages [][]string
	pg := 0

	for {
		select {
		case ready := <-load:
			if ready == true {
				env.Draw() <- audioBtn
			}
		// case pages = <-content:
		// filename = getSpeech(string(strings.Join(pages[pg], " ")), fmt.Sprintf("pg-%v.mp3", pg))
		// audio, err := os.Open(filename)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// go playBack(audio, start, pause, restart, done)
		// case p := <-page:
		// 	switch p {
		// 	case "prev":
		// 		if pg > 0 {
		// 			pg--
		// 		}
		// 	case "next":
		// 		if pg < len(pages)-1 {
		// 			pg++
		// 		}
		// 	}
		// done <- true
		// filename = getSpeech(string(strings.Join(pages[pg], " ")), fmt.Sprintf("pg-%v.mp3", pg))
		// audio, err := os.Open(filename)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// go playBack(audio, start, pause, restart, done)
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				if e.Point.In(pngR) {
					if audioStarted {
						restart <- true
					} else {
						audio, err := os.Open(fmt.Sprintf("pg-%v.mp3", pg))
						if err != nil {
							log.Fatal(err)
						}
						go playBack(audio, start, pause, restart, done)
						start <- true
						audioStarted = true
					}
				}
			case win.KbDown:
				if e.Key == win.KeySpace {
					pause <- true
				}
			}
		}
	}
}
