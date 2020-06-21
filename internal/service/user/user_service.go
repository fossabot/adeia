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
	CreateUser(name, email, empID, designation string) error
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
func (i *Impl) CreateUser(name, email, empID, designation string) error {
	// check if user already exists
	usr, err := i.usrRepo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("cannot find existing user with the provided email: %v", err)
	}
	if usr != nil {
		return errors.New("user already exists with the provided email")
	}

	// user does not exist, so create one
	u := &model.User{
		EmployeeID:  empID,
		Name:        name,
		Email:       email,
		Password:    "dummy",
		Designation: designation,
		IsActivated: false,
	}
	_, err = i.usrRepo.Insert(u)
	return err
}
