package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"strings"

	"golang.org/x/image/math/fixed"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
)

const (
	MAXWIDTH  = 1200
	TEXTWIDTH = 900
	HEIGHT    = 900
	FONTSZ    = 16
	FONT_REG  = "../../fonts/Karma/Karma-Regular.ttf"
	FONT_BOLD = "../../fonts/Karma/Karma-Bold.ttf"
	FONT_H    = 20
	NEWLINE   = byte('\n')
	maxLineW  = 125
)

// Content contains a buffer with text content; each line of text corresponds to text up until the next newline character is found
// TODO: changing approach to parsing - first take in each paragraph, then have separate function to split each p into smaller parts (maybe even as words) in order to
// figure out a way to track location of each word and how it should be tracked so that each word can be clicked on to search
// TODO: take paragraph breaks into account, and split at newline if exists within the lineMaxW

func (c *Content) Store(pg []byte, face *font.Face) error {
	newPg := Pgraph{}
	line := strings.TrimSuffix(string(pg), "\n")
	sentences := strings.SplitAfter(line, ". ")
	for _, lin := range sentences {
		adv := font.MeasureString(*face, lin)
		sentence := Sentence{line: lin, numBytes: len(lin), advance: adv}
		newPg.lines = append(newPg.lines, sentence)
	}
	newPg.num = len(sentences)
	c.pgraph = append(c.pgraph, newPg)
	c.num++
	return nil
}

func (c *Content) parseText(filename string, face font.Face) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file! %v\n", err)
		return -1, err
	}
	c.fullText = content
	// TODO: refactor to simplify
	c.formatLines()
	buffer := bytes.NewBuffer(content)
	p, err := buffer.ReadBytes('\n')
	c.Store(p, &face)
	for p != nil {
		p, err = buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading buffer: %v\n", err)
			return -1, err
		}
		c.Store(p, &face)
	}
	return 0, nil
}

func (cont *Content) loadTxt(face font.Face) func(drw draw.Image) image.Rectangle {
	load := func(drw draw.Image) image.Rectangle {
		// coordinates refer to the destination image's coordinate space
		page := image.Rect(0, 0, 900, 900)
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for i, lns := range cont.format {
			text := &font.Drawer{
				Dst:  drw,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(FONTSZ, (FONT_H*(i+1))+FONTSZ),
			}
			cont.format[i].bounds, _ = text.BoundString(lns.txt)
			text.DrawString(lns.txt)
		}
		return page
	}
	return load
}

func highlightLine(face font.Face, cont *Content, p image.Point) func(drw draw.Image) image.Rectangle {
	var line image.Rectangle
	load := func(drw draw.Image) image.Rectangle {
		for _, ln := range cont.format {
			rct := ln.bounds
			if p.Y >= (rct.Min.Y).Floor() && p.Y <= (rct.Max.Y).Ceil() && p.X <= TEXTWIDTH {
				line = image.Rect((rct.Min.X).Floor(), (rct.Min.Y).Floor(), (rct.Max.X).Floor(), (rct.Max.Y).Floor())
				draw.Draw(drw, line, &image.Uniform{color.RGBA{0, 0, 255, 100}}, image.ZP, draw.Over)
			}
		}
		return line
	}
	return load
}

func Text(env gui.Env, textFile string) {
	fontFaces := loadFonts(FONT_REG, FONT_BOLD)
	cont := NewContent()
	_, err := cont.parseText(textFile, fontFaces["regular"])
	if err != nil {
		error.Error(err)
		// panic("panic! text file not properly loaded")
	}
	loadText := cont.loadTxt(fontFaces["regular"])
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
				fmt.Println(e.X, e.Y)
			case win.MoUp:
				loadText = cont.loadTxt(fontFaces["regular"])
				highlight := highlightLine(fontFaces["bold"], cont, image.Pt(e.X, e.Y))
				env.Draw() <- loadText
				env.Draw() <- highlight
			}
		}
	}
}
