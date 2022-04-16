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
	"golang.org/x/image/font"
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

func parseText(cont *Content, filename string, face font.Face) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file! %v\n", err)
		return -1, err
	}
	cont.fullText = content
	// TODO: refactor to simplify
	cont.formatLines()
	buffer := bytes.NewBuffer(content)
	p, err := buffer.ReadBytes('\n')
	cont.Store(p, &face)
	for p != nil {
		p, err = buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading buffer: %v\n", err)
			return -1, err
		}
		cont.Store(p, &face)
	}
	return 0, nil
}

func loadTxt(face font.Face, cont *Content) func(drw draw.Image) image.Rectangle {
	load := func(drw draw.Image) image.Rectangle {
		// coordinates refer to the destination image's coordinate space
		page := image.Rect(0, 0, 900, 600)
		draw.Draw(drw, page, image.White, page.Min, draw.Src)
		for i, lns := range cont.format {
			text := &font.Drawer{
				Dst:  drw,
				Src:  image.Black,
				Face: face,
				Dot:  fixed.P(FONTSZ, (FONT_H*(i+1))+FONTSZ),
			}
			bounds, _ := font.BoundString(face, lns.txt)
			lns.bounds = bounds
			fmt.Print(bounds)
			text.DrawString(lns.txt)
		}
		return page
	}
	return load
}

func Text(env gui.Env, textFile string) {
	fontFaces := loadFonts(FONT_REG, FONT_BOLD)
	cont := NewContent()
	_, err := parseText(cont, textFile, fontFaces["regular"])
	if err != nil {
		error.Error(err)
		// panic("panic! text file not properly loaded")
	}
	loadText := loadTxt(fontFaces["regular"], cont)
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
			case win.MoUp:
				loadText = loadTxt(fontFaces["bold"], cont)
				env.Draw() <- loadText
			}
		}
	}
}
