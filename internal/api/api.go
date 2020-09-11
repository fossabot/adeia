package api

import (
	"fmt"

	"adeia-api/internal/api/server"
	"adeia-api/internal/cache"
	"adeia-api/internal/config"
	"adeia-api/internal/controller"
	"adeia-api/internal/db"
	"adeia-api/internal/mailer"
	"adeia-api/internal/repo"
	"adeia-api/internal/service"
	"adeia-api/internal/util/log"
)

var (
	dbConn         repo.DB
	cacheConn      service.Cache
	mailerInstance service.Mailer
)

// Start starts the API server.
//
// Error handling when serving requests is handled by the `server`. Only other
// panic-able errors (errors that happen on things that are absolutely
// necessary) are returned to main and panic-ed.
func Start(conf *config.Config) error {
	// init all necessary services
	if err := initServices(conf); err != nil {
		return err
	}

	defer func() {
		log.Debug("closing connections...")
		_ = log.Sync()
		_ = dbConn.Close()
		_ = cacheConn.Close()
	}()

	controller.Init(conf, dbConn, cacheConn, mailerInstance)
	s := server.New(&conf.ServerConfig)
	s.AddRoutes()
	// start serving
	s.Serve()

	return nil
}

func initServices(conf *config.Config) (err error) {
	// init logger
	if err = log.Init(&conf.LoggerConfig); err != nil {
		return fmt.Errorf("cannot initialize logger: %v", err)
	}
	log.Debug("successfully initialized logger")

	// init db
	dbConn, err = db.New(&conf.DBConfig)
	if err != nil {
		return fmt.Errorf("cannot initialize connection to db: %v", err)
	}
	log.Debug("successfully initialized database connection")

	// init cache
	cacheConn, err = cache.New(&conf.CacheConfig)
	if err != nil {
		return fmt.Errorf("cannot initialize cache: %v", err)
	}
	log.Debug("successfully initialized cache connection")

	// init mailer
	mailerInstance, err = mailer.New(&conf.MailerConfig)
	if err != nil {
		return fmt.Errorf("cannot initialize mailer: %v", err)
	}

	return nil
}
