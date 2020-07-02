package controller

import (
	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/model"
	log "adeia-api/internal/util/logger"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// UserRoutes returns a slice containing all user-related routes.
func HolidayRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/holidays/", CreateHoliday(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month", GetHolidayByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year", GetHolidayByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month/date/:date", GetHolidayByDate(), middleware.Nil),
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

func GetHolidayByYear() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		year, _ := strconv.Atoi(params.ByName("year"))
		holidays, err := holidayService.GetHolidayByDate(model.Date{Year:year}, model.Year)
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
		holidays, err := holidayService.GetHolidayByDate(model.Date{Month:month,Year:year}, model.Month)
		if err != nil {
			log.Error(err)
			_, _ = writer.Write(nil)
		}
		result, _ := json.Marshal(holidays)
		_, _ = writer.Write(result)
	}
}

func GetHolidayByDate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		year, _ := strconv.Atoi(params.ByName("year"))
		month, _ := strconv.Atoi(params.ByName("month"))
		date, _ := strconv.Atoi(params.ByName("date"))
		holidays, err := holidayService.GetHolidayByDate(model.Date{DayOfMonth:date,Month:month,Year:year}, model.Year)
		if err != nil {
			log.Error(err)
			_, _ = writer.Write(nil)
		}
		result, _ := json.Marshal(holidays)
		_, _ = writer.Write(result)
	}
}