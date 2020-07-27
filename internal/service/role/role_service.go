package role

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	roleRepo "adeia-api/internal/repo/role"
	"adeia-api/internal/util"
	"adeia-api/internal/util/log"
)

// Service contains all role-related business logic.
type Service interface {
	CreateRole(name string) (*model.Role, error)
}

// Impl is a Service implementation.
type Impl struct {
	roleRepo roleRepo.Repo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache) Service {
	return &Impl{roleRepo.New(d)}
}

// CreateRole creates a role.
func (i *Impl) CreateRole(name string) (*model.Role, error) {
	// check for existing role
	if existingRole, err := i.roleRepo.GetByName(name); err != nil {
		log.Errorf("cannot find if role already exists with the provided name %s: %v", name, err)
		return nil, util.ErrDatabaseError
	} else if existingRole != nil {
		log.Errorf("role already exists with provided name %s", name)
		return nil, util.ErrResourceAlreadyExists
	}

	// create new role
	r := model.Role{Name: name}
	id, err := i.roleRepo.Insert(&r)
	if err != nil {
		log.Errorf("cannot create new role: %v", err)
		return nil, util.ErrDatabaseError
	}

	// return created role
	r.ID = id
	return &r, nil
}
