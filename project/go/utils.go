package main

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

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
	txtImg := image.NewRGBA(bounds)
	draw.Draw(txtImg, bounds, image.White, bounds.Min, draw.Src)
	text.Dst = txtImg
	text.DrawString(s)
	return text.Dst, Formatted{txt: s, span: txtAdv, bounds: txtBnds}
}

func formatTextImages(wordList []string, minY int, minX int, margin int, lineHeight int, face font.Face) []imageObj {
	var images []imageObj
	x0 := minX + margin
	y0 := minY + margin
	y1 := minY + margin + lineHeight
	for _, w := range wordList {
		img, format := drawText(w, face)
		x1 := x0 + img.Bounds().Dx()
		fontR := image.Rect(x0, y0, x1, y1)
		images = append(images, imageObj{text: format, img: img, placement: fontR})
		y0 += lineHeight
		y1 += lineHeight
	}
	return images
}

func makeInsetR(r image.Rectangle, margin int) image.Rectangle {
	x0 := r.Min.X + margin
	y0 := r.Min.Y + margin
	sz := r.Size().Sub(image.Pt(margin, margin))

	return image.Rect(x0, y0, r.Min.X+sz.X, r.Min.Y+sz.Y)
}

func makeHeaderLeftR(r image.Rectangle, header image.Image, margin int) image.Rectangle {
	x0 := r.Bounds().Min.X + margin
	y0 := r.Bounds().Min.Y + margin
	x1 := x0 + header.Bounds().Max.X
	y1 := y0 + header.Bounds().Max.Y

	return image.Rect(x0, y0, x1, y1)
}

func drawBtn(text string, face font.Face, r image.Rectangle, loc image.Point) draw.Image {
	rLoc := r.Add(loc)
	button := image.NewRGBA(rLoc)
	txt := &font.Drawer{
		Dst:  button,
		Src:  image.Black,
		Face: face,
	}
	bounds, adv := txt.BoundString(text)
	width := rLoc.Bounds().Dx()
	height := rLoc.Bounds().Dy()
	paddingW := (width - int(adv)) / 2
	paddingH := height - ((height - int(bounds.Max.Y.Ceil())) / 2)
	txt.Dot = fixed.P(paddingW, paddingH)
	txt.DrawString(text)

	return txt.Dst
}

func getPNG(pngFile string) (image.Image, error) {
	if !strings.HasSuffix(pngFile, "png") {
		return nil, FileError{
			filename: pngFile,
			err:      errors.New("wrong file type, must be PNG"),
		}
	}
	pngBytes, err := os.ReadFile(pngFile)
	if err != nil {
		return nil, FileError{
			filename: pngFile,
			err:      err,
		}
	}
	pngReader := bytes.NewReader(pngBytes)
	pngImage, err := png.Decode(pngReader)
	if err != nil {
		return nil, FileError{
			filename: pngFile,
			err:      err,
		}
	}
	return pngImage, nil
}
