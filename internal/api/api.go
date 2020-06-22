package api

import (
	"fmt"

	"adeia-api/internal/api/server"
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/util/log"
)

// Start starts the API server.
//
// Error handling when serving requests is handled by the `server`. Only other
// panic-able errors (errors that happen on things that are absolutely
// necessary) are returned to main and panic-ed.
func Start() error {
	// init logger
	if err := log.Init(); err != nil {
		return fmt.Errorf("cannot initialize logger: %v", err)
	}
	defer log.Sync()
	log.Debug("successfully initialized logger")

	// init db connection
	dbConn, err := db.New()
	if err != nil {
		return fmt.Errorf("cannot initialize connection to db: %v", err)
	}
	defer dbConn.Close()
	log.Debug("successfully initialized database connection")

	// init cache
	cacheConn, err := cache.New()
	if err != nil {
		log.Warnf("cannot initialize cache: %v\nrunning in cache-less mode...", err)
	} else {
		log.Debug("successfully initialized cache")
	}
	defer cacheConn.Close()

	s := server.New(dbConn, cacheConn)
	s.AddRoutes()
	// start serving
	s.Serve()

	return nil
}
