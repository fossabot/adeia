package main

import (
	"adeia-api/internal/config"
	"adeia-api/internal/logger"
	"adeia-api/internal/server"
	"fmt"
	"os"
)

func onError(msg string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, msg+": %v", err)
	os.Exit(1)
}

func main() {
	// load config
	conf, err := config.Load("./config/config.yaml")
	if err != nil {
		onError("cannot load config", err)
	}

	// init logger
	err = logger.Init(&conf.Logger)
	if err != nil {
		onError("cannot initialize logger", err)
	}

	log := logger.Get()
	defer log.Sync()

	apiServer := server.NewAPIServer(conf)
	apiServer.AddRoutes()
	if err := apiServer.Serve(); err != nil {
		log.Panicf("error while serving: %v", err)
	}
}
