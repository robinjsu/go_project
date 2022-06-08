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

// func getGoogleClient() (*tts.Client, context.Context) {
// 	ctx := context.Background()
// 	jsonCreds, err := os.ReadFile("./tts_client_secret.json")
// 	if err != nil {
// 		log.Fatal(FileError{"tts_client_secret.json", err})
// 	}
// 	ttsClient, err := tts.NewClient(ctx, option.WithCredentialsJSON(jsonCreds))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return ttsClient, ctx
// }

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

func getSpeech(content string, outDir string, outFile string, client *tts.Client, ctx context.Context) string {
	// ctx := context.Background()

	// ttsClient, err := tts.NewClient(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer ttsClient.Close()

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

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	err = ioutil.WriteFile(outFile, resp.AudioContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return outFile
}

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

func TextToSpeech(env gui.Env, load chan bool, content <-chan [][]string, client *tts.Client, ctx context.Context) {
	var pages [][]string
	var streamers []beep.StreamSeekCloser
	var title string
	var authorized bool
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

	// creds, _ := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	// if creds == "" {
	// 	audioDir = "example_audio"
	// 	authorized = false
	// } else {
	// 	authorized = true
	// }
	// client, ctx := getGoogleClient()
	if client != nil {
		authorized = true
	} else {
		audioDir = "example_audio"
	}
	defer client.Close()

	for {
		select {
		case ready := <-load:
			if ready == true {
				env.Draw() <- audioBtn
			}
		case pages = <-content:
			for i, page := range pages {
				if page != nil {
					audioFile := fmt.Sprintf("%s/%s/pg-%v.mp3", audioDir, title, i)
					if authorized {
						getSpeech(string(strings.Join(page, " ")), title, audioFile, client, ctx)
					}
					audio, err := os.Open(audioFile)
					if err != nil {
						log.Fatal(FileError{audioFile, err})
					}
					stream, _, err := mp3.Decode(audio)
					if err != nil {
						log.Fatal(FileError{audioFile, err})
					}
					streamers = append(streamers, stream)
				}
			}
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
					if audioStarted && paused == false {
						restart <- true
					} else if audioStarted && paused == true {
						paused = false
						pause <- false
					} else {
						audio, err := os.Open(fmt.Sprintf("%s/%s/pg-%v.mp3", audioDir, title, pg))
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
					paused = true
					pause <- paused
				case e.Point.In(iconsR[2]):
					if pg > 0 {
						pg--
						done <- true
						streamers[pg].Seek(0)
						newStreamer <- streamers[pg]
					}
				case e.Point.In(iconsR[3]):
					if pg < len(pages)-1 {
						pg++
						done <- true
						streamers[pg].Seek(0)
						newStreamer <- streamers[pg]
					}
				}
			case win.KbDown:
				if e.Key == win.KeySpace {
					paused = !paused
					pause <- paused
				}
			case win.PathDrop:
				path := e.FilePath
				dirs := strings.Split(path, "/")
				title = strings.TrimSuffix(dirs[len(dirs)-1], ".txt")
				if _, err := os.Stat(fmt.Sprintf("%s", audioDir)); err != nil {
					err = os.Mkdir(fmt.Sprintf("%s", audioDir), 0777)
					err = os.Mkdir(fmt.Sprintf("%s/%s", audioDir, title), 0777)
				} else if _, err = os.Stat(fmt.Sprintf("%s/%s", audioDir, title)); err != nil {
					err = os.Mkdir(fmt.Sprintf("%s/%s", audioDir, title), 0777)
				}

			}
		}

	}
}
