package server

import (
	"fetch-me-if-you-read-me/imaginer"
	"fmt"
	"net/http"
)

type imagesCreate struct {
	imaginer *imaginer.Imaginer
}

func newImagesCreate(imaginer *imaginer.Imaginer) *imagesCreate {
	return &imagesCreate{
		imaginer,
	}
}

func (c *imagesCreate) createImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(c)
	w.WriteHeader(http.StatusOK)
}
