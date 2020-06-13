package api

import (
	"fmt"

	"adeia-api/internal/server"
	"adeia-api/internal/service/cache"
	"adeia-api/internal/service/db"
	log "adeia-api/internal/utils/logger"
)

// Start starts the API server.
func Start() error {
	// Error handling when serving requests is handled by the server. Only other
	// panic-able errors (errors that happen on things that are absolutely
	// necessary) are returned to main and panicked.

	// init logger
	if err := log.Init(); err != nil {
		return fmt.Errorf("cannot initialize logger: %v", err)
	}
	defer log.Sync()

	// init db connection
	if err := db.Init(); err != nil {
		return fmt.Errorf("cannot initialize connection to db: %v", err)
	}
	defer db.Close()

	// init cache
	if err := cache.Init(); err != nil {
		log.Warnf("cannot initialize cache: %v\nrunning in cache-less mode...", err)
	}
	defer cache.Close()

	s := server.New()
	s.AddRoutes()
	// start serving
	s.Serve()

	return nil
}
