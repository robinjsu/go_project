package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"

	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Sentence struct {
	line     string
	numBytes int
	advance  fixed.Int26_6
}
type Pgraph struct {
	// single paragraph, split into slice of strings split at each space
	lines []Sentence
	// number of bytes in the entire string
	num int
}
type Content struct {
	fullText []byte
	format   []Formatted
	pgraph   []Pgraph
	num      int
}

func NewContent() *Content {
	c := Content{}
	return &c
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

type Formatted struct {
	txt    string
	span   int
	bounds fixed.Rectangle26_6
}

func endsInSpace(lookAhead []byte) bool {
	lastChar := lookAhead[len(lookAhead)-1]
	secondToLastChar := lookAhead[len(lookAhead)-2]
	return unicode.IsSpace(rune(lastChar)) && unicode.IsSpace(rune(secondToLastChar))
}

func (c *Content) formatLines() []Formatted {
	var fmtLines []Formatted
	var p []byte
	var idx int

	buffer := bufio.NewReader(bytes.NewBuffer(c.fullText))
	lookAhead, err := buffer.Peek(maxLineW + 1)
	if err != nil {
		panic(err)
	}
	if !endsInSpace(lookAhead) {
		idx = bytes.LastIndexAny(lookAhead, " ")
	} else {
		idx = maxLineW
	}
	p = make([]byte, idx, idx)
	n, err := buffer.Read(p)
	if err != nil {
		panic(err)
	}
	fmtLines = append(fmtLines, Formatted{txt: string(lookAhead[0:idx]), span: idx})

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
		if !endsInSpace(lookAhead) {
			idx = bytes.LastIndex(lookAhead, []byte(" "))
		} else {
			idx = maxLineW
		}
		p = make([]byte, idx, idx)
		n, err = buffer.Read(p)
		if err != nil && err != io.EOF {
			panic(err)
		}
		fmtLines = append(fmtLines, Formatted{txt: string(p), span: idx})
	}
	c.format = fmtLines
	return fmtLines
}
