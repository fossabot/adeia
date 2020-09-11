package user

import (
	"adeia/internal/model"
	"adeia/internal/service"
	"adeia/internal/util"
	"adeia/internal/util/crypto"
	"adeia/internal/util/log"
)

type Service struct {
	cache    service.Cache
	mailer   service.Mailer
	userRepo service.UserRepo
}

// New creates a new Service.
func New(u service.UserRepo, c service.Cache, m service.Mailer) *Service {
	return &Service{
		cache:    c,
		mailer:   m,
		userRepo: u,
	}
}

// ActivateUser activates an user account.
func (s *Service) ActivateUser(empID, email, password string) (*model.User, error) {
	// check if user exists
	usr, err := s.userRepo.GetByEmpID(empID)
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
	if err := s.userRepo.UpdatePasswordAndIsActivated(usr, hash, true); err != nil {
		log.Errorf("cannot update user: %v", err)
		return nil, util.ErrDatabaseError
	}

	return usr, nil
}

// CreateUser creates a new user.
func (s *Service) CreateUser(name, email, empID, designation string) (*model.User, error) {
	// check if user already exists
	usr, err := s.userRepo.GetByEmailInclDeleted(email)
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
	if _, err = s.userRepo.Insert(u); err != nil {
		log.Error("cannot create new user: %v", err)
		return nil, util.ErrDatabaseError
	}
	return u, nil
}

// DeleteUser deletes a user.
func (s *Service) DeleteUser(empID string) error {
	rowsAffected, err := s.userRepo.DeleteByEmpID(empID)
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
func (s *Service) GetUserByEmpID(empID string) (*model.User, error) {
	usr, err := s.userRepo.GetByEmpID(empID)
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
func (s *Service) GetUserByID(id int) (*model.User, error) {
	usr, err := s.userRepo.GetByID(id)
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
func (s *Service) LoginUser(email, password string) (*model.User, error) {
	// check if user exists
	usr, err := s.userRepo.GetByEmail(email)
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
