package database

import (
	log "adeia-api/internal/utils/logger"
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	config "github.com/spf13/viper"
)

type DB struct {
	*sqlx.DB
}

var (
	// dbConn is the db connection instance.
	dbConn *DB
	// data-source name
	dsn string
	// initConn is used to ensure that dbConn is initialized only once.
	initConn = new(sync.Once)
)

func Init() error {
	err := errors.New("config already loaded")

	initConn.Do(func() {
		err = nil
		dsn = buildDSN()
		c, e := newConn("postgres")
		if e != nil {
			err = e
			return
		}

		dbConn = &DB{c}
	})

	return err
}

func Close() error {
	return dbConn.Close()
}

func GetConn() *DB {
	return dbConn
}

func newConn(driver string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
	}
	log.Debug("Successfully connected to database: " + driver)
	return db, nil
}

func buildDSN() string {
	// helper to get config
	getConfig := func(k string) string {
		return config.GetString("database." + k)
	}

	// get values from config
	params := map[string]string{
		"dbname":   getConfig("dbname"),
		"user":     getConfig("user"),
		"password": getConfig("password"),
		"host":     getConfig("host"),
		"port":     getConfig("port"),
		"sslmode":  getConfig("sslmode"),
	}
	sslParams := map[string]string{
		"sslcert":     getConfig("sslcert"),
		"sslkey":      getConfig("sslkey"),
		"sslrootcert": getConfig("sslrootcert"),
	}

	dsn := ""
	for k, v := range params {
		dsn += k + "='" + v + "'"
	}

	// add ssl params if ssl is enabled
	if params["sslmode"] != "disable" {
		for k, v := range sslParams {
			dsn += k + "='" + v + "'"
		}
	}

	return dsn
}
