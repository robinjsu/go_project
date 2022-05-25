package main

import (
	"image"
	"image/color"
)

const (
	MAXWIDTH       = 1200
	MAXHEIGHT      = 1000
	FONTSZ         = 16
	FONT_REG       = "../../fonts/Anonymous_Pro/AnonymousPro-Regular.ttf"
	FONT_BOLD      = "../../fonts/Anonymous_Pro/AnonymousPro-Bold.ttf"
	NEWLINE        = byte('\n')
	MAXLINE_DEF    = 50
	MIN_X_SEARCH   = 800
	MAX_X_SEARCH   = 1200
	MIN_X_TEXT     = 0
	MAX_X_TEXT     = 800
	MIN_X_DEF      = 800
	MIN_Y_DEF      = 55
	LINES_PER_PAGE = 26
)

var (
	BLACK         = color.RGBA{13, 19, 33, 255}
	PRUSSIAN_BLUE = color.RGBA{29, 45, 68, 255}
	DARK_BLUE     = color.RGBA{62, 92, 118, 255}
	// SHADOW           = color.RGBA{116, 140, 171, 255}
	EGGSHELL         = color.RGBA{240, 235, 216, 255}
	HIGHLIGHT_GRAY   = color.RGBA{0, 0, 255, 75}
	LIGHT_GRAY       = color.RGBA{0, 0, 255, 50}
	PURE_BLUE        = color.RGBA{0, 0, 255, 100}
	SHADOW           = color.RGBA{116, 140, 171, 255}
	DEEP_BLUE        = color.RGBA{0, 0, 100, 255}
	PRUSSIAN_BLUE_TP = color.RGBA{29, 45, 68, 100}
	DARK_BLUE_TP     = color.RGBA{62, 92, 118, 100}
	SHADOW_TP        = color.RGBA{116, 140, 171, 100}
	EGGSHELL_TP      = color.RGBA{240, 235, 216, 100}
	DEEP_BLUE_TP     = color.RGBA{0, 0, 100, 100}
	MARGIN           = 5
	defCorner        = image.Rect(800, 50, 1200, 1200)
	wordCorner       = image.Rect(800, 0, 1200, 50)
	textBounds       = image.Rect(0, 0, 800, 900)
	buttonR          = image.Rect(0, 0, 75, 30)
	footer           = image.Rect(0, 900, 800, 1000)
)
