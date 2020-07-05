package repo

import (
	"adeia-api/internal/model"
)

// UserRepo is an interface that represents the list of functions that need to be
// implemented for the User model, by the repo.
type UserRepo interface {
	GetByEmail(email string) (*model.User, error)
	GetByEmailInclDeleted(email string) (*model.User, error)
	GetByEmpID(empID string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Insert(u *model.User) (lastInsertID int, err error)
	UpdatePasswordAndIsActivated(u *model.User, password string, isActivated bool) error
	DeleteByEmpID(empID string) (rowsAffected int64, err error)
}

type HolidayRepo interface {
	GetByEpoch(epoch int64) ([]*model.Holiday, error)
	GetByYear(year int) ([]*model.Holiday, error)
	GetByYearAndMonth(year, month int) ([]*model.Holiday, error)
	GetByYMD(year, month, date int) ([]*model.Holiday, error)
	GetByID(id int) (*model.Holiday, error)
	Insert(u *model.Holiday) (int, error)
}
