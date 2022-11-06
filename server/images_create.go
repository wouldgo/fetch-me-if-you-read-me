package server

import (
	"errors"
	"fetch-me-if-you-read-me/imaginer"
	logging "fetch-me-if-you-read-me/logger"
	"fetch-me-if-you-read-me/model"

	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ImageCreation struct {
	UsedIn string
}

func (c *imagesCreate) createImage(w http.ResponseWriter, r *http.Request) {
	var anImageCreation ImageCreation
	err := decodeJSONBody(w, r, &anImageCreation)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			c.logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	uuid, err := c.model.PrepareImage(anImageCreation.UsedIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newLocation := fmt.Sprintf("images/%s", uuid.String())

	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusTemporaryRedirect)

	w.Write([]byte(""))
}

type imagesCreate struct {
	logger   *zap.SugaredLogger
	imaginer *imaginer.Imaginer
	model    *model.Model
}

func newImagesCreate(logger *logging.Logger, imaginer *imaginer.Imaginer, model *model.Model) *imagesCreate {
	return &imagesCreate{
		logger:   logger.Log,
		imaginer: imaginer,
		model:    model,
	}
}
