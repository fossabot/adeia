package model

import "time"

// Holiday represents the Holiday model.
type Holiday struct {
	// ID represents the holiday ID. It is an auto-incremented db field.
	ID int `db:"id" json:"id"`

	// Name represents the name of the holiday.
	Name string `db:"name" json:"name"`

	// HolidayType represents the type of the holiday.
	HolidayType string `db:"type" json:"type"`

	// HolidayDate represents the date on which the holiday occurs.
	HolidayDate time.Time `db:"date" json:"date"`
}
