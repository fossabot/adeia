package repo

import (
	"adeia-api/internal/model"

	"github.com/jmoiron/sqlx"
)

// UserRepo is an interface that represents the list of functions that need to be
// implemented for the User model, by the repo.
type UserRepo interface {
	Insert(u *model.User) error
	InsertWithTx(tx *sqlx.Tx, u *model.User) error
	GetByEmpID(empID string) (*model.User, error)
	GetByID(id int) (*model.User, error)
}
