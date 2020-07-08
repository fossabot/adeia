package repo

import (
	"database/sql"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
)

const (
	queryHolidayInsert                = "INSERT INTO holidays (date, name, type) VALUES (:date, :name, :type) RETURNING id"
	queryHolidayByID                  = "SELECT * FROM holidays WHERE id=$1"
	queryHolidayByDate                = "SELECT * FROM holidays WHERE EXTRACT(EPOCH FROM date)=$1"
	queryHolidayByYear                = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1"
	queryHolidayByYearAndMonth        = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1 AND EXTRACT(MONTH FROM date)=$2"
	queryHolidayByYearAndMonthAndDate = "SELECT * FROM holidays WHERE EXTRACT(YEAR FROM date)=$1 AND EXTRACT(MONTH FROM date)=$2 AND EXTRACT(DAY FROM date)=$3"
	queryUpdateNameAndType            = "UPDATE holidays SET name=$1, type=$2 WHERE id=$3"
	queryDeleteByID                   = "DELETE FROM holidays WHERE id=$1"
)

// HolidayRepoImpl is an implementation of HolidayRepo for Postgres.
type HolidayRepoImpl struct {
	db db.DB
}

// NewHolidayRepo creates a new HolidayRepo.
func NewHolidayRepo(d db.DB) HolidayRepo {
	return &HolidayRepoImpl{d}
}

// Insert inserts a holiday using the db connection instance and returns the LastInsertID.
func (i *HolidayRepoImpl) Insert(h *model.Holiday) (int, error) {
	stmt, err := i.db.PrepareNamed(queryHolidayInsert)
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
func (i *HolidayRepoImpl) GetByID(id int) (*model.Holiday, error) {
	var u model.Holiday
	if err := i.db.Get(&u, queryHolidayByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// GetByEpoch gets holidays by epoch.
func (i *HolidayRepoImpl) GetByEpoch(epoch int64) ([]*model.Holiday, error) {
	return i.get(queryHolidayByDate, epoch)
}

// GetByYear gets holidays by year.
func (i *HolidayRepoImpl) GetByYear(year int) ([]*model.Holiday, error) {
	return i.get(queryHolidayByYear, year)
}

// GetByYearAndMonth gets holidays by year and month.
func (i *HolidayRepoImpl) GetByYearAndMonth(year, month int) ([]*model.Holiday, error) {
	return i.get(queryHolidayByYearAndMonth, year, month)
}

// GetByYMD gets holidays by year, month and date.
func (i *HolidayRepoImpl) GetByYMD(year, month, dateOfMonth int) ([]*model.Holiday, error) {
	return i.get(queryHolidayByYearAndMonthAndDate, year, month, dateOfMonth)
}

func (i *HolidayRepoImpl) get(query string, args ...interface{}) ([]*model.Holiday, error) {
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
func (i *HolidayRepoImpl) UpdateNameAndType(id int, name, holidayType string) (int64, error) {
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
func (i *HolidayRepoImpl) DeletedByID(id int) (int64, error) {
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
