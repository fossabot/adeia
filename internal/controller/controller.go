package controller

import (
	"net/http"

	"adeia-api/internal/cache"
	"adeia-api/internal/db"
	holidayService "adeia-api/internal/service/holiday"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
	"adeia-api/internal/util"
	"adeia-api/internal/util/log"
	"adeia-api/internal/util/mail"
)

type ProtectedHandler struct {
	PermissionName string
	Handler        http.HandlerFunc
}

func (p *ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// perform role checking using context, like !contains(userRoles, p.PermissionName)
	if false {
		// user doesn't have access
		log.Debug("user not authorized to perform action")
		util.RespondWithError(w, util.ErrUnauthorized)
		return
	}

	// user has access, so continue
	p.Handler.ServeHTTP(w, r)
}

var (
	usrSvc     userService.Service
	holidaySvc holidayService.Service
	sessionSvc sessionService.Service
)

// Init initializes all services that are used in the controllers.
func Init(d db.DB, c cache.Cache, m mail.Mailer) {
	usrSvc = userService.New(d, c, m)
	sessionSvc = sessionService.New(c.GetInstance())
	holidaySvc = holidayService.New(d, c)
}
