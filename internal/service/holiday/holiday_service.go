package user

import (
	"adeia-api/internal/service"
	"errors"
	"fmt"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
)


// Impl is a Service implementation.
type Impl struct {
	holidayRepo repo.HolidayRepo
}

// New creates a new Service.
func New(d db.DB, c cache.Cache) service.HolidayService {
	holiday := repo.NewHolidayRepo(d)
	return &Impl{holiday}
}

// CreateUser creates a new user.
func (i *Impl) CreateHoliday(holiday model.Holiday) error {
	// check if user already exists
	usr, err := i.holidayRepo.GetByDate(holiday.HolidayDate)
	if err != nil {
		return fmt.Errorf("cannot find existing holiday with the provided date: %v", err)
	}
	if usr != nil {
		return errors.New("holiday already exists with the provided date")
	}
	_, err = i.holidayRepo.Insert(&holiday)
	return err
}
