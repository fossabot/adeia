package repo

import (
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"database/sql"
)

type LeaveTypeRepoImpl struct {
	db db.DB
}

const (
	queryLeaveTypeInsert     = "INSERT INTO leave_types (name, days) VALUES (:name, :days) RETURNING id"
	queryLeaveTypeByID       = "SELECT * FROM leave_types WHERE id=$1"
	queryUpdateLeaveType     = "UPDATE leave_types SET name=$1, days=$2 WHERE id=$3"
	queryDeleteLeaveTypeByID = "DELETE FROM leave_types WHERE id=$1"
)

func NewLeaveTypeRepo(d db.DB) LeaveTypeRepo {
	return &LeaveTypeRepoImpl{d}
}

func (i *LeaveTypeRepoImpl) GetByID(id int) (*model.LeaveType, error) {
	result, err := i.get(queryLeaveTypeByID, id)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (i *LeaveTypeRepoImpl) UpdateNameAndDays(id int, name string, days int) (int64, error) {
	return i.execute(queryUpdateLeaveType, model.LeaveType{ID: id, Name: name, Days: days}, false)
}

func (i *LeaveTypeRepoImpl) DeletedByID(id int) (int64, error) {
	return i.execute(queryDeleteLeaveTypeByID, model.LeaveType{ID: id}, false)
}

func (i *LeaveTypeRepoImpl) Insert(u *model.LeaveType) (int64, error) {
	return i.execute(queryLeaveTypeInsert, model.LeaveType{Name: u.Name, Days: u.Days}, true)
}

func (i *LeaveTypeRepoImpl) get(query string, args ...interface{}) ([]*model.LeaveType, error) {
	var u []*model.LeaveType
	if err := i.db.Select(&u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (i *LeaveTypeRepoImpl) execute(query string, leaveType model.LeaveType, isInsert bool) (int64, error) {
	var err error
	result, err := i.db.NamedExec(query, leaveType)
	if err == nil {
		if isInsert {
			return result.LastInsertId()
		}
		return result.RowsAffected()
	}
	return 0, err
}
