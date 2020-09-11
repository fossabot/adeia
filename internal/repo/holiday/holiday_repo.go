package holiday

import (
	"adeia/internal/model"
	"adeia/internal/repo"
)

const (
	queryInsert               = "INSERT INTO holidays (date, name, type) VALUES (:date, :name, :type) RETURNING id"
	queryByID                 = "SELECT * FROM holidays WHERE id=$1"
	queryByEpoch              = "SELECT * FROM holidays WHERE EXTRACT(EPOCH FROM date)=$1"
	queryByYear               = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1"
	queryByYearAndMonth       = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1 AND EXTRACT(MONTH FROM date)=$2"
	queryByYearAndMonthAndDay = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1 AND EXTRACT(MONTH FROM date)=$2 AND EXTRACT(DAY FROM date)=$3"
	queryUpdateNameAndType    = "UPDATE holidays SET name=$1, type=$2 WHERE id=$3"
	queryDeleteByID           = "DELETE FROM holidays WHERE id=$1"
)

type Repo struct {
	db repo.DB
}

// New creates a new Repo.
func New(d repo.DB) *Repo {
	return &Repo{d}
}

// Insert inserts a holiday using the db connection instance and returns the LastInsertID.
func (r *Repo) Insert(h *model.Holiday) (int, error) {
	return r.db.InsertNamed(queryInsert, h)
}

// GetByID gets a holiday from db using the id.
func (r *Repo) GetByID(id int) (*model.Holiday, error) {
	var u model.Holiday
	if ok, err := r.db.GetOne(&u, queryByID, id); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &u, nil
}

// GetByEpoch gets holidays by epoch.
func (r *Repo) GetByEpoch(epoch int64) ([]*model.Holiday, error) {
	return r.getMany(queryByEpoch, epoch)
}

// GetByYear gets holidays by year.
func (r *Repo) GetByYear(year int) ([]*model.Holiday, error) {
	return r.getMany(queryByYear, year)
}

// GetByYearAndMonth gets holidays by year and month.
func (r *Repo) GetByYearAndMonth(year, month int) ([]*model.Holiday, error) {
	return r.getMany(queryByYearAndMonth, year, month)
}

// GetByYMD gets holidays by year, month and day.
func (r *Repo) GetByYMD(year, month, day int) ([]*model.Holiday, error) {
	return r.getMany(queryByYearAndMonthAndDay, year, month, day)
}

// UpdateNameAndType updates a holiday's name and type.
func (r *Repo) UpdateNameAndType(id int, name, holidayType string) (int64, error) {
	return r.db.Update(queryUpdateNameAndType, name, holidayType, id)
}

// DeletedByID deletes a holiday by id.
func (r *Repo) DeletedByID(id int) (int64, error) {
	return r.db.Delete(queryDeleteByID, id)
}

// get is a generic getter that other `Get*` methods wrap on.
func (r *Repo) get(query string, args ...interface{}) (*model.Holiday, error) {
	h := model.Holiday{}
	if ok, err := r.db.GetOne(&h, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return &h, nil
}

func (r *Repo) getMany(query string, args ...interface{}) ([]*model.Holiday, error) {
	var u []*model.Holiday
	if ok, err := r.db.GetMany(&u, query, args...); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	return u, nil
}
