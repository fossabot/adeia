package user

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
	"adeia-api/internal/service"
	"errors"
	"fmt"
	"time"
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

func (i *Impl) CreateHoliday(holiday model.Holiday) error {
	existingHoliday, err := i.holidayRepo.GetByEpoch(time.Parse("",holiday.HolidayDate).Unix())
	if err != nil {
		return fmt.Errorf("cannot find existing holiday with the provided date: %v", err)
	}
	if existingHoliday != nil {
		return errors.New("holiday already exists with the provided date")
	}
	_, err = i.holidayRepo.Insert(&holiday)
	return err
}

func (i *Impl) GetHolidayByDate(date model.Date, timeUnit model.TimeUnit) (*[]model.Holiday, error) {
	var err = errors.New("Time Unit Not Found")
	var holiday *[]model.Holiday
	switch timeUnit {
	case model.Year:
		holiday, err = i.holidayRepo.GetByYear(date.Year)
		break
	case model.Month:
		holiday, err = i.holidayRepo.GetByYear(date.Year)
		break
	case model.Epoch:
		holiday, err = i.holidayRepo.GetByEpoch(date.Epoch)
		break
	}
	if err != nil {
		return nil, err
	}
	return holiday, nil
}
