package repo

import (
	"database/sql"
	"time"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
)

const (
	queryHolidayInsert = "INSERT INTO holiday (date, name, type) VALUES (:date, :name, :type) RETURNING id"
	queryHolidayById    = "SELECT * FROM holiday WHERE id=$1"
	queryHolidayByDate = "SELECT * FROM holiday WHERE date=$1"
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
	u := model.Holiday{}
	if err := i.db.Get(&u, queryHolidayById, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (i *HolidayRepoImpl) GetByDate(date time.Time) (*model.Holiday, error) {
	u := model.Holiday{}
	if err := i.db.Get(&u, queryHolidayByDate, date); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

