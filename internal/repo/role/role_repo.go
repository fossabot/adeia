package role

import (
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
)

const (
	queryInsert       = "INSERT INTO roles (name) VALUES (:name) RETURNING id"
	queryByID         = "SELECT * FROM roles WHERE id=$1"
	queryByName       = "SELECT * FROM roles WHERE name=$1"
	queryIfNameExists = "SELECT * FROM roles WHERE name=$1 AND id!=$2"
	queryUpdateName   = "UPDATE roles SET name=$1 WHERE id=$2"
)

type Repo struct {
	db repo.DB
}

// New creates a new Repo.
func New(d repo.DB) *Repo {
	return &Repo{d}
}

// Insert inserts a new role into the database.
func (r *Repo) Insert(role *model.Role) (lastInsertID int, err error) {
	return r.db.InsertNamed(queryInsert, role)
}

// GetByName gets a role by its name.
func (r *Repo) GetByName(name string) (*model.Role, error) {
	return r.get(queryByName, name)
}

// CheckIfNameExists checks if name already exists except in the row identified
// by id. This is useful to ensure that the update route (PUT) remains idempotent.
func (r *Repo) CheckIfNameExists(name string, id int) (*model.Role, error) {
	return r.get(queryIfNameExists, name, id)
}

// GetByID gets a role by its id.
func (r *Repo) GetByID(id int) (*model.Role, error) {
	return r.get(queryByID, id)
}

// UpdateName updates the name of a role.
func (r *Repo) UpdateName(roleID int, name string) (rowsAffected int64, err error) {
	return r.db.Update(queryUpdateName, name, roleID)
}

func (r *Repo) get(query string, args ...interface{}) (*model.Role, error) {
	var role model.Role
	if ok, err := r.db.GetOne(&role, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &role, nil
}

func (r *Repo) getMany(query string, args ...interface{}) ([]*model.Role, error) {
	var roles []*model.Role
	if ok, err := r.db.GetMany(&roles, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return roles, nil
}
