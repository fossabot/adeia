package role

import (
	"adeia-api/internal/model"
	"adeia-api/internal/service"
	"adeia-api/internal/util"
	"adeia-api/internal/util/log"
)

type Service struct {
	roleRepo service.RoleRepo
}

// New creates a new Service.
func New(r service.RoleRepo) *Service {
	return &Service{r}
}

// CreateRole creates a role.
func (s *Service) CreateRole(name string) (*model.Role, error) {
	// check for existing role
	if existingRole, err := s.roleRepo.GetByName(name); err != nil {
		log.Errorf("cannot find if role already exists with the provided name %s: %v", name, err)
		return nil, util.ErrDatabaseError
	} else if existingRole != nil {
		log.Errorf("role already exists with provided name %s", name)
		return nil, util.ErrResourceAlreadyExists
	}

	// create new role
	r := model.Role{Name: name}
	id, err := s.roleRepo.Insert(&r)
	if err != nil {
		log.Errorf("cannot create new role: %v", err)
		return nil, util.ErrDatabaseError
	}

	// return created role
	r.ID = id
	return &r, nil
}

// UpdateByID updates a role identified by the roleID.
func (s *Service) UpdateByID(roleID int, role *model.Role) error {
	// check for existing role, because name of role is unique
	if existingRole, err := s.roleRepo.CheckIfNameExists(role.Name, roleID); err != nil {
		log.Errorf("cannot find if role already exists with the provided name %s: %v", role.Name, err)
		return util.ErrDatabaseError
	} else if existingRole != nil {
		log.Errorf("role already exists with provided name %s", role.Name)
		return util.ErrResourceAlreadyExists
	}

	rowsAffected, err := s.roleRepo.UpdateName(roleID, role.Name)
	if err != nil {
		log.Errorf("Database Error: %v", err)
		return util.ErrDatabaseError
	} else if rowsAffected == 0 {
		log.Errorf("no holiday found with provided id: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}
