package main

import (
	"adeia-api/internal/config"
	"adeia-api/internal/server"
	"adeia-api/internal/service/cache"
	"adeia-api/internal/service/db"
	log "adeia-api/internal/utils/logger"
	"fmt"
	"os"
)

func onError(msg string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, msg+": %v", err)
	os.Exit(1)
}

func main() {
	// load config
	err := config.LoadConf()
	if err != nil {
		onError("cannot load config", err)
	}

	// init logger
	err = log.InitLogger()
	if err != nil {
		onError("cannot initialize logger", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	// init cache
	err = cache.Init()
	if err != nil {
		log.Warnf("cannot initialize cache: %v\nrunning in cache-less mode", err)
	}

	// init db connection
	err = db.Init()
	if err != nil {
		log.Panicf("cannot initialize connection to db: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	// start serving
	apiServer := server.NewAPIServer()
	apiServer.AddRoutes()
	if err := apiServer.Serve(); err != nil {
		log.Panicf("error while serving: %v", err)
	}
}
