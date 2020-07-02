package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib" // import postgres driver
	"github.com/jmoiron/sqlx"
	config "github.com/spf13/viper"
)

// DB represents the funcs required for database access.
type DB interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Close() error
}

// PostgresDB is a DB implementation.
type PostgresDB struct {
	*sqlx.DB
}

// New creates a new db connection instance using values from the config.
func New() (DB, error) {
	dsn := buildDSN()
	d, err := sqlx.Connect(config.GetString("database.driver"), dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{d}, nil
}

// buildDSN is a helper to build the DSN string.
func buildDSN() string {
	getConf := func(k string) string {
		return config.GetString("database." + k)
	}

	// get values from config
	params := map[string]string{
		"dbname":   getConf("dbname"),
		"user":     getConf("user"),
		"password": getConf("password"),
		"host":     getConf("host"),
		"port":     getConf("port"),
		"sslmode":  getConf("sslmode"),
	}
	sslParams := map[string]string{
		"sslcert":     getConf("sslcert"),
		"sslkey":      getConf("sslkey"),
		"sslrootcert": getConf("sslrootcert"),
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
