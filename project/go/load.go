package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	anchor = fixed.Point26_6{
		X: fixed.I(50),
		Y: fixed.I(50),
	}
	fontBounds = fixed.Rectangle26_6{}
	fontAdv    = fixed.I(0)
	bg         = image.Rect(0, 0, MAXWIDTH, MAXHEIGHT)
	bgImage    = image.NewRGBA(bg)
	fileString = ""
	letters    = []rune{}
)

func setFont(face font.Face) {
	// regFont := fontFaces["regular"]
	fontBounds, fontAdv = font.BoundString(face, "A")
	fmt.Println(fontAdv)
	anchor = anchor.Add(fixed.P(0, fontBounds.Max.Y.Ceil()))
}

func makeSplash() func(draw.Image) image.Rectangle {
	drawSplash := func(drw draw.Image) image.Rectangle {

		draw.Draw(drw, bg, image.NewUniform(TEAL), image.ZP, draw.Src)
		return bg
	}
	return drawSplash
}

// func printLetter(r rune, face font.Face) draw.Image {
// 	fileString = fmt.Sprintf("%s%s", fileString, string(r))
// 	letters = append(letters, r)
// 	text := &font.Drawer{
// 		Src:  image.Black,
// 		Face: face,
// 		Dot:  anchor,
// 	}

// 	bounds, _ := text.BoundString(fileString)
// 	letterBox := image.Rect(
// 		bounds.Min.X.Floor(),
// 		bounds.Min.Y.Floor(),
// 		bounds.Max.X.Ceil(),
// 		bounds.Max.Y.Ceil(),
// 	)
// 	text.Dst = image.NewRGBA(letterBox)
// 	text.DrawString(string(r))
// 	anchor = text.Dot
// 	return text.Dst
// }
func printLetter(message string, face font.Face) draw.Image {
	// fileString = fmt.Sprintf("%s%s", fileString, string(r))
	// letters = append(letters, r)
	text := &font.Drawer{
		Src:  image.Black,
		Face: face,
		Dot:  anchor,
	}

	bounds, _ := text.BoundString(message)
	letterBox := image.Rect(
		bounds.Min.X.Floor(),
		bounds.Min.Y.Floor(),
		bounds.Max.X.Ceil(),
		bounds.Max.Y.Ceil(),
	)
	text.Dst = image.NewRGBA(letterBox)
	text.DrawString(message)
	anchor = text.Dot
	return text.Dst
}

func deleteLetter(face font.Face) func(draw.Image) image.Rectangle {
	anchor = anchor.Sub(fixed.Point26_6{X: fontAdv, Y: 0})

	fmt.Println(anchor)
	lastRune := letters[len(letters)-1]
	letters = letters[:len(letters)-1]
	fileString = strings.TrimSuffix(fileString, string(lastRune))
	del := &font.Drawer{
		Src:  image.NewUniform(color.White),
		Face: face,
		Dot:  anchor,
	}

	bounds, _ := del.BoundString(string(lastRune))
	letterBox := image.Rect(
		bounds.Min.X.Floor(),
		bounds.Min.Y.Floor(),
		bounds.Max.X.Ceil(),
		bounds.Max.Y.Ceil(),
	)
	del.Dst = image.NewRGBA(letterBox)
	del.DrawString(string(lastRune))

	delete := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, letterBox, image.White, image.ZP, draw.Over)
		return drw.Bounds()
	}
	return delete
}

func drawLetters(letterImg draw.Image) func(draw.Image) image.Rectangle {
	typing := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, letterImg.Bounds(), letterImg, letterImg.Bounds().Min, draw.Over)
		return drw.Bounds()
	}
	return typing
}

func printKey(k win.Key) {
	fmt.Print(k)
}

func Load(env gui.Env, face font.Face, filepath chan<- string) {
	setFont(face)
	env.Draw() <- makeSplash()
	textImg := printLetter("DRAG A .TXT FILE OVER THIS WINDOW TO START...", face)
	env.Draw() <- drawLetters(textImg)
	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.PathDrop:
				filepath <- e.FilePath
			case win.KbType:
				// textImg := printLetter(e.Rune, face)
				// env.Draw() <- drawLetters(textImg)
			case win.KbDown:
				switch e.Key {
				// case win.KeyEnter:
				// 	filepath <- fileString
				case win.KeyBackspace:
					if fileString != "" {
						// img := deleteLetter(fontFaces["regular"])
						env.Draw() <- deleteLetter(face)
					}

				}
				// case win.MoDown:
				// case win.MoUp:
			}
		}
	}
}
