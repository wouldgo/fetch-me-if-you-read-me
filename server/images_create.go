package server

import (
	"fetch-me-if-you-read-me/imaginer"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	id := uuid.New()
	newLocation := fmt.Sprintf("images/%s", id.String())

	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(""))
}
