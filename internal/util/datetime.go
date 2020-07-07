package util

import "time"

// GetYMDFromTime returns year, month and date from the provided time.Time.
func GetYMDFromTime(time time.Time) (year, month, day int) {
	return time.Year(), int(time.Month()), time.Day()
}

// GetTime returns time.Time from the provided year, month and date.
func GetTime(year, month, date int) time.Time {
	return time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.UTC)
}
