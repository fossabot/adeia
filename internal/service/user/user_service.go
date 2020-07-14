package user

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	userRepo "adeia-api/internal/repo/user"
	"adeia-api/internal/util"
	"adeia-api/internal/util/crypto"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/mail"
)

// Service contains all user-related business logic.
type Service interface {
	ActivateUser(empID, email, password string) (*model.User, error)
	CreateUser(name, email, empID, designation string) (*model.User, error)
	DeleteUser(empID string) error
	GetUserByEmpID(empID string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	LoginUser(email, password string) (*model.User, error)
}

// Impl is a Service implementation.
type Impl struct {
	cache    cache.Cache
	mailer   mail.Mailer
	userRepo userRepo.Repo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache, m mail.Mailer) Service {
	u := userRepo.New(d)
	return &Impl{
		cache:    c,
		mailer:   m,
		userRepo: u,
	}
}

// ActivateUser activates an user account.
func (i *Impl) ActivateUser(empID, email, password string) (*model.User, error) {
	// check if user exists
	usr, err := i.userRepo.GetByEmpID(empID)
	if err != nil {
		log.Errorf("cannot find user by empID and email: %v", err)
		return nil, util.ErrDatabaseError
	} else if usr == nil || usr.Email != email {
		log.Errorf("email and empID do not match: %v", err)
		return nil, util.ErrResourceNotFound
	}

	if usr.IsActivated {
		log.Error("user already activated")
		return nil, util.ErrBadRequest.Msg("Account already activated")
	}

	// user exists, hash and store password
	hash, err := crypto.HashPassword(password)
	if err != nil {
		log.Errorf("cannot generate hash for password: %v", err)
		return nil, util.ErrInternalServerError
	}
	if err := i.userRepo.UpdatePasswordAndIsActivated(usr, hash, true); err != nil {
		log.Errorf("cannot update user: %v", err)
		return nil, util.ErrDatabaseError
	}

	return usr, nil
}

// CreateUser creates a new user.
func (i *Impl) CreateUser(name, email, empID, designation string) (*model.User, error) {
	// check if user already exists
	usr, err := i.userRepo.GetByEmailInclDeleted(email)
	if err != nil {
		log.Errorf("cannot find if an user already exists with the provided email: %v", err)
		return nil, util.ErrDatabaseError
	} else if usr != nil {
		log.Warnf("user already exists with the provided email %s", email)
		return nil, util.ErrResourceAlreadyExists
	}

	// user does not exist, so create one
	u := &model.User{
		Name:        name,
		Email:       email,
		Password:    "",
		Designation: designation,
		IsActivated: false,
	}
	if empID == "" {
		u.EmployeeID = crypto.NewEmpID()
	} else {
		u.EmployeeID = empID
	}

	// create user
	if _, err = i.userRepo.Insert(u); err != nil {
		log.Error("cannot create new user: %v", err)
		return nil, util.ErrDatabaseError
	}
	return u, nil
}

// DeleteUser deletes a user.
func (i *Impl) DeleteUser(empID string) error {
	rowsAffected, err := i.userRepo.DeleteByEmpID(empID)
	if err != nil {
		log.Errorf("cannot delete user: %v", err)
		return util.ErrDatabaseError
	} else if rowsAffected == 0 {
		log.Errorf("no user found with empID: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}

// GetUserByEmpID gets a user using the provided empID.
func (i *Impl) GetUserByEmpID(empID string) (*model.User, error) {
	usr, err := i.userRepo.GetByEmpID(empID)
	if err != nil {
		log.Errorf("cannot find user with the provided employee ID: %v", err)
		return nil, util.ErrDatabaseError
	} else if usr == nil {
		log.Warnf("user does not exist for the provided employee ID %v", empID)
		return nil, util.ErrResourceNotFound
	}

	return usr, nil
}

// GetUserByID gets a user using the provided id.
func (i *Impl) GetUserByID(id int) (*model.User, error) {
	usr, err := i.userRepo.GetByID(id)
	if err != nil {
		log.Errorf("cannot find user with the provided ID: %v", err)
		return nil, util.ErrDatabaseError
	} else if usr == nil {
		log.Warnf("user does not exist for the provided ID %v", id)
		return nil, util.ErrResourceNotFound
	}

	return usr, nil
}

// LoginUser logs in a user.
func (i *Impl) LoginUser(email, password string) (*model.User, error) {
	// check if user exists
	usr, err := i.userRepo.GetByEmail(email)
	if err != nil {
		log.Errorf("cannot find user by email: %v", err)
		return nil, util.ErrDatabaseError
	} else if usr == nil {
		log.Errorf("user not found with the specified email: %v", err)
		return nil, util.ErrUnauthorized.Msg("Wrong credentials")
	}

	// check if account is activated
	if !usr.IsActivated {
		log.Error("account not activated")
		return nil, util.ErrAccountNotActivated
	}

	// user exists, check credentials
	match, err := crypto.ComparePwdHash(password, usr.Password)
	if err != nil || !match {
		log.Errorf("error comparing password and hash: %v", err)
		return nil, util.ErrUnauthorized.Msg("Wrong credentials")
	}

	return usr, nil
}
