package controller

import (
	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	holidayService "adeia-api/internal/service/holiday"
	leavetypeService "adeia-api/internal/service/leavetype"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
	"adeia-api/internal/util/mail"
)

var (
	usrSvc       userService.Service
	holidaySvc   holidayService.Service
	sessionSvc   sessionService.Service
	leaveTypeSvc leavetypeService.Service
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache, m mail.Mailer) {
	usrSvc = userService.New(d, c, m)
	sessionSvc = sessionService.New(c.GetInstance())
	holidaySvc = holidayService.New(d, c)
	leaveTypeSvc = leavetypeService.New(d, c)
}
