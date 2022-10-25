package server

import (
	"fetch-me-if-you-read-me/imaginer"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
)

type ServerConfs struct {
	Host string
	Port string
}

type Server struct {
	mux.Router
	listenString string
}

func New(confs *ServerConfs, imaginer *imaginer.Imaginer) (*Server, error) {

	listenString := fmt.Sprintf("%s:%s", confs.Host, confs.Port)
	router := &Server{
		*mux.NewRouter(),
		listenString,
	}

	createImage := newImagesCreate(imaginer)
	imageGet := newImagesGet(imaginer)

	router.Path("/images").
		Methods("POST").
		HandlerFunc(createImage.createImage)

	router.Path("/images/{uuid}").
		Methods("GET").
		HandlerFunc(imageGet.imageGet)

	return router, nil
}

func (server *Server) Listen() error {
	err := http.ListenAndServe(server.listenString, server)
	if err != nil {
		return err
	}
	return nil
}
