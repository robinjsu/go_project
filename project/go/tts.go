package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func getSpeech(content string) string {
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

	filename := "audioOutput.mp3"
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func playBack(audio *os.File, start chan bool, playback chan bool) {
	stream, format, err := mp3.Decode(audio)
	if err != nil {
		log.Fatal(err)
	}
	stream.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	controller := &beep.Ctrl{Streamer: stream, Paused: false}
	for {
		select {
		case <-start:
			go speaker.Play(controller)
		case <-playback:
			fmt.Println("received signal")
			speaker.Lock()
			controller.Paused = true
			speaker.Unlock()
		}
	}
}

func TextToSpeech(env gui.Env, load chan bool, content <-chan *Content) {
	audioPng, err := getPNG("images/audio_icon.png")
	if err != nil {
		fmt.Println(err)
	}
	audioBtn, pngR := drawAudio(audioPng)
	// start := make(chan bool)
	// pause := make(chan bool)

	for {
		select {
		case ready := <-load:
			if ready == true {
				env.Draw() <- audioBtn
				load <- true
			}
		case text := <-content:
			pages := makePages(text.wrapped, LINES_PER_PAGE)
			filename := getSpeech(string(strings.Join(pages[0], " ")))
			_, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}

			// go playBack(audio, start, pause)
		// start <- true

		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				fmt.Println(e)
				if e.Point.In(pngR) {
					fmt.Println("found audio btn")
					// speaker.Play(playback)
				}
			case win.KbDown:
				if e.Key == win.KeySpace {
					fmt.Println(e)
					// pause <- true
				}
			}
		}
	}
}
