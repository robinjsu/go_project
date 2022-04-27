package main

import (
	"image"
	"image/color"
)

const (
	MAXWIDTH       = 1200
	MAXHEIGHT      = 1000
	FONTSZ         = 16
	FONT_REG       = "../../fonts/Karma/Karma-Regular.ttf"
	FONT_BOLD      = "../../fonts/Karma/Karma-Bold.ttf"
	NEWLINE        = byte('\n')
	MAXLINE_TEXT   = 100
	MAXLINE_DEF    = 50
	MIN_X_SEARCH   = 800
	MAX_X_SEARCH   = 1200
	MIN_X_TEXT     = 0
	MAX_X_TEXT     = 800
	MIN_X_DEF      = 800
	MIN_Y_DEF      = 325
	LINES_PER_PAGE = 26
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
	buttonR        = image.Rect(0, 0, 100, 50)
)
