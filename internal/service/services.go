package service

import "adeia-api/internal/model"

// Service contains all user-related business logic.
type UserService interface {
	CreateUser(name, email, empID, designation string) error
}

type HolidayService interface {
	CreateHoliday(holiday model.Holiday) error
	GetHolidayByYear(year int) (*[]model.Holiday, error)
	GetHolidayByYearAndMonth(year, month int) (*[]model.Holiday, error)
}
