package server

import (
	"fetch-me-if-you-read-me/imaginer"
	logging "fetch-me-if-you-read-me/logger"
	"fetch-me-if-you-read-me/model"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServerConfs struct {
	Host string
	Port string
}

type Server struct {
	mux.Router
	listenString string
	logger       *zap.SugaredLogger
}

func New(confs *ServerConfs, logger *logging.Logger, imaginer *imaginer.Imaginer, model *model.Model) (*Server, error) {
	listenString := fmt.Sprintf("%s:%s", confs.Host, confs.Port)
	router := &Server{
		*mux.NewRouter(),
		listenString,
		logger.Log,
	}

	logger.Log.Debugf("Creating server on %s ...", listenString)
	createImage := newImagesCreate(logger, imaginer, model)
	imageGet := newImagesGet(logger, imaginer, model)
	statusHandlerFunc := newStatus(logger, model)

	router.
		Methods("GET").
		Path("/live").
		HandlerFunc(statusHandlerFunc.statusHandler)

	router.
		Path("/images").
		Methods("POST").
		HandlerFunc(createImage.createImage)

	router.Path("/images/{uuid}").
		Methods("HEAD", "GET", "POST").
		HandlerFunc(imageGet.imageGet)

	return router, nil
}

func (server *Server) Listen() error {
	err := http.ListenAndServe(server.listenString, server)
	if err != nil {
		return err
	}
	server.logger.Infof("Listening on server %s...", server.listenString)
	return nil
}
