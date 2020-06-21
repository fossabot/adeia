 package model

 import "time"

 type Holiday struct {
	ID int `db:"id"`
	Name string `db:"name"`
	HolidayType string `db:"type"`
	HolidayDate time.Time `db:"date"`
}