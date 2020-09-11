package controller

import (
	"adeia/internal/config"
	"adeia/internal/repo"
	holidayRepo "adeia/internal/repo/holiday"
	roleRepo "adeia/internal/repo/role"
	sessionRepo "adeia/internal/repo/session"
	userRepo "adeia/internal/repo/user"
	"adeia/internal/service"
	holidayService "adeia/internal/service/holiday"
	roleService "adeia/internal/service/role"
	sessionService "adeia/internal/service/session"
	userService "adeia/internal/service/user"
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
