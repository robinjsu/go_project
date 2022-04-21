package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Formatted struct {
	txt    string
	span   fixed.Int26_6
	bounds fixed.Rectangle26_6
}
type Content struct {
	fullText []byte
	format   []Formatted
}

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c *Content) parseText(filename string, face font.Face) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file! %v\n", err)
		return -1, err
	}
	c.fullText = content
	c.formatLines()

	return 0, nil
}

func (c *Content) formatLines() []Formatted {
	var fmtLines []Formatted
	var p []byte
	var idx int

	buffer := bufio.NewReaderSize(bytes.NewBuffer(c.fullText), len(c.fullText))
	lookAhead, err := buffer.Peek(maxLineW + 1)
	if err != nil {
		panic(err)
	}
	if bytes.ContainsRune(lookAhead, rune('\n')) {
		idx = bytes.IndexRune(lookAhead, rune('\n')) + 1
	} else if !endsInSpace(lookAhead) {
		idx = bytes.LastIndexAny(lookAhead, " ") + 1
	} else {
		idx = maxLineW
	}
	p = make([]byte, idx, idx)
	n, err := buffer.Read(p)
	if err != nil {
		panic(err)
	}
	fmtLines = append(fmtLines, Formatted{txt: string(lookAhead[0:idx])})

	for n != 0 {
		length := buffer.Buffered()
		if length < maxLineW {
			lookAhead, err = buffer.Peek(length)
		} else {
			lookAhead, err = buffer.Peek(maxLineW + 1)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		if bytes.ContainsRune(lookAhead, rune('\n')) {
			idx = bytes.IndexRune(lookAhead, rune('\n')) + 1
		} else if !endsInSpace(lookAhead) {
			idx = bytes.LastIndexAny(lookAhead, " ") + 1
		} else {
			idx = maxLineW
		}
		p = make([]byte, idx, idx)
		n, err = buffer.Read(p)
		if err != nil && err != io.EOF {
			panic(err)
		}
		ptrim := strings.TrimSuffix(string(p), "\n")
		fmtLines = append(fmtLines, Formatted{txt: ptrim})
	}
	c.format = fmtLines
	return fmtLines
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

func loadFonts(fonts ...string) map[string]font.Face {
	fontFaces := make(map[string]font.Face)
	for _, f := range fonts {
		// parse bytes and return a pointer to a Font type object
		fnt, err := parseFont(f)
		if err != nil {
			error.Error(err)
			panic("panic! TTF file not properly loaded")
		}
		// create face, which provides the `glyph mask images`
		face := tt.NewFace(fnt, &tt.Options{
			// options... here just font size (0 is 12-point default)
			Size: FONTSZ,
		})
		switch {
		case strings.Contains(f, "Regular"):
			fontFaces["regular"] = face
		case strings.Contains(f, "Bold"):
			fontFaces["bold"] = face
		}
	}
	return fontFaces
}

func endsInSpace(lookAhead []byte) bool {
	if len(lookAhead) == 0 {
		return true
	}
	lastChar := lookAhead[len(lookAhead)-1]
	secondToLastChar := lookAhead[len(lookAhead)-2]
	return unicode.IsSpace(rune(lastChar)) && unicode.IsSpace(rune(secondToLastChar))
}
