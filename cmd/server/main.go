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

var termination []os.Signal = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
var is = make([]interface{}, len(termination))

func main() {

	for i, v := range termination {
		is[i] = v
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, termination...)

	options, err := parseOptions()
	if err != nil {

		panic(err)
	}

	imaginer, imaginerErr := imaginer.New(options.Imaginer)
	if imaginerErr != nil {

		panic(imaginerErr)
	}

	options.Logger.Log.Info("Setup model")

	model, modelErr := model.New(options.Logger, options.PostgresqlConfigurations)
	if modelErr != nil {

		panic(modelErr)
	}
	defer model.Dispose()

	options.Logger.Log.Info("Setup http server")
	httpServer, httpServerError := server.New(options.Server, options.Logger, imaginer, model)
	if httpServerError != nil {

		panic(httpServerError)
	}

	go httpServer.Listen()

	options.Logger.Log.Infof("Waiting %+v...", is)
	signal := <-stop
	options.Logger.Log.Infof("Stopping due to %s", signal.String())
}
