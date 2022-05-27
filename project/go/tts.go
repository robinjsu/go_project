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

func getBtnIcons() []image.Image {
	imgFiles := []string{"images/play.png", "images/pause.png", "images/prev.png", "images/next.png"}
	var images []image.Image
	for _, img := range imgFiles {
		png, err := getPNG(img)
		if err != nil {
			fmt.Println(err)
		}
		images = append(images, png)
	}
	return images
}

func drawAudio(images []image.Image) (func(draw.Image) image.Rectangle, []image.Rectangle) {
	imgSz := images[0].Bounds().Size()
	footerSz := footer.Bounds().Size()
	pos := footer.Min.Add(image.Pt(10, (footerSz.Y-imgSz.Y)/2))
	pngR := image.Rectangle{pos, pos.Add(imgSz)}
	iconArea := image.Rectangle{pos, pos.Add(image.Pt(imgSz.X*4, imgSz.Y))}
	icons := image.NewRGBA(iconArea)
	iconsR := []image.Rectangle{}

	for range images {
		iconsR = append(iconsR, pngR)
		pngR = pngR.Add(image.Pt(imgSz.X, 0))
	}

	drawPNG := func(drw draw.Image) image.Rectangle {
		for i, img := range iconsR {
			draw.Draw(icons, img, images[i], image.ZP, draw.Over)
		}
		draw.Draw(drw, icons.Bounds(), icons, icons.Bounds().Min, draw.Over)
		return icons.Bounds()
	}
	return drawPNG, iconsR
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

// TODO: streamer still nto working!
func playBack(audio *os.File, start <-chan bool, pause <-chan bool, restart <-chan bool, done <-chan bool, newStream <-chan beep.StreamSeekCloser) {
	var controller *beep.Ctrl
	var streamer beep.StreamSeekCloser
	_, format, err := mp3.Decode(audio)
	if err != nil {
		log.Fatal(err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	controller = &beep.Ctrl{Streamer: nil, Paused: false}
	for {
		select {
		case streamer = <-newStream:
			controller.Streamer = streamer
		case <-start:
			speaker.Play(controller)
		case paused := <-pause:
			speaker.Lock()
			controller.Paused = paused
			speaker.Unlock()
		case <-restart:
			err := streamer.Seek(0)
			if err != nil {
				log.Fatal(ReadError{err})
			}
		case <-done:
			controller.Streamer = nil
		}
	}
}

func TextToSpeech(env gui.Env, load chan bool, content <-chan [][]string) {
	var pages [][]string

	audioStarted := false
	paused := false
	audioBtns := getBtnIcons()
	audioBtn, iconsR := drawAudio(audioBtns)
	newStreamer := make(chan beep.StreamSeekCloser)
	start := make(chan bool)
	restart := make(chan bool)
	pause := make(chan bool)
	done := make(chan bool)

	pg := 0

	for {
		select {
		case ready := <-load:
			if ready == true {
				env.Draw() <- audioBtn
			}
		case pages = <-content:
			fmt.Println(pages[0])
			// getSpeech(string(strings.Join(pages[0], " ")), fmt.Sprintf("%s/pg-%v.mp3", audioDir, i))
		case e, ok := <-env.Events():
			if !ok {
				// audioFiles, err := os.ReadDir(audioDir)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// for _, audio := range audioFiles {
				// 	os.Remove(audio.Name())
				// }
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				switch {
				case e.Point.In(iconsR[0]):
					fmt.Println("play")
					if audioStarted && paused == false {
						restart <- true
					} else if audioStarted && paused == true {
						paused = false
						pause <- false
					} else {
						audio, err := os.Open(fmt.Sprintf("%s/pg-%v.mp3", audioDir, pg))
						if err != nil {
							log.Fatal(err)
						}
						audioStream, _, _ := mp3.Decode(audio)
						go playBack(audio, start, pause, restart, done, newStreamer)
						newStreamer <- audioStream
						start <- true
						audioStarted = true
					}
				case e.Point.In(iconsR[1]):
					fmt.Println("pause")
					paused = true
					pause <- paused
				case e.Point.In(iconsR[2]):
					fmt.Println("prev")
					if pg > 0 {
						pg--
						done <- true
						audio, err := os.Open(fmt.Sprintf("%s/pg-%v.mp3", audioDir, pg))
						if err != nil {
							log.Fatal(err)
						}
						streamer, _, err := mp3.Decode(audio)
						if err != nil {
							log.Fatal(err)
						}
						newStreamer <- streamer
					}
				case e.Point.In(iconsR[3]):
					fmt.Println("next")
					if pg < len(pages)-1 {
						pg++
						done <- true
						audio, err := os.Open(fmt.Sprintf("%s/pg-%v.mp3", audioDir, pg))
						if err != nil {
							log.Fatal(err)
						}
						streamer, _, _ := mp3.Decode(audio)
						newStreamer <- streamer
					}
				}
			case win.KbDown:
				if e.Key == win.KeySpace {
					paused = !paused
					pause <- paused
				}
			}
		}

	}
}
