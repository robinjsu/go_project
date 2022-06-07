package main

import (
	"bufio"
	"bytes"
	"image"
	"io"
	"os"
	"strings"

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
	fullText  []byte
	wrapped   []string
	formatted []Formatted
}

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c *Content) parseText(filename string, face font.Face, areaR *image.Rectangle) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return -1, FileError{
			filename: filename,
			err:      err,
		}
	}
	c.fullText = content
	maxWidth := calculateLineWidth(face, areaR.Dx())
	c.wrapped, c.formatted, err = formatLines(c.fullText, maxWidth)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

func formatLines(fullText []byte, maxLineW int) ([]string, []Formatted, error) {
	var lines []string
	var fmtLines []Formatted
	var p []byte
	var lookAhead []byte
	var idx int
	var err error

	buffer := bufio.NewReaderSize(bytes.NewBuffer(fullText), len(fullText))
	length := len(fullText)
	if length < maxLineW {
		lookAhead, err = buffer.Peek(length)
	} else {
		lookAhead, err = buffer.Peek(maxLineW)
	}
	if err != nil {
		return nil, nil, ReadError{err}
	}
	idx = findWrapIdx(lookAhead, maxLineW)
	p = make([]byte, idx, idx)
	n, err := buffer.Read(p)
	if err != nil {
		return nil, nil, ReadError{err}
	}
	ptrim := strings.TrimSuffix(string(p), "\n")
	fmtLines = append(fmtLines, Formatted{txt: ptrim})
	lines = append(lines, ptrim)

	for n != 0 {
		length = buffer.Buffered()
		if length < maxLineW {
			lookAhead, err = buffer.Peek(length)
		} else {
			lookAhead, err = buffer.Peek(maxLineW)
		}
		if err != nil && err != io.EOF {
			return nil, nil, ReadError{err}
		}
		idx = findWrapIdx(lookAhead, maxLineW)
		p = make([]byte, idx, idx)
		n, err = buffer.Read(p)
		if err != nil && err != io.EOF {
			return nil, nil, ReadError{err}
		}
		if n > 0 {
			ptrim := strings.TrimSuffix(string(p), "\n")
			fmtLines = append(fmtLines, Formatted{txt: ptrim})
			lines = append(lines, ptrim)
		}
	}
	return lines, fmtLines, nil
}

func findWrapIdx(b []byte, maxWidth int) int {
	if bytes.ContainsRune(b, rune('\n')) {
		return (bytes.IndexRune(b, rune('\n')) + 1)
	}
	return bytes.LastIndexAny(b, " ") + 1
}

func splitStr(lookup string) []string {
	var list []string
	splitWords := strings.Split(lookup, " ")
	for _, wd := range splitWords {
		word := strings.Trim(wd, " ,.!?';:“”’\"()")
		list = append(list, word)
	}
	return list
}

func parseFont(file string) (*tt.Font, error) {
	ttfFile, err := os.ReadFile(file)
	if err != nil {
		return nil, FileError{
			filename: file,
			err:      err,
		}
	}
	ttf, err := tt.Parse(ttfFile)
	if err != nil {
		return nil, FileError{
			filename: "ttf file bytes",
			err:      err,
		}
	}

	return ttf, nil
}

func loadFonts(fontSize float64, fonts ...string) (map[string]font.Face, error) {
	fontFaces := make(map[string]font.Face)
	for _, f := range fonts {
		// parse bytes and return a pointer to a Font type object
		fnt, err := parseFont(f)
		if err != nil {
			return nil, err
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
	return fontFaces, nil
}

func copyFonts(fontFaces map[string]font.Face) map[string]font.Face {
	return fontFaces
}

func wrapDef(s string, wrapIdx int) []string {
	var lines []string
	if len(s) < wrapIdx {
		return append(lines, s)
	}
	for len(s) > 0 {
		if len(s) < wrapIdx {
			lines = append(lines, s)
			break
		} else {
			idx := strings.LastIndexAny(s[:wrapIdx], " ")
			lines = append(lines, s[:idx])
			s = s[idx:]
		}
	}
	return lines
}

func makePages(content []string, linesPerPage int) [][]string {
	var pages [][]string
	contentCopy := content
	for len(contentCopy) > 0 {
		if len(contentCopy) < linesPerPage {
			pages = append(pages, contentCopy)
			break
		}
		pages = append(pages, contentCopy[:linesPerPage])
		contentCopy = contentCopy[linesPerPage:]
	}
	return pages
}

func calculateLineWidth(face font.Face, contentArea int) int {
	charWidth := font.MeasureString(face, "A")
	return int(contentArea / charWidth.Ceil())
}
