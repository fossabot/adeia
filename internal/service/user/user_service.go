package user

import (
	"strings"

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
	CreateUser(name, email, empID, designation string) (string, error)
	GetUserByEmpID(id string) (*model.User, error)
}

// Impl is a Service implementation.
type Impl struct {
	usrRepo repo.UserRepo
	mailer  mail.Mailer
}

// New creates a new Service.
func New(d db.DB, c cache.Cache, m mail.Mailer) Service {
	u := repo.NewUserRepo(d)
	return &Impl{
		usrRepo: u,
		mailer:  m,
	}
}

// CreateUser creates a new user.
func (i *Impl) CreateUser(name, email, empID, designation string) (string, error) {
	// check if user already exists
	usr, err := i.usrRepo.GetByEmail(email)
	if err != nil {
		log.Errorf("cannot find if an user already exists with the provided email: %v", err)
		return "", util.ErrInternalServerError
	} else if usr != nil {
		log.Warnf("user already exists with the provided email %s", email)
		return "", util.ErrResourceAlreadyExists
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
		u.EmployeeID = model.NewEmpID()
	} else {
		u.EmployeeID = strings.ToUpper(empID)
	}

	// add email verification token
	randBytes, err := crypto.GenerateRandomBytes(128)
	if err != nil {
		log.Error("error generating random bytes: %v", err)
		return "", util.ErrInternalServerError
	}
	token := crypto.EncodeBase64(randBytes)
	u.EmailVerificationToken = crypto.Hash(randBytes)

	// create user
	if _, err = i.usrRepo.Insert(u); err != nil {
		log.Error("cannot create new user: %v", err)
		return "", util.ErrInternalServerError
	}

	// send verification email
	newEmail := mail.NewEmailBuilder().
		To([]string{u.Email}).
		Subject("Verify your email address - Adeia").
		Build()
	emailData := mail.TemplateEmailVerifyData
	// TODO: use URL from front-end's verification page
	emailData.Link = "http://frontend/verify?token=" + token
	if err := i.mailer.Send(newEmail, mail.TemplateEmailVerify, emailData); err != nil {
		log.Error("cannot send verification email: %v", err)
		return "", util.ErrInternalServerError
	}

	return u.EmployeeID, nil
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
