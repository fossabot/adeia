package role

import (
	"database/sql"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
)

const (
	queryInsert = "INSERT INTO roles (name) VALUES (:name) RETURNING id"
	queryByID   = "SELECT * FROM roles WHERE id=$1"
	queryByName = "SELECT * FROM roles WHERE name=$1"
)

// Repo is an interface that represents the list of functions that need to be
// implemented for the Role model, by the repo.
type Repo interface {
	Insert(r *model.Role) (lastInsertID int, err error)
	GetByID(id int) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
}

// Impl is an implementation of Repo for Postgres.
type Impl struct {
	db db.DB
}

// New creates a new Repo.
func New(d db.DB) Repo {
	return &Impl{d}
}

// Insert inserts a new role into the database.
func (i *Impl) Insert(r *model.Role) (lastInsertID int, err error) {
	return i.db.Insert(queryInsert, r)
}

// GetByName gets a role by its name.
func (i *Impl) GetByName(name string) (*model.Role, error) {
	return i.get(queryByName, name)
}

// GetByID gets a role by its id.
func (i *Impl) GetByID(id int) (*model.Role, error) {
	return i.get(queryByID, id)
}

func (i *Impl) get(query string, args ...interface{}) (*model.Role, error) {
	var r model.Role
	if err := i.db.Get(&r, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}
