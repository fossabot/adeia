package repository

import (
	"adeia-api/internal/model"
	"adeia-api/internal/service/db"

	"github.com/jmoiron/sqlx"
)

const (
	queryUserInsert  = "INSERT INTO users (employee_id, name, email) VALUES (:employee_id, :name, :email)"
	queryUserByID    = "SELECT * FROM users WHERE id=$1"
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

// InsertWithTx inserts a user using the provided transaction.
func (p *UserRepoImpl) InsertWithTx(tx *sqlx.Tx, u *model.User) error {
	_, err := tx.NamedExec(queryUserInsert, u)
	return err
}

// Insert inserts a user using the db connection instance.
func (p *UserRepoImpl) Insert(u *model.User) error {
	_, err := p.db.NamedExec(queryUserInsert, u)
	return err
}

// GetByID gets a user from db using the id.
func (p *UserRepoImpl) GetByID(id int) (*model.User, error) {
	u := model.User{}
	if err := p.db.Get(&u, queryUserByID, id); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmpID gets a user from db using the empId.
func (p *UserRepoImpl) GetByEmpID(empID string) (*model.User, error) {
	u := model.User{}
	if err := p.db.Get(&u, queryUserByEmpID, empID); err != nil {
		return nil, err
	}
	return &u, nil
}
