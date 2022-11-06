package main

import (
	"fetch-me-if-you-read-me/imaginer"
	"fetch-me-if-you-read-me/model"
	"fetch-me-if-you-read-me/server"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/breml/rootcerts"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	options, err := parseOptions()
	if err != nil {

		panic(err)
	}

	imaginer, imaginerErr := imaginer.New(options.Imaginer)
	if imaginerErr != nil {

		panic(imaginerErr)
	}

	model, modelErr := model.New(options.Logger, options.PostgresqlConfigurations)
	if modelErr != nil {

		panic(modelErr)
	}
	defer model.Dispose()

	httpServer, httpServerError := server.New(options.Server, options.Logger, imaginer, model)
	if httpServerError != nil {

		panic(httpServerError)
	}

	go httpServer.Listen()
	signal := <-stop
	options.Logger.Log.Infof("Stopping due to %s", signal.String())
}
