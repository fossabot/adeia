package user

import (
	"time"

	"adeia/internal/model"
	"adeia/internal/repo"
)

const (
	queryInsert = "INSERT INTO users (employee_id, name, email, password, designation, " +
		"is_activated) VALUES (:employee_id, :name, :email, :password, " +
		":designation, :is_activated) RETURNING id"
	queryUpdatePwdAndIsActivated = "UPDATE users SET password=:password, is_activated=:is_activated " +
		"WHERE id=:id AND deleted_at IS NULL"
	queryByID               = "SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL"
	queryByEmail            = "SELECT * FROM users WHERE email=$1 AND deleted_at IS NULL"
	queryByEmailInclDeleted = "SELECT * FROM users WHERE email=$1"
	queryByEmpID            = "SELECT * FROM users WHERE employee_id=$1 AND deleted_at IS NULL"
	queryDeleteByEmpID      = "UPDATE users SET deleted_at=$1 WHERE employee_id=$2 AND deleted_at IS NULL"
)

// Repo is an implementation of Repo for Postgres.
type Repo struct {
	db repo.DB
}

// New creates a new Repo.
func New(d repo.DB) *Repo {
	return &Repo{d}
}

// Insert inserts a user using the db connection instance and returns the LastInsertID.
func (r *Repo) Insert(u *model.User) (lastInsertID int, err error) {
	return r.db.InsertNamed(queryInsert, u)
}

// GetByID gets a user from db using the id.
func (r *Repo) GetByID(id int) (*model.User, error) {
	return r.get(queryByID, id)
}

// GetByEmail gets a user from db using the email.
func (r *Repo) GetByEmail(email string) (*model.User, error) {
	return r.get(queryByEmail, email)
}

// GetByEmailInclDeleted gets a user from db using the email, including deleted accounts.
func (r *Repo) GetByEmailInclDeleted(email string) (*model.User, error) {
	return r.get(queryByEmailInclDeleted, email)
}

// GetByEmpID gets a user from db using the empId.
func (r *Repo) GetByEmpID(empID string) (*model.User, error) {
	return r.get(queryByEmpID, empID)
}

// UpdatePasswordAndIsActivated updates the user with the specified password and isActivated.
func (r *Repo) UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error {
	u.Password = password
	u.IsActivated = isActivated
	if _, err := r.db.UpdateNamed(queryUpdatePwdAndIsActivated, u); err != nil {
		return err
	}
	return nil
}

// DeleteByEmpID deletes a user with the empID.
func (r *Repo) DeleteByEmpID(empID string) (rowsAffected int64, err error) {
	return r.db.Delete(queryDeleteByEmpID, time.Now().UTC(), empID)
}

// get is a generic getter that other `Get*` methods wrap on.
func (r *Repo) get(query string, args ...interface{}) (*model.User, error) {
	u := model.User{}
	if ok, err := r.db.GetOne(&u, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &u, nil
}

func (r *Repo) getMany(query string, args ...interface{}) ([]*model.User, error) {
	var u []*model.User
	if ok, err := r.db.GetMany(&u, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return u, nil
}
