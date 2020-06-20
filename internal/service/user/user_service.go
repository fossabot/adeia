package user

import (
	"errors"
	"fmt"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
)

// Service contains all user-related business logic.
type Service interface {
	CreateUser(empID string) (*model.User, error)
}

// Impl is a Service implementation.
type Impl struct {
	usrRepo repo.UserRepo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache) Service {
	u := repo.NewUserRepo(d)
	return &Impl{u}
}

// CreateUser creates a new user.
func (i *Impl) CreateUser(empID string) (*model.User, error) {
	// check if user already exists
	usr, err := i.usrRepo.GetByEmpID(empID)
	if err != nil {
		return nil, fmt.Errorf("cannot find existing user with provided empID: %v", err)
	}

	if usr != nil {
		return nil, errors.New("user already exists with provided empID")
	}

	// user does not exist, so create one
	u := &model.User{
		EmployeeID: empID,
		Name:       "",
		Email:      "",
		Password:   "",
	}
	err = i.usrRepo.Insert(u)
	if err != nil {
		return nil, fmt.Errorf("cannot create new user: %v", err)
	}

	// successfully inserted user
	return u, nil
}
