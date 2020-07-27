package session

import (
	"database/sql"
	"time"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
)

const (
	queryInsert = "INSERT INTO sessions (user_id, refresh_token, refresh_token_expires) VALUES " +
		"(:user_id, :refresh_token, :refresh_token_expires) RETURNING id"
	queryByUserIDAndRefreshToken       = "SELECT * FROM sessions WHERE user_id=$1 AND refresh_token=$2"
	queryUpdateRefreshToken            = "UPDATE sessions SET refresh_token=$1, refresh_token_expires=$2 WHERE id=$3"
	queryDeleteByUserIDAndRefreshToken = "DELETE FROM sessions WHERE user_id=$1 AND refresh_token=$2"
)

// Repo is an interface that represents the list of functions that need to be
// implemented for the Session model, by the repo.
type Repo interface {
	Insert(s *model.Session) (int, error)
	GetByUserIDAndRefreshToken(id int, refreshToken []byte) (*model.Session, error)
	UpdateRefreshToken(id int, refreshToken []byte, expires time.Time) (rowsAffected int64, err error)
	DeleteByUserIDAndRefreshToken(id int, refreshToken []byte) (rowsAffected int64, err error)

	exec(query string, args ...interface{}) (rowsAffected int64, err error)
}

// Impl is an implementation of Repo for Postgres.
type Impl struct {
	db db.DB
}

// New creates a new Repo.
func New(d db.DB) Repo {
	return &Impl{d}
}

// Insert inserts a new session into the database.
func (i *Impl) Insert(s *model.Session) (lastInsertID int, err error) {
	return i.db.Insert(queryInsert, s)
}

// GetByUserIDAndRefreshToken gets a session identified by its userID and refreshToken.
func (i *Impl) GetByUserIDAndRefreshToken(id int, refreshToken []byte) (*model.Session, error) {
	s := model.Session{}
	if err := i.db.Get(&s, queryByUserIDAndRefreshToken, id, refreshToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// UpdateRefreshToken updates a session with the provided refreshToken and expiry.
func (i *Impl) UpdateRefreshToken(id int, refreshToken []byte, expires time.Time) (rowsAffected int64, err error) {
	return i.exec(queryUpdateRefreshToken, refreshToken, expires, id)
}

// DeleteByUserIDAndRefreshToken deletes a session with the provided userID and refreshToken.
func (i *Impl) DeleteByUserIDAndRefreshToken(id int, refreshToken []byte) (rowsAffected int64, err error) {
	return i.exec(queryDeleteByUserIDAndRefreshToken, id, refreshToken)
}

func (i *Impl) exec(query string, args ...interface{}) (rowsAffected int64, err error) {
	result, err := i.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
