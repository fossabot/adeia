package main

import (
	"adeia-api/internal/config"
	"adeia-api/internal/server"
	"log"
)

func main() {
	conf, err := config.Load("./config/config.yaml")
	if err != nil {
		log.Panicf("cannot load config: %v", err)
	}

	apiServer := server.NewAPIServer(conf)
	apiServer.AddRoutes()
	if err := apiServer.Serve(); err != nil {
		log.Panicf("error while serving: %v", err)
	}
}
