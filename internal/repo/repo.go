package repo

import "adeia-api/internal/model"

// UserRepo is an interface that represents the list of functions that need to be
// implemented for the User model, by the repo.
type UserRepo interface {
	GetByEmail(email string) (*model.User, error)
	GetByEmpID(empID string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Insert(u *model.User) (lastInsertID int, err error)
	UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error
	DeleteByEmpID(empID string) (rowsAffected int64, err error)
}
