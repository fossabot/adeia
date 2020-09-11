package service

import (
	"time"

	"adeia-api/internal/model"

	"github.com/jordan-wright/email"
)

type Cache interface {
	Close() error
	Delete(keys ...string) error
	Expire(key string, seconds int) error
	Get(rcv interface{}, key string) error
	Set(key string, value string) error
	SetWithExpiry(key, value string, seconds int) error
}

// Mailer represents a mailer service.
type Mailer interface {
	Send(e *email.Email, template string, data interface{}) error
}

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

// SessionRepo is an interface that represents the list of functions that need to be
// implemented for the Session model, by the repo.
type SessionRepo interface {
	Insert(s *model.Session) (lastInsertID int, err error)
	GetByUserIDAndRefreshToken(id int, refreshToken []byte) (*model.Session, error)
	UpdateRefreshToken(id int, refreshToken []byte, expires time.Time) (rowsAffected int64, err error)
	DeleteByUserIDAndRefreshToken(id int, refreshToken []byte) (rowsAffected int64, err error)
}

// HolidayRepo is an interface that represents the list of functions that need to be
// implemented for the Holiday model, by the repo.
type HolidayRepo interface {
	DeletedByID(id int) (int64, error)
	GetByEpoch(epoch int64) ([]*model.Holiday, error)
	GetByID(id int) (*model.Holiday, error)
	GetByYear(year int) ([]*model.Holiday, error)
	GetByYearAndMonth(year, month int) ([]*model.Holiday, error)
	GetByYMD(year, month, day int) ([]*model.Holiday, error)
	Insert(u *model.Holiday) (int, error)
	UpdateNameAndType(id int, name, holidayType string) (int64, error)
}

// RoleRepo is an interface that represents the list of functions that need to be
// implemented for the Role model, by the repo.
type RoleRepo interface {
	CheckIfNameExists(name string, id int) (*model.Role, error)
	GetByID(id int) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Insert(r *model.Role) (lastInsertID int, err error)
	UpdateName(roleID int, name string) (rowsAffected int64, err error)
}
