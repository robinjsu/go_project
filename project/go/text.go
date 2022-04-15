package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"strings"

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

// Content contains a buffer with text content; each line of text corresponds to text up until the next newline character is found
// TODO: changing approach to parsing - first take in each paragraph, then have separate function to split each p into smaller parts (maybe even as words) in order to
// figure out a way to track location of each word and how it should be tracked so that each word can be clicked on to search

type Sentence struct {
	line      string
	numBytes  int
	numPixels int
}
type Pgraph struct {
	// single paragraph, split into slice of strings split at each space
	lines []Sentence
	// number of bytes in the entire string
	num int
}
type Content struct {
	pgraph []Pgraph
	num    int
}

func NewContent() *Content {
	c := Content{}
	return &c
}

// func (cont *Content) append(line string) {
// 	// fmt.Printf("%s | ", string(line))
// 	cont.lines = append(cont.lines, line)
// 	cont.num++
// }

func (c *Content) Store(pg []byte) error {
	newPg := Pgraph{}
	line := strings.TrimSuffix(string(pg), "\n")
	sentences := strings.SplitAfter(line, ". ")
	for _, lin := range sentences {
		sentence := Sentence{line: lin, numBytes: len(lin), numPixels: 0}
		newPg.lines = append(newPg.lines, sentence)
	}
	newPg.num = len(sentences)
	c.pgraph = append(c.pgraph, newPg)
	c.num++
	return nil
}

func parseText(cont *Content, filename string, bounds int) (int, error) {
	// c := NewContent()
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file! %v\n", err)
		return -1, err
	}

	buffer := bytes.NewBuffer(content)
	p, err := buffer.ReadBytes('\n')
	cont.Store(p)
	for p != nil {
		p, err = buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading buffer: %v\n", err)
			return -1, err
		}
		cont.Store(p)
	}
	return 0, nil
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
		for i, pg := range cont.pgraph {
			text := &font.Drawer{
				Dst:  drw,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(FONTSZ, (FONT_H*(i+1))+FONTSZ),
			}
			text.DrawString(pg.lines[0].line)
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
