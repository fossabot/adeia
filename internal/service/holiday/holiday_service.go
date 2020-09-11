package user

import (
	"time"

	"adeia/internal/model"
	"adeia/internal/service"
	"adeia/internal/util"
	"adeia/internal/util/constants"
	"adeia/internal/util/log"
)

type Service struct {
	holidayRepo service.HolidayRepo
}

// New creates a new Service.
func New(h service.HolidayRepo) *Service {
	return &Service{h}
}

// CreateHoliday creates a holiday.
func (s *Service) CreateHoliday(holiday model.Holiday) (*model.Holiday, error) {
	existingHoliday, err := s.holidayRepo.GetByYMD(util.GetYMDFromTime(holiday.HolidayDate))
	if err != nil {
		log.Errorf("Error while fetching holiday from Database : %v", err)
		return nil, util.ErrDatabaseError
	} else if existingHoliday != nil {
		log.Errorf("Holiday already exists : %v", existingHoliday)
		return nil, util.ErrResourceAlreadyExists
	}

	id, err := s.holidayRepo.Insert(&holiday)
	if err != nil {
		log.Errorf("cannot create new holiday: %v", err)
		return nil, util.ErrDatabaseError
	}
	holiday.ID = id
	return &holiday, nil
}

// GetHolidaysByDate gets all holidays using the provided date.
func (s *Service) GetHolidaysByDate(date time.Time, granularity constants.TimeUnit) ([]*model.Holiday, error) {
	var (
		err      error
		holidays []*model.Holiday
	)

	switch granularity {
	case constants.Year:
		holidays, err = s.holidayRepo.GetByYear(date.Year())
	case constants.Month:
		holidays, err = s.holidayRepo.GetByYearAndMonth(date.Year(), int(date.Month()))
	case constants.DayOfMonth:
		holidays, err = s.holidayRepo.GetByYMD(util.GetYMDFromTime(date))
	case constants.Epoch:
		holidays, err = s.holidayRepo.GetByEpoch(date.Unix())
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
func (s *Service) GetHolidayByID(id int) (*model.Holiday, error) {
	holiday, err := s.holidayRepo.GetByID(id)
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
func (s *Service) UpdateByID(id int, holiday *model.Holiday) error {
	rowsAffected, err := s.holidayRepo.UpdateNameAndType(id, holiday.Name, holiday.HolidayType)
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
func (s *Service) DeleteByID(id int) error {
	rowsAffected, err := s.holidayRepo.DeletedByID(id)
	if err != nil {
		log.Errorf("Database Error : %v", err)
		return util.ErrDatabaseError
	} else if rowsAffected == 0 {
		log.Errorf("no holiday found with id: %v", err)
		return util.ErrResourceNotFound
	}

	return nil
}
