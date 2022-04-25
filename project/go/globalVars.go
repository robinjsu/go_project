package main

import (
	"image"
	"image/color"
)

const (
	MAXWIDTH     = 1200
	MAXHEIGHT    = 900
	TEXTWIDTH    = 800
	FONTSZ       = 16
	FONT_REG     = "../../fonts/Karma/Karma-Regular.ttf"
	FONT_BOLD    = "../../fonts/Karma/Karma-Bold.ttf"
	NEWLINE      = byte('\n')
	MAXLINE_TEXT = 125
	MAXLINE_DEF  = 40
	MIN_X_SEARCH = 805
	MAX_Y_SEARCH = 1200
	MIN_X_TEXT   = 0
)

var (
	HIGHLIGHT_GRAY = color.RGBA{0, 0, 255, 100}
	LIGHT_GRAY     = color.RGBA{0, 0, 255, 50}
	PURE_BLUE      = color.RGBA{0, 0, 255, 100}
	TEAL           = color.RGBA{0, 255, 200, 255}
	DEEP_BLUE      = color.RGBA{0, 0, 100, 255}
	MARGIN         = 5
	defCorner      = image.Rect(800, 300, 1200, 1200)
	wordCorner     = image.Rect(800, 0, 1200, 300)
	textBounds     = image.Rect(0, 0, 800, 900)
)
