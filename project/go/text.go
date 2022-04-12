package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"

	"golang.org/x/image/math/fixed"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	FONTSZ     = 14
	FONTFAMILY = "../../fonts/Karma/Karma-Regular.ttf"
	FONT_H     = 20
	NEWLINE    = byte('\n')
)

// func check_err(err error) error {
// 	if err != nil && err != io.EOF {
// 		return err
// 	}
// }

type Content struct {
	lines []string
	num   int
}

func NewContent() *Content {
	c := Content{num: 0}
	return &c
}

func (cont *Content) append(line string) {
	fmt.Printf("%s | ", string(line))
	cont.lines = append(cont.lines, line)
	cont.num++
}

func parseText(cont *Content, filename string, bounds int) (int, error) {
	// textFile, err := os.Open(filename)
	textBytes, err := os.ReadFile(filename)
	if err != nil {
		return -1, err
	}

	// reader := bufio.NewReader(strings.NewReader(string(textBytes)))
	// buffer := make([]byte, bounds)
	buffer := bytes.NewBuffer(textBytes)

	for err != io.EOF {
		// nBytes, err := reader.Peek(bounds)
		nBytes := buffer.Next(bounds)
		// if err != nil {
		// 	return -1, err
		// }
		idx := bytes.IndexByte(nBytes, '\n')
		if idx != -1 {
			// buffer, err = reader.ReadBytes('\n')
			line, err := buffer.ReadBytes('\n')
			if err != nil {
				return -1, err
			}
			cont.append(string(line))
		} else {
			pbytes := make([]byte, bounds)
			// _, err = reader.Read(buffer)
			_, err = buffer.Read(pbytes)
			if err != nil {
				return -1, err
			}
			cont.append(string(pbytes))
		}
	}

	return cont.num, nil
}

func parseFont(file string) (*tt.Font, error) {
	ttfFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ttf, err := tt.Parse(ttfFile)
	if err != nil {
		return nil, err
	}

	return ttf, nil
}

func Text(env gui.Env, textFile string) {
	// parse bytes and return a pointer to a Font type object
	fontstyle, err := parseFont(FONTFAMILY)
	if err != nil {
		error.Error(err)
		panic("panic! TTF file not properly loaded")
	}
	// create face, which provides the `glyph mask images``
	face := tt.NewFace(fontstyle, &tt.Options{
		// options... here just font size (0 is 12-point default)
		Size: 0,
	})

	cont := NewContent()
	_, err = parseText(cont, textFile, 100)
	if err != nil {
		error.Error(err)
		// panic("panic! text file not properly loaded")
	}

	loadText := func(drw draw.Image) image.Rectangle {
		page := image.Rect(0, 0, 450, 600)
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for i := 0; i < 10; i++ {
			// line := image.Rect(0, FONT_H*(i+1), 900, (FONT_H*(i+1))+FONTSZ)
			text := &font.Drawer{
				Dst:  drw,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(FONTSZ, (FONT_H*(i+1))+FONTSZ),
			}
			// draw.Draw(drw, line, image.White, line.Min, draw.Src)
			text.DrawString(cont.lines[i])
		}

		return page
	}

	env.Draw() <- loadText

	for {
		select {
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				fmt.Println(e.String())
			}
		}
	}
}
