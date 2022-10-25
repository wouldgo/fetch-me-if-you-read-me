package imager

import (
	"image"
	"image/color"

	"github.com/google/uuid"
)

var CyanColor = color.RGBA{100, 200, 200, 0xff}

type ImaginerConfs struct {
	Color  *color.RGBA
	Width  uint
	Height uint
}

type Imaginer struct {
	color  color.RGBA
	width  uint
	height uint
}

type Image struct {
	Image *image.RGBA
	Id    string
}

func New(conf *ImaginerConfs) (*Imaginer, error) {
	color := CyanColor
	var width, height uint = 1, 1

	if conf.Color != nil {

		color = *conf.Color
	}

	if conf.Width != 0 {

		width = conf.Width
	}

	if conf.Height != 0 {
		height = conf.Height
	}

	return &Imaginer{
		color:  color,
		width:  uint(width),
		height: uint(height),
	}, nil
}

func (imager *Imaginer) MakeImage() *Image {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{int(imager.width), int(imager.height)}

	imageBounds := image.Rectangle{upLeft, lowRight}
	img := image.NewRGBA(imageBounds)

	for x := 0; x < imageBounds.Dx(); x++ {
		for y := 0; y < imageBounds.Dy(); y++ {
			img.Set(x, y, imager.color)
		}
	}

	return &Image{
		Image: img,
		Id:    uuid.NewString(),
	}
}
