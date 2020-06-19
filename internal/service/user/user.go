package user

import (
	"errors"
	"fmt"

	"adeia-api/internal/model"
	"adeia-api/internal/repository"
)

var usrRepo = &repository.UserRepoImpl{}

type Service interface {
	CreateUser(empID string) (*model.User, error)
}

type Impl struct{}

func (i *Impl) CreateUser(empID string) (*model.User, error) {
	// check if user already exists
	usr, err := usrRepo.GetByEmpID(empID)
	if err != nil {
		return nil, fmt.Errorf("cannot find existing user with provided empID: %v", err)
	}

	if usr != nil {
		return nil, errors.New("user already exists with provided empID")
	}

	// user does not exist, so create one
	u := &model.User{
		EmployeeID: "",
		Name:       "",
		Email:      "",
		Password:   "",
	}
	err = usrRepo.Insert(u)
	if err != nil {
		return nil, fmt.Errorf("cannot create new user: %v", err)
	}

	// successfully inserted user
	return u, nil
}
