package user

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo"
	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/log"
	"database/sql"
	"time"
)

type Service interface {
	CreateHoliday(holiday model.Holiday) (*model.Holiday, error)
	GetHolidayByDate(date time.Time, timeUnit constants.TimeUnit) ([]*model.Holiday, error)
	GetHolidayById(id int) (*model.Holiday, error)
	UpdateById(holiday model.Holiday, id int) error
	DeleteById(id int) error
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

func (i *Impl) CreateHoliday(holiday model.Holiday) (*model.Holiday, error) {
	existingHoliday, err := i.holidayRepo.GetByYMD(util.GetYMDFromTime(holiday.HolidayDate))
	if err != nil {
		log.Errorf("Error while fetching holiday from Database : %v", err)
		return nil, util.ErrDatabaseError
	}
	if existingHoliday != nil {
		log.Errorf("Holiday already exists : %v", existingHoliday)
		return nil, util.ErrResourceAlreadyExists
	}
	holidayId, err := i.holidayRepo.Insert(&holiday)
	holiday.ID = holidayId
	return &holiday, err
}

func (i *Impl) GetHolidayByDate(date time.Time, granularity constants.TimeUnit) ([]*model.Holiday, error) {
	var err error
	var holiday []*model.Holiday
	switch granularity {
	case constants.Year:
		holiday, err = i.holidayRepo.GetByYear(date.Year())
		break
	case constants.Month:
		holiday, err = i.holidayRepo.GetByYearAndMonth(date.Year(), int(date.Month()))
		break
	case constants.DateOfMonth:
		holiday, err = i.holidayRepo.GetByYMD(util.GetYMDFromTime(date))
		break
	case constants.Epoch:
		holiday, err = i.holidayRepo.GetByEpoch(date.Unix())
		break
	}
	if err != nil {
		return nil, util.ErrDatabaseError.Msgf("Error : %v", err)
	}
	return holiday, nil
}

func (i *Impl) GetHolidayById(id int) (*model.Holiday, error) {
	if holiday, err := i.holidayRepo.GetByID(id); err != nil {
		log.Errorf("Database Error : %v", err)
		return nil, util.ErrDatabaseError
	} else {
		return holiday, nil
	}

}

func (i *Impl) UpdateById(holiday model.Holiday, id int) error {
	holiday.ID = id
	err := i.holidayRepo.UpdateNameAndType(holiday)
	if err == sql.ErrNoRows {
		return util.ErrResourceNotFound
	} else if err != nil {
		log.Errorf("Database Error : %v", err)
		return util.ErrDatabaseError
	} else {
		return nil
	}
}

func (i *Impl) DeleteById(id int) error {
	holiday := model.Holiday{ID: id}
	err := i.holidayRepo.DeletedById(holiday)
	if err == sql.ErrNoRows {
		return util.ErrResourceNotFound
	} else if err != nil {
		log.Errorf("Database Error : %v", err)
		return util.ErrDatabaseError
	} else {
		return nil
	}
}
