package model

import "time"

type Holiday struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	HolidayType string    `db:"type" json:"type"`
	HolidayDate time.Time `db:"date" json:"date"`
}

type Date struct {
	DayOfMonth int
	Month int
	Year int
	Epoch int64
}

type TimeUnit string

 const (
 	 Epoch TimeUnit = "Epoch"
	 Month TimeUnit = "Month"
	 Year TimeUnit = "Year"
 )