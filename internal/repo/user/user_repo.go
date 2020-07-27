package user

import (
	"database/sql"
	"time"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
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

// Repo is an interface that represents the list of functions that need to be
// implemented for the User model, by the repo.
type Repo interface {
	GetByEmail(email string) (*model.User, error)
	GetByEmailInclDeleted(email string) (*model.User, error)
	GetByEmpID(empID string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Insert(u *model.User) (lastInsertID int, err error)
	UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error
	DeleteByEmpID(empID string) (rowsAffected int64, err error)
}

// Impl is an implementation of Repo for Postgres.
type Impl struct {
	db db.DB
}

// New creates a new Repo.
func New(d db.DB) Repo {
	return &Impl{d}
}

// Insert inserts a user using the db connection instance and returns the LastInsertID.
func (i *Impl) Insert(u *model.User) (lastInsertID int, err error) {
	return i.db.Insert(queryInsert, u)
}

// GetByID gets a user from db using the id.
func (i *Impl) GetByID(id int) (*model.User, error) {
	return i.get(queryByID, id)
}

// GetByEmail gets a user from db using the email.
func (i *Impl) GetByEmail(email string) (*model.User, error) {
	return i.get(queryByEmail, email)
}

// GetByEmailInclDeleted gets a user from db using the email, including deleted accounts.
func (i *Impl) GetByEmailInclDeleted(email string) (*model.User, error) {
	return i.get(queryByEmailInclDeleted, email)
}

// GetByEmpID gets a user from db using the empId.
func (i *Impl) GetByEmpID(empID string) (*model.User, error) {
	return i.get(queryByEmpID, empID)
}

// get is a generic getter that other `Get*` methods wrap on.
func (i *Impl) get(query string, args ...interface{}) (*model.User, error) {
	u := model.User{}
	if err := i.db.Get(&u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// UpdatePasswordAndIsActivated updates the user with the specified password and isActivated.
func (i *Impl) UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error {
	u.Password = password
	u.IsActivated = isActivated
	if _, err := i.db.NamedExec(queryUpdatePwdAndIsActivated, u); err != nil {
		return err
	}
	return nil
}

// DeleteByEmpID deletes a user with the empID.
func (i *Impl) DeleteByEmpID(empID string) (rowsAffected int64, err error) {
	result, err := i.db.Exec(queryDeleteByEmpID, time.Now().UTC(), empID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
