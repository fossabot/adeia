package controller

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	"adeia-api/internal/service"
	holidaySvc "adeia-api/internal/service/holiday"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
	"adeia-api/internal/util/mail"
)

var (
	usrSvc         userService.Service
	holidayService service.HolidayService
	session        sessionService.Service
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache, m mail.Mailer) {
	usrSvc = userService.New(d, c, m)
	session = sessionService.NewService(c.GetInstance())
	holidayService = holidaySvc.New(d, c)
}
