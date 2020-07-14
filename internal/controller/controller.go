package controller

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	holidayService "adeia-api/internal/service/holiday"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
	"adeia-api/internal/util/mail"
)

var (
	usrSvc     userService.Service
	holidaySvc holidayService.Service
	sessionSvc sessionService.Service
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache, m mail.Mailer) {
	usrSvc = userService.New(d, c, m)
	sessionSvc = sessionService.New(d)
	holidaySvc = holidayService.New(d, c)
}
