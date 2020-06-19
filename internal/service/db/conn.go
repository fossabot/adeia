package db

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // import postgres driver
	config "github.com/spf13/viper"
)

// DB represents the database object.
type DB struct {
	*sqlx.DB
}

var (
	// dbConn is the db connection instance.
	dbConn *DB
	// initConn is used to ensure that dbConn is initialized only once.
	initConn = new(sync.Once)
)

// Init creates a new db connection instance using values from the config.
func Init() error {
	err := errors.New("config already loaded")

	initConn.Do(func() {
		err = nil
		dsn := buildDSN()
		driver := config.GetString("database.driver")
		c, e := newConn(driver, dsn)
		if e != nil {
			err = e
			return
		}

		dbConn = &DB{c}
	})

	return err
}

// newConn creates a new connection instance.
func newConn(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// buildDSN is a helper to build the DSN string.
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

// Close closes the connection.
func Close() error {
	return dbConn.Close()
}

// GetConn returns the connection instance.
func GetConn() *DB {
	return dbConn
}

// SetConn sets the connection instance.
func SetConn(d *DB) {
	dbConn = d
}

// NamedExec wraps sqlx.DB's NamedExec.
func NamedExec(query string, arg interface{}) (sql.Result, error) {
	return dbConn.NamedExec(query, arg)
}

// NamedExecContext wraps sqlx.DB's NamedExecContext.
func NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return dbConn.NamedExecContext(ctx, query, arg)
}

// Get wraps sqlx.DB's Get.
func Get(dest interface{}, query string, args ...interface{}) error {
	return dbConn.Get(dest, query, args...)
}

// GetContext wraps sqlx.DB's GetContext.
func GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return dbConn.GetContext(ctx, dest, query, args...)
}
