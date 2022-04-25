package main

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type imageObj struct {
	text      Formatted
	img       image.Image
	placement image.Rectangle
}

// https://github.com/faiface/gui/blob/master/examples/imageviewer/util.go#L66
func drawText(s string, face font.Face) (image.Image, Formatted) {
	text := &font.Drawer{
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, face.Metrics().Height.Ceil()),
	}
	txtBnds, txtAdv := text.BoundString(s)
	bounds := image.Rect(
		txtBnds.Min.X.Floor(),
		txtBnds.Min.Y.Floor(),
		txtBnds.Max.X.Ceil(),
		txtBnds.Max.Y.Ceil(),
	)
	text.Dst = image.NewRGBA(bounds)
	text.DrawString(s)
	return text.Dst, Formatted{txt: s, span: txtAdv, bounds: txtBnds}
}

func makeInsetRect(r image.Rectangle, margin int) image.Rectangle {
	x0 := defCorner.Min.X + margin
	y0 := defCorner.Min.Y + margin
	sz := defCorner.Size().Sub(image.Pt(margin, margin))

	return image.Rect(x0, y0, x0+sz.X, y0+sz.Y)
}

func makeHeaderR(r image.Rectangle, header image.Image, margin int) image.Rectangle {
	x0 := r.Bounds().Min.X + margin
	y0 := r.Bounds().Min.Y + margin
	x1 := x0 + header.Bounds().Max.X
	y1 := y0 + header.Bounds().Max.Y

	return image.Rect(x0, y0, x1, y1)
}
