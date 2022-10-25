package imager

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateImage(t *testing.T) {
	imager, err := New(&ImaginerConfs{})

	assert.NotNil(t, imager, "New client has to be not nil")
	assert.Nil(t, err, "Error has to be nil")

	image := imager.MakeImage()

	assert.NotNil(t, image, "Image has to be not nil")

	assert.Equal(t, image.Image.Rect.Dx(), 1, "Image has to be 1 px width")
	assert.Equal(t, image.Image.Rect.Dy(), 1, "Image has to be 1 px height")
}

func TestCreateImageWithConf(t *testing.T) {
	imager, err := New(&ImaginerConfs{
		Color:  &color.RGBA{200, 100, 100, 0xff},
		Width:  10,
		Height: 5,
	})

	assert.NotNil(t, imager, "New client has to be not nil")
	assert.Nil(t, err, "Error has to be nil")

	assert.Equal(t, imager.color.R, uint8(200), "Red value has to be 200")
	assert.Equal(t, imager.color.G, uint8(100), "Green value has to be 100")
	assert.Equal(t, imager.color.B, uint8(100), "Blu value has to be 100")

	image := imager.MakeImage()

	assert.NotNil(t, image, "Image has to be not nil")

	assert.Equal(t, image.Image.Rect.Dx(), 10, "Image has to be 10 px width")
	assert.Equal(t, image.Image.Rect.Dy(), 5, "Image has to be 5 px height")
}
