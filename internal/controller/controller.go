package controller

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/service"
	holidaySvc "adeia-api/internal/service/holiday"
	userService "adeia-api/internal/service/user"
)

var (
	usrSvc userService.Service
	holidayService service.HolidayService
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache) {
	usrSvc = userService.New(d, c)
	holidayService = holidaySvc.New(d,c)
}
