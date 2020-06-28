package user

import (
	"strconv"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
	"adeia-api/internal/util"
	"adeia-api/internal/util/crypto"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/mail"
)

// Service contains all user-related business logic.
type Service interface {
	CreateUser(name, email, empID, designation string) (*model.User, error)
	GetUserByEmpID(empID string) (*model.User, error)
	ActivateUser(empID, email, password string) (*model.User, error)
	DeleteUser(empID string) error
	LoginUser(email, password string) (sessID string, err error)
	GetAuthenticated(sessID string) (*model.User, error)
}

// Impl is a Service implementation.
type Impl struct {
	cache   cache.Cache
	mailer  mail.Mailer
	usrRepo repo.UserRepo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache, m mail.Mailer) Service {
	u := repo.NewUserRepo(d)
	return &Impl{
		cache:   c,
		mailer:  m,
		usrRepo: u,
	}
}

// CreateUser creates a new user.
func (i *Impl) CreateUser(name, email, empID, designation string) (*model.User, error) {
	// check if user already exists
	usr, err := i.usrRepo.GetByEmail(email)
	if err != nil {
		log.Errorf("cannot find if an user already exists with the provided email: %v", err)
		return nil, util.ErrInternalServerError
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
	if _, err = i.usrRepo.Insert(u); err != nil {
		log.Error("cannot create new user: %v", err)
		return nil, util.ErrInternalServerError
	}
	return u, nil
}

// GetUserByEmpID gets a user using the provided empID.
func (i *Impl) GetUserByEmpID(empID string) (*model.User, error) {
	usr, err := i.usrRepo.GetByEmpID(empID)
	if err != nil {
		log.Errorf("cannot find user with the provided employee ID: %v", err)
		return nil, util.ErrInternalServerError
	} else if usr == nil {
		log.Warnf("user does not exist for the provided employee ID %v", empID)
		return nil, util.ErrResourceNotFound
	}

	return usr, nil
}

// ActivateUser activates a user account.
func (i *Impl) ActivateUser(empID, email, password string) (*model.User, error) {
	// check if user exists
	usr, err := i.usrRepo.GetByEmpIDAndEmail(empID, email)
	if err != nil {
		log.Errorf("cannot find user by empID and email: %v", err)
		return nil, util.ErrInternalServerError
	} else if usr == nil {
		log.Errorf("user not found with the specified empID and email: %v", err)
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
	if err := i.usrRepo.UpdatePasswordAndIsActivated(usr, hash, true); err != nil {
		log.Errorf("cannot update user: %v", err)
		return nil, util.ErrInternalServerError
	}

	return usr, nil
}

// DeleteUser deletes a user.
func (i *Impl) DeleteUser(empID string) error {
	rowsAffected, err := i.usrRepo.DeleteByEmpID(empID)
	if err != nil {
		log.Errorf("cannot delete user: %v", err)
		return util.ErrInternalServerError
	} else if rowsAffected == 0 {
		log.Errorf("no user found with empID: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}

// LoginUser logs in a user.
func (i *Impl) LoginUser(email, password string) (sessID string, err error) {
	// check if user exists
	usr, err := i.usrRepo.GetByEmail(email)
	if err != nil {
		log.Errorf("cannot find user by email: %v", err)
		return "", util.ErrInternalServerError
	} else if usr == nil {
		log.Errorf("user not found with the specified email: %v", err)
		return "", util.ErrUnauthorized.Msg("Wrong credentials")
	}

	// check if account is activated
	if !usr.IsActivated {
		log.Error("account not activated")
		return "", util.ErrAccountNotActivated
	}

	// user exists, check credentials
	match, err := crypto.ComparePwdHash(password, usr.Password)
	if err != nil || !match {
		log.Errorf("error comparing password and hash: %v", err)
		return "", util.ErrUnauthorized.Msg("Wrong credentials")
	}

	// credentials match, create session
	sessID, err = crypto.NewUUID()
	if err != nil {
		log.Errorf("cannot generate new UUID: %v", err)
		return "", util.ErrInternalServerError
	}

	// add session to cache
	if err := i.cache.StoreSession(sessID, usr.ID); err != nil {
		log.Errorf("cannot cache session: %v", err)
		return "", util.ErrInternalServerError
	}

	return sessID, nil
}

// GetAuthenticated gets the user identified by their sessID.
func (i *Impl) GetAuthenticated(sessID string) (*model.User, error) {
	// get user id for associated sessID
	var idStr string
	if err := i.cache.GetSession(&idStr, sessID); err != nil {
		log.Errorf("cannot get user id from cache for sessID: %v", err)
		return nil, util.ErrUnauthorized
	}

	// get user
	id, _ := strconv.Atoi(idStr)
	usr, err := i.usrRepo.GetByID(id)
	if err != nil {
		log.Errorf("cannot get user for id: %v", err)
		return nil, util.ErrUnauthorized
	} else if usr == nil {
		log.Errorf("user does not exist for the provided id %v", id)
		return nil, util.ErrUnauthorized
	}

	return usr, nil
}
