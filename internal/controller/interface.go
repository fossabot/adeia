package controller

import (
	"net/http"
	"time"

	"adeia-api/internal/model"
	"adeia-api/internal/util/constants"
)

type UserService interface {
	ActivateUser(empID, email, password string) (*model.User, error)
	CreateUser(name, email, empID, designation string) (*model.User, error)
	DeleteUser(empID string) error
	GetUserByEmpID(empID string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	LoginUser(email, password string) (*model.User, error)
}

// SessionService contains all session-related business logic.
type SessionService interface {
	NewSession(id int, empID string) (accessToken, refreshToken string, err error)
	RefreshToken(id int, empID, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	AddRefreshTokenCookie(w http.ResponseWriter, refreshToken string)
	ReadRefreshTokenCookie(r *http.Request) (string, error)
	ParseAccessToken(jwt string) (id string, err error)
	Destroy(id int, refreshToken string) error
}

// RoleService contains all role-related business logic.
type RoleService interface {
	CreateRole(name string) (*model.Role, error)
	UpdateByID(roleID int, role *model.Role) error
}

// HolidayService contains all holiday-related business logic.
type HolidayService interface {
	CreateHoliday(holiday model.Holiday) (*model.Holiday, error)
	GetHolidaysByDate(date time.Time, timeUnit constants.TimeUnit) ([]*model.Holiday, error)
	GetHolidayByID(id int) (*model.Holiday, error)
	UpdateByID(id int, holiday *model.Holiday) error
	DeleteByID(id int) error
}
