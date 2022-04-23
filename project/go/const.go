package main

import (
	"image"
	"image/color"
)

var (
	HIGHLIGHT_GRAY = color.RGBA{0, 0, 255, 100}
	LIGHT_GRAY     = color.RGBA{0, 0, 255, 50}
	PURE_BLUE      = color.RGBA{0, 0, 255, 100}
	TEAL           = color.RGBA{0, 150, 100, 255}
	DEEP_BLUE      = color.RGBA{0, 0, 100, 255}
	MARGIN         = 5
	defCorner      = image.Rect(900, 300, 1200, 1200)
	wordCorner     = image.Rect(900, 0, 1200, 300)
	textBounds     = image.Rect(0, 0, 900, 900)
)
