package server

import (
	"fetch-me-if-you-read-me/imaginer"
	"fmt"
	"net/http"
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
	fmt.Println(c)
	w.WriteHeader(http.StatusOK)
}
