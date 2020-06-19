package user

import (
	"errors"
	"fmt"

	"adeia-api/internal/model"
	"adeia-api/internal/repository"
	"adeia-api/internal/service/cache"
	"adeia-api/internal/service/db"
)

type Service interface {
	CreateUser(empID string) (*model.User, error)
}

type Impl struct {
	usrRepo repository.UserRepo
}

func New(d db.DB, c cache.Cache) Service {
	u := repository.NewUserRepo(d)
	return &Impl{u}
}

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
