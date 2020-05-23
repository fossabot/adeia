package main

import (
	"adeia-api/internal/config"
	log "adeia-api/internal/logger"
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
	err = log.InitLogger(&conf.Logger)
	if err != nil {
		onError("cannot initialize logger", err)
	}

	defer func() {
		_ = log.Sync()
	}()

	apiServer := server.NewAPIServer(conf)
	apiServer.AddRoutes()
	if err := apiServer.Serve(); err != nil {
		log.Panicf("error while serving: %v", err)
	}
}
