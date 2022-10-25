package main

import (
	"fetch-me-if-you-read-me/imaginer"
	"fetch-me-if-you-read-me/server"

	_ "github.com/breml/rootcerts"
)

func main() {
	options, err := parseOptions()
	if err != nil {

		panic(err)
	}

	imaginer, imaginerErr := imaginer.New(options.Imaginer)

	if imaginerErr != nil {

		panic(imaginerErr)
	}

	httpServer, httpServerError := server.New(options.Server, imaginer)
	if httpServerError != nil {

		panic(httpServerError)
	}

	httpServer.Listen()
}
