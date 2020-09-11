package session

import (
	"time"

	"adeia/internal/model"
	"adeia/internal/repo"
)

const (
	queryInsert = "INSERT INTO sessions (user_id, refresh_token, refresh_token_expires) VALUES " +
		"(:user_id, :refresh_token, :refresh_token_expires) RETURNING id"
	queryByUserIDAndRefreshToken       = "SELECT * FROM sessions WHERE user_id=$1 AND refresh_token=$2"
	queryUpdateRefreshToken            = "UPDATE sessions SET refresh_token=$1, refresh_token_expires=$2 WHERE id=$3"
	queryDeleteByUserIDAndRefreshToken = "DELETE FROM sessions WHERE user_id=$1 AND refresh_token=$2"
)

// Repo is an implementation of Repo for Postgres.
type Repo struct {
	db repo.DB
}

// New creates a new Repo.
func New(d repo.DB) *Repo {
	return &Repo{d}
}

// Insert inserts a new session into the database.
func (r *Repo) Insert(s *model.Session) (lastInsertID int, err error) {
	return r.db.InsertNamed(queryInsert, s)
}

// GetByUserIDAndRefreshToken gets a session identified by its userID and refreshToken.
func (r *Repo) GetByUserIDAndRefreshToken(id int, refreshToken []byte) (*model.Session, error) {
	s := model.Session{}
	if ok, err := r.db.GetOne(&s, queryByUserIDAndRefreshToken, id, refreshToken); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &s, nil
}

// UpdateRefreshToken updates a session with the provided refreshToken and expiry.
func (r *Repo) UpdateRefreshToken(id int, refreshToken []byte, expires time.Time) (rowsAffected int64, err error) {
	return r.db.Update(queryUpdateRefreshToken, refreshToken, expires, id)
}

// DeleteByUserIDAndRefreshToken deletes a session with the provided userID and refreshToken.
func (r *Repo) DeleteByUserIDAndRefreshToken(id int, refreshToken []byte) (rowsAffected int64, err error) {
	return r.db.Delete(queryDeleteByUserIDAndRefreshToken, id, refreshToken)
}
