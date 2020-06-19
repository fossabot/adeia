package controller

import (
	"net/http"

	"adeia-api/internal/middleware"
	"adeia-api/internal/route"
	log "adeia-api/internal/utils/logger"
)

func UserRoutes() []*route.Route {
	routes := []*route.Route{
		// create new user
		route.New(http.MethodPost, "/users/", CreateUser(), middleware.Nil),
	}
	return routes
}

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: perform validation

		// call service
		_, err := usrSvc.CreateUser()
		if err != nil {
			log.Errorf("cannot create new user: %v", err)
		}
	}
}
