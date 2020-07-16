package holiday

import (
	"database/sql"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
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

// Repo is an interface that represents the list of functions that need to be
// implemented for the Holiday model, by the repo.
type Repo interface {
	DeletedByID(id int) (int64, error)
	GetByEpoch(epoch int64) ([]*model.Holiday, error)
	GetByID(id int) (*model.Holiday, error)
	GetByYear(year int) ([]*model.Holiday, error)
	GetByYearAndMonth(year, month int) ([]*model.Holiday, error)
	GetByYMD(year, month, day int) ([]*model.Holiday, error)
	Insert(u *model.Holiday) (int, error)
	UpdateNameAndType(id int, name, holidayType string) (int64, error)
}

// Impl is an implementation of Repo for Postgres.
type Impl struct {
	db db.DB
}

// New creates a new Repo.
func New(d db.DB) Repo {
	return &Impl{d}
}

// Insert inserts a holiday using the db connection instance and returns the LastInsertID.
func (i *Impl) Insert(h *model.Holiday) (int, error) {
	stmt, err := i.db.PrepareNamed(queryInsert)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastInsertID int
	if err := stmt.Get(&lastInsertID, h); err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// GetByID gets a holiday from db using the id.
func (i *Impl) GetByID(id int) (*model.Holiday, error) {
	var u model.Holiday
	if err := i.db.Get(&u, queryByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByEpoch gets holidays by epoch.
func (i *Impl) GetByEpoch(epoch int64) ([]*model.Holiday, error) {
	return i.get(queryByEpoch, epoch)
}

// GetByYear gets holidays by year.
func (i *Impl) GetByYear(year int) ([]*model.Holiday, error) {
	return i.get(queryByYear, year)
}

// GetByYearAndMonth gets holidays by year and month.
func (i *Impl) GetByYearAndMonth(year, month int) ([]*model.Holiday, error) {
	return i.get(queryByYearAndMonth, year, month)
}

// GetByYMD gets holidays by year, month and day.
func (i *Impl) GetByYMD(year, month, day int) ([]*model.Holiday, error) {
	return i.get(queryByYearAndMonthAndDay, year, month, day)
}

func (i *Impl) get(query string, args ...interface{}) ([]*model.Holiday, error) {
	var u []*model.Holiday
	if err := i.db.Select(&u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// UpdateNameAndType updates a holiday's name and type.
func (i *Impl) UpdateNameAndType(id int, name, holidayType string) (int64, error) {
	result, err := i.db.Exec(queryUpdateNameAndType, name, holidayType, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// DeletedByID deletes a holiday by id.
func (i *Impl) DeletedByID(id int) (int64, error) {
	result, err := i.db.Exec(queryDeleteByID, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
