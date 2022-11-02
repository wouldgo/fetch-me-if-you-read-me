package server

import (
	"bytes"
	"fetch-me-if-you-read-me/imaginer"
	"fmt"
	"image/jpeg"
	"net/http"

	"github.com/gorilla/mux"
)

type imagesGet struct {
	imaginer *imaginer.Imaginer
}

func newImagesGet(imaginer *imaginer.Imaginer) *imagesGet {
	return &imagesGet{
		imaginer,
	}
}

func (c *imagesGet) imageGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageUUID := vars["uuid"]

	w.WriteHeader(http.StatusOK)

	image := c.imaginer.MakeImage()
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image.Image, &jpeg.Options{
		Quality: 1,
	})

	if err != nil {
		panic(err.Error())
	}
	fmt.Print(imageUUID)
	w.Write(buf.Bytes())
}
