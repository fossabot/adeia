package user

import (
	"errors"
	"fmt"
	"time"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
)

type Service interface {
	CreateHoliday(holiday model.Holiday) error
	GetHolidayByYear(year int) (*[]model.Holiday, error)
	GetHolidayByYearAndMonth(year, month int) (*[]model.Holiday, error)
}

// Impl is a Service implementation.
type Impl struct {
	holidayRepo repo.HolidayRepo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache) Service {
	holiday := repo.NewHolidayRepo(d)
	return &Impl{holiday}
}

func (i *Impl) CreateHoliday(holiday model.Holiday) error {
	existingHoliday, err := i.holidayRepo.GetByDate(holiday.HolidayDate)
	if err != nil {
		return fmt.Errorf("cannot find existing holiday with the provided date: %v", err)
	}
	if existingHoliday != nil {
		return errors.New("holiday already exists with the provided date")
	}
	_, err = i.holidayRepo.Insert(&holiday)
	return err
}

func (i *Impl) GetHolidayByDate(date time.Time) (*model.Holiday, error) {
	holiday, err := i.holidayRepo.GetByDate(date)
	if err != nil {
		return nil, err
	}
	return holiday, nil
}

func (i *Impl) GetHolidayByYear(year int) (*[]model.Holiday, error) {
	holiday, err := i.holidayRepo.GetByYear(year)
	if err != nil {
		return nil, err
	}
	return holiday, nil
}

func (i *Impl) GetHolidayByYearAndMonth(year, month int) (*[]model.Holiday, error) {
	holiday, err := i.holidayRepo.GetByYearAndMonth(year, month)
	if err != nil {
		return nil, err
	}
	return holiday, nil
}
