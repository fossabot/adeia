package util

import "time"

func GetYMDFromTime(time time.Time) (year, month, day int){
	return time.Year(), int(time.Month()), time.Day()
}


func GetTime(year, month, date int) time.Time {
	return time.Date(year,time.Month(month),date,0,0,0,0,time.UTC)
}