package controller

import (
	"adeia-api/internal/service/cache"
	"adeia-api/internal/service/db"
	userService "adeia-api/internal/service/user"
)

var (
	usrSvc userService.Service
)

func Init(d db.DB, c cache.Cache) {
	usrSvc = userService.New(d, c)
}
