package controller

import (
	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/model"
	log "adeia-api/internal/util/logger"
	"encoding/json"
	"net/http"
)

// UserRoutes returns a slice containing all user-related routes.
func HolidayRoutes() []*route.Route {
	routes := []*route.Route{
		// create new user
		route.New(http.MethodPost, "/holidays/", CreateHoliday(), middleware.Nil),
	}
	return routes
}

func CreateHoliday() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var holiday model.Holiday
		if err := json.NewDecoder(request.Body).Decode(&holiday); err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Error(err)
			return
		}
		// TODO: perform validation
		// create user
		if err := holidayService.CreateHoliday(holiday); err != nil {
			log.Errorf("cannot create new holiday: %v", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Location", "http://google.com")
		writer.WriteHeader(http.StatusCreated)

	}
}