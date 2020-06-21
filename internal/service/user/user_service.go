package user

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
	"adeia-api/internal/util"
	log "adeia-api/internal/util/logger"
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
		log.Errorf("cannot find if an user already exists with the provided email: %v", err)
		return util.ErrInternalServerError
	} else if usr != nil {
		log.Warnf("user already exists with the provided email %s", email)
		return util.ErrResourceAlreadyExists
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
	if _, err = i.usrRepo.Insert(u); err != nil {
		log.Error("cannot create new user: %v", err)
		return util.ErrInternalServerError
	}

	return nil
}
