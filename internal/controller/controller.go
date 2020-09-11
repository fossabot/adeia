package controller

import (
	"adeia-api/internal/config"
	"adeia-api/internal/repo"
	holidayRepo "adeia-api/internal/repo/holiday"
	roleRepo "adeia-api/internal/repo/role"
	sessionRepo "adeia-api/internal/repo/session"
	userRepo "adeia-api/internal/repo/user"
	"adeia-api/internal/service"
	holidayService "adeia-api/internal/service/holiday"
	roleService "adeia-api/internal/service/role"
	sessionService "adeia-api/internal/service/session"
	userService "adeia-api/internal/service/user"
)

var (
	HolidaySvc HolidayService
	RoleSvc    RoleService
	SessionSvc SessionService
	UserSvc    UserService
)

// Init initializes all services that are used in the controllers.
func Init(conf *config.Config, d repo.DB, c service.Cache, m service.Mailer) {
	// init repository layer
	h := holidayRepo.New(d)
	r := roleRepo.New(d)
	s := sessionRepo.New(d)
	u := userRepo.New(d)

	// init service layer
	HolidaySvc = holidayService.New(h)
	RoleSvc = roleService.New(r)
	SessionSvc = sessionService.New(s, conf.ServerConfig.JWTSecret)
	UserSvc = userService.New(u, c, m)
}
