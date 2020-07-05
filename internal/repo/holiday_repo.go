package repo

import (
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"database/sql"
)

const (
	queryHolidayInsert                = "INSERT INTO holiday (date, name, type) VALUES (:date, :name, :type) RETURNING id"
	queryHolidayById                  = "SELECT * FROM holiday WHERE id=$1"
	queryHolidayByDate                = "SELECT * FROM holiday WHERE EXTRACT(EPOCH FROM date)=$1"
	queryHolidayByYear                = "SELECT * FROM holiday WHERE EXTRACT(YEAR FROM date)=$1"
	queryHolidayByYearAndMonth        = "SELECT * FROM holiday WHERE EXTRACT(YEAR FROM date)=$1 and EXTRACT(MONTH from date)=$2"
	queryHolidayByYearAndMonthAndDate = "SELECT * FROM holiday WHERE EXTRACT(YEAR FROM date)=$1 and EXTRACT(MONTH from date)=$2 and EXTRACT(DAY from date)=$3"
)

// HolidayRepoImpl is an implementation of HolidayRepo for Postgres.
type HolidayRepoImpl struct {
	db db.DB
}

func NewHolidayRepo(d db.DB) HolidayRepo {
	return &HolidayRepoImpl{d}
}

func (i *HolidayRepoImpl) Insert(u *model.Holiday) (int, error) {
	stmt, err := i.db.PrepareNamed(queryHolidayInsert)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastInsertID int
	if err := stmt.Get(&lastInsertID, u); err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (i *HolidayRepoImpl) GetByID(id int) (*model.Holiday, error) {
	var u model.Holiday
	if err := i.db.Get(&u, queryHolidayById, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (i *HolidayRepoImpl) GetByEpoch(epoch int64) ([]*model.Holiday, error) {
	return i.get(queryHolidayByDate, epoch)
}

func (i *HolidayRepoImpl) GetByYear(year int) ([]*model.Holiday, error) {
	return i.get(queryHolidayByYear, year)
}

func (i *HolidayRepoImpl) GetByYearAndMonth(year, month int) ([]*model.Holiday, error) {
	return i.get(queryHolidayByYearAndMonth, year, month)
}

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