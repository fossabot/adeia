package db

import (
	"database/sql"
	"strconv"

	"adeia-api/internal/config"

	"github.com/jmoiron/sqlx"
)

// PostgresDB is a repo.DB implementation.
type PostgresDB struct {
	*sqlx.DB
}

// New creates a new db connection instance using values from the config.
func New(conf *config.DBConfig) (*PostgresDB, error) {
	d, err := sqlx.Connect(conf.Driver, buildDSN(conf))
	if err != nil {
		return nil, err
	}

	return &PostgresDB{d}, nil
}

// buildDSN is a helper to build the DSN string.
func buildDSN(conf *config.DBConfig) string {
	// get values from config
	params := map[string]string{
		"dbname":   conf.DBName,
		"user":     conf.User,
		"password": conf.Password,
		"host":     conf.Host,
		"port":     strconv.Itoa(conf.Port),
		"sslmode":  conf.SSLMode,
	}
	sslParams := map[string]string{
		"sslcert":     conf.SSLCert,
		"sslkey":      conf.SSLKey,
		"sslrootcert": conf.SSLRootCert,
	}

	dsn := ""
	for k, v := range params {
		dsn += k + "='" + v + "'"
	}

	// add sslParams to dsn if ssl is enabled
	if params["sslmode"] != "disable" {
		for k, v := range sslParams {
			dsn += k + "='" + v + "'"
		}
	}

	return dsn
}

func (p *PostgresDB) Insert(query string, args ...interface{}) (lastInsertID int, err error) {
	err = p.QueryRow(query, args...).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// InsertNamed is a generic database insert that returns the last-insert ID for Postgres.
// This method only accepts named queries.
func (p *PostgresDB) InsertNamed(namedQuery string, arg interface{}) (lastInsertID int, err error) {
	query, args, err := sqlx.Named(namedQuery, arg)
	if err != nil {
		return 0, err
	}

	query = p.Rebind(query)
	err = p.Get(&lastInsertID, query, args...)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// GetMany is a generic database select that returns multiple records.
func (p *PostgresDB) GetMany(dest interface{}, query string, args ...interface{}) (ok bool, err error) {
	err = p.Select(dest, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetOne is a generic database select that returns a single record.
func (p *PostgresDB) GetOne(dest interface{}, query string, args ...interface{}) (ok bool, err error) {
	if err := p.Get(dest, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Update is a generic database update. It wraps the exec method.
func (p *PostgresDB) Update(query string, args ...interface{}) (rowsAffected int64, err error) {
	return p.exec(query, args...)
}

// UpdateNamed is named-query variant of Update. It wraps the execNamed method.
func (p *PostgresDB) UpdateNamed(query string, arg interface{}) (rowsAffected int64, err error) {
	return p.execNamed(query, arg)
}

// Delete is a generic database delete. It wraps the exec method.
func (p *PostgresDB) Delete(query string, args ...interface{}) (rowsAffected int64, err error) {
	return p.exec(query, args...)
}

func (p *PostgresDB) exec(query string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := p.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (p *PostgresDB) execNamed(namedQuery string, arg interface{}) (rowsAffected int64, err error) {
	result, err := p.NamedExec(namedQuery, arg)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
