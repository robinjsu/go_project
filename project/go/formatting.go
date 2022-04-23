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
	_, c.format = formatLines(c.fullText, MAXLINEWIDTH)

	return 0, nil
}

func formatLines(fullText []byte, maxLineW int) ([]string, []Formatted) {
	var lines []string
	var fmtLines []Formatted
	var p []byte
	var lookAhead []byte
	var idx int
	var err error

	buffer := bufio.NewReaderSize(bytes.NewBuffer(fullText), len(fullText))
	length := buffer.Buffered()
	if length < maxLineW+1 {
		lookAhead, err = buffer.Peek(length)
	} else {
		lookAhead, err = buffer.Peek(maxLineW + 1)
	}
	if err != nil {
		panic(err)
	}
	idx = findWrapIdx(lookAhead, maxLineW)
	p = make([]byte, idx, idx)
	n, err := buffer.Read(p)
	if err != nil {
		panic(err)
	}
	ptrim := strings.TrimSuffix(string(p), "\n")
	fmtLines = append(fmtLines, Formatted{txt: ptrim})
	lines = append(lines, ptrim)

	for n != 0 {
		length = buffer.Buffered()
		if length < maxLineW {
			lookAhead, err = buffer.Peek(length)
		} else {
			lookAhead, err = buffer.Peek(maxLineW + 1)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		idx = findWrapIdx(lookAhead, maxLineW)
		p = make([]byte, idx, idx)
		n, err = buffer.Read(p)
		if err != nil && err != io.EOF {
			panic(err)
		}
		ptrim := strings.TrimSuffix(string(p), "\n")
		fmtLines = append(fmtLines, Formatted{txt: ptrim})
		lines = append(lines, ptrim)
	}

	return lines, fmtLines
}

func findWrapIdx(b []byte, maxWidth int) int {
	if bytes.ContainsRune(b, rune('\n')) {
		return (bytes.IndexRune(b, rune('\n')) + 1)
	} else if !endsInSpace((b)) {
		return (bytes.LastIndexAny(b, " ") + 1)
	}
	return maxWidth
}

func endsInSpace(lookAhead []byte) bool {
	if len(lookAhead) > 1 {
		lastChar := lookAhead[len(lookAhead)-1]
		secondToLastChar := lookAhead[len(lookAhead)-2]
		return unicode.IsSpace(rune(lastChar)) && unicode.IsSpace(rune(secondToLastChar))
	}
	return true
}

func splitStr(lookup string) []string {
	var list []string
	splitWords := strings.Split(lookup, " ")
	for _, wd := range splitWords {
		word := strings.Trim(wd, " ,.!?';:“”’\"()")
		if !isCommon(word) {
			list = append(list, word)
		}
	}
	return list
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

func loadFonts(fontSize float64, fonts ...string) map[string]font.Face {
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
			Size: fontSize,
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

func (word *Word) formatDefs() Word {
	splitIdx := 40
	for i, d := range word.Def {
		s := fmt.Sprintf(" - (%s) %s", d.PartOfSpeech, d.Definition)
		fmtDefs := wrapDef(s, splitIdx)
		word.Def[i].Formatted = fmtDefs
	}
	return *word
}

func wrapDef(s string, wrapIdx int) []string {
	var lines []string
	if len(s) < wrapIdx {
		return append(lines, s)
	}
	for len(s) > 0 {
		if len(s) < 40 {
			lines = append(lines, s)
			break
		} else {
			lines = append(lines, s[:wrapIdx])
			s = s[wrapIdx:]
		}
	}
	return lines
}
