package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/model"
	"adeia-api/internal/util/log"

	"github.com/julienschmidt/httprouter"
)

// HolidayRoutes returns a slice containing all holiday-related routes.
func HolidayRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/holidays/", CreateHoliday(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month", GetHolidayByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year", GetHolidayByYear(), middleware.Nil),
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
		if err := holidaySvc.CreateHoliday(holiday); err != nil {
			log.Errorf("cannot create new holiday: %v", err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Location", "http://google.com")
		writer.WriteHeader(http.StatusCreated)

	}
}

func GetHolidayByYear() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		year, _ := strconv.Atoi(params.ByName("year"))
		holidays, err := holidaySvc.GetHolidayByYear(year)
		if err != nil {
			log.Error(err)
			_, _ = writer.Write(nil)
		}
		result, _ := json.Marshal(holidays)
		_, _ = writer.Write(result)
	}
}

func GetHolidayByYearAndMonth() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		year, _ := strconv.Atoi(params.ByName("year"))
		month, _ := strconv.Atoi(params.ByName("month"))
		holidays, err := holidaySvc.GetHolidayByYearAndMonth(year, month)
		if err != nil {
			log.Error(err)
			_, _ = writer.Write(nil)
		}
		result, _ := json.Marshal(holidays)
		_, _ = writer.Write(result)
	}
}
