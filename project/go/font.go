package main

import (
	"os"
	"strings"

	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

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
