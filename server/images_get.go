package server

import (
	"bytes"
	"fetch-me-if-you-read-me/imaginer"
	logging "fetch-me-if-you-read-me/logger"
	"fetch-me-if-you-read-me/model"

	"image/jpeg"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type imagesGet struct {
	logger   *zap.SugaredLogger
	imaginer *imaginer.Imaginer
	model    *model.Model
}

func (c *imagesGet) imageGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageFk := vars["uuid"]

	imageFkUUID, err := uuid.Parse(imageFk)
	if err != nil {
		c.logger.Errorf(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	image := c.imaginer.MakeImage()
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, image.Image, &jpeg.Options{
		Quality: 1,
	})

	if err != nil {
		c.logger.Errorf(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())

	meta := make(map[string][]string)

	for key, value := range r.Header {
		meta[key] = value
	}

	for key, value := range r.Trailer {
		meta[key] = value
	}

	c.model.ImageFetched(imageFkUUID, r.RemoteAddr, meta)
}

func newImagesGet(logger *logging.Logger, imaginer *imaginer.Imaginer, model *model.Model) *imagesGet {
	return &imagesGet{
		logger:   logger.Log,
		imaginer: imaginer,
		model:    model,
	}
}
