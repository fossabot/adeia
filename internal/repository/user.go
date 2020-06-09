package repository

import (
	"adeia-api/internal/model"
	"adeia-api/internal/service/db"

	"github.com/jmoiron/sqlx"
)

const (
	queryUserInsert  = "INSERT INTO users (employee_id, name, email) VALUES (:employee_id, :name, :email)"
	queryUserByEmpID = "SELECT * FROM users WHERE employee_id=$1"
)

// UserRepoImpl is an implementation of UserRepo for Postgres.
type UserRepoImpl struct{}

// NewUserRepo creates a new UserRepoImpl.
func NewUserRepo() UserRepo {
	return &UserRepoImpl{}
}

// InsertWithTx inserts a user using the provided transaction.
func (p *UserRepoImpl) InsertWithTx(tx *sqlx.Tx, u *model.User) error {
	_, err := tx.NamedExec(queryUserInsert, u)
	return err
}

// Insert inserts a user using the db connection instance.
func (p *UserRepoImpl) Insert(u *model.User) error {
	_, err := db.NamedExec(queryUserInsert, u)
	return err
}

// GetByEmpID gets a user from db using the empId.
func (p *UserRepoImpl) GetByEmpID(empID string) (*model.User, error) {
	u := model.User{}
	err := db.Get(&u, queryUserByEmpID, empID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
