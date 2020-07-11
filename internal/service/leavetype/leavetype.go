package leavetype

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
)

type Service interface {
	CreateLeaveType(holiday model.LeaveType) (*model.LeaveType, error)
	GetLeaveTypeByID(id int) (*model.LeaveType, error)
	UpdateByID(id int, holiday *model.LeaveType) error
	DeleteByID(id int) error
}

type LeaveTypeServiceImpl struct{
	repo repo.LeaveTypeRepo
}

func New(d db.DB, c cache.Cache) Service {
	leaveTypeRepo := repo.NewLeaveTypeRepo(d)
	return &LeaveTypeServiceImpl{leaveTypeRepo}
}
func (i *LeaveTypeServiceImpl) CreateLeaveType(leaveType model.LeaveType) (*model.LeaveType, error) {
	id, err := i.repo.Insert(&leaveType)
	if(err!=nil){
		return nil, err
	}
	leaveType.ID = int(id)
	return &leaveType,nil
}

func (i *LeaveTypeServiceImpl) GetLeaveTypeByID(id int) (*model.LeaveType, error) {
	return i.repo.GetByID(id)
}

func (i *LeaveTypeServiceImpl) UpdateByID(id int, leaveType *model.LeaveType) error {
	_ , err := i.repo.UpdateNameAndDays(id, leaveType.Name, leaveType.Days)
	if err!=nil {
		return err
	}
	return nil
}

func (i *LeaveTypeServiceImpl) DeleteByID(id int) error {
	_ , err := i.repo.DeletedByID(id)
	if err!=nil {
		return err
	}
	return nil
}
