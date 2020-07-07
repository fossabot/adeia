package user

import (
	"time"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/log"
)

// Service contains all holiday-related business logic.
type Service interface {
	CreateHoliday(holiday model.Holiday) (*model.Holiday, error)
	GetHolidaysByDate(date time.Time, timeUnit constants.TimeUnit) ([]*model.Holiday, error)
	GetHolidayByID(id int) (*model.Holiday, error)
	UpdateByID(id int, holiday *model.Holiday) error
	DeleteByID(id int) error
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

// CreateHoliday creates a holiday.
func (i *Impl) CreateHoliday(holiday model.Holiday) (*model.Holiday, error) {
	existingHoliday, err := i.holidayRepo.GetByYMD(util.GetYMDFromTime(holiday.HolidayDate))
	if err != nil {
		log.Errorf("Error while fetching holiday from Database : %v", err)
		return nil, util.ErrDatabaseError
	} else if existingHoliday != nil {
		log.Errorf("Holiday already exists : %v", existingHoliday)
		return nil, util.ErrResourceAlreadyExists
	}

	id, err := i.holidayRepo.Insert(&holiday)
	holiday.ID = id
	return &holiday, err
}

// GetHolidaysByDate gets all holidays using the provided date.
func (i *Impl) GetHolidaysByDate(date time.Time, granularity constants.TimeUnit) ([]*model.Holiday, error) {
	var (
		err      error
		holidays []*model.Holiday
	)

	switch granularity {
	case constants.Year:
		holidays, err = i.holidayRepo.GetByYear(date.Year())
	case constants.Month:
		holidays, err = i.holidayRepo.GetByYearAndMonth(date.Year(), int(date.Month()))
	case constants.DateOfMonth:
		holidays, err = i.holidayRepo.GetByYMD(util.GetYMDFromTime(date))
	case constants.Epoch:
		holidays, err = i.holidayRepo.GetByEpoch(date.Unix())
	default:
		log.Error("specified granularity cannot be used")
		return nil, util.ErrInternalServerError
	}

	if err != nil {
		log.Errorf("cannot fetch holiday from db: %v", err)
		return nil, util.ErrDatabaseError
	}
	return holidays, nil
}

// GetHolidayByID gets a holiday by id.
func (i *Impl) GetHolidayByID(id int) (*model.Holiday, error) {
	holiday, err := i.holidayRepo.GetByID(id)
	if err != nil {
		log.Errorf("Database Error : %v", err)
		return nil, util.ErrDatabaseError
	} else if holiday == nil {
		log.Error("no holiday found for id")
		return nil, util.ErrResourceNotFound
	}

	return holiday, nil
}

// UpdateByID updates a holiday by id.
func (i *Impl) UpdateByID(id int, holiday *model.Holiday) error {
	rowsAffected, err := i.holidayRepo.UpdateNameAndType(id, holiday.Name, holiday.HolidayType)
	if err != nil {
		log.Errorf("Database Error: %v", err)
		return util.ErrDatabaseError
	} else if rowsAffected == 0 {
		log.Errorf("no holiday found with provided id: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}

// DeleteByID deletes a holiday by id.
func (i *Impl) DeleteByID(id int) error {
	rowsAffected, err := i.holidayRepo.DeletedByID(id)
	if err != nil {
		log.Errorf("Database Error : %v", err)
		return util.ErrDatabaseError
	} else if rowsAffected == 0 {
		log.Errorf("no holiday found with id: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}
