package repo

import (
	"database/sql"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
)

const (
	queryUserInsert = "INSERT INTO users (employee_id, name, email, password, designation, " +
		"is_activated) VALUES (:employee_id, :name, :email, :password, :designation, :is_activated) RETURNING id"
	queryUserByID    = "SELECT * FROM users WHERE id=$1"
	queryUserByEmail = "SELECT * FROM users WHERE email=$1"
	queryUserByEmpID = "SELECT * FROM users WHERE employee_id=$1"
)

// UserRepoImpl is an implementation of UserRepo for Postgres.
type UserRepoImpl struct {
	db db.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(d db.DB) UserRepo {
	return &UserRepoImpl{d}
}

// Insert inserts a user using the db connection instance and returns the LastInsertID.
func (i *UserRepoImpl) Insert(u *model.User) (int, error) {
	stmt, err := i.db.PrepareNamed(queryUserInsert)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastInsertID int
	if err := stmt.Get(&lastInsertID, u); err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// GetByID gets a user from db using the id.
func (i *UserRepoImpl) GetByID(id int) (*model.User, error) {
	u := model.User{}
	if err := i.db.Get(&u, queryUserByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByEmail gets a user from db using the email.
func (i *UserRepoImpl) GetByEmail(email string) (*model.User, error) {
	u := model.User{}
	if err := i.db.Get(&u, queryUserByEmail, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByEmpID gets a user from db using the empId.
func (i *UserRepoImpl) GetByEmpID(empID string) (*model.User, error) {
	u := model.User{}
	if err := i.db.Get(&u, queryUserByEmpID, empID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
