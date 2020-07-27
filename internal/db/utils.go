package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type txFn func(tx *sqlx.Tx) error

// WithTx is a wrapper that executes a given callback inside a transaction. Any
// error returned by the callback will rollback the transaction. If no error is
// returned, the transaction is committed.
func WithTx(db *sqlx.DB, fn txFn) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if ex := tx.Rollback(); ex != nil {
			return fmt.Errorf("rollback for '%w' failed with: %v", err, ex)
		}
		return err
	}
	return tx.Commit()
}

// Insert is a generic database insert that returns the last-insert ID for Postgres.
// This method only accepts named queries.
func (p *PostgresDB) Insert(namedQuery string, arg interface{}) (lastInsertID int, err error) {
	stmt, err := p.PrepareNamed(namedQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	if err := stmt.Get(&lastInsertID, arg); err != nil {
		return 0, err
	}
	return lastInsertID, nil
}
