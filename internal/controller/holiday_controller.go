package controller

import (
	"net/http"
	"strconv"
	"time"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/model"
	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/validation"

	"github.com/arkn98/httprouter"
)

// HolidayRoutes returns a slice containing all holiday-related routes.
func HolidayRoutes() []*route.Route {
	return []*route.Route{
		route.New(http.MethodPost, "/holidays", CreateHoliday(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month", GetHolidaysByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year", GetHolidaysByYear(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month/day/:day", GetHolidaysByDay(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/id/:id", GetHolidayByID(), middleware.Nil),
		route.New(http.MethodPut, "/holidays/id/:id", UpdateHolidayByID(), middleware.Nil),
		route.New(http.MethodDelete, "/holidays/id/:id", DeleteHolidayByID(), middleware.Nil),
	}
}

// CreateHoliday creates a new holiday.
func CreateHoliday() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Date string `json:"date"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"name": validation.ValidateResourceName(r.Name),
				"type": validation.ValidateResourceName(r.Type),
				"date": validation.ValidateDate(r.Date),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var rBody request
		if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
			return
		}

		// validate request
		if err := validator(rBody).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// create holiday
		date, _ := time.Parse(time.RFC3339, rBody.Date)
		holiday := model.Holiday{Name: rBody.Name, HolidayType: rBody.Type, HolidayDate: date}
		response, err := holidaySvc.CreateHoliday(holiday)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		w.Header().Set("Location", "/holidays/"+string(response.ID))
		util.RespondWithJSON(w, http.StatusCreated, response)
	}
}

// GetHolidaysByYear returns the holidays in the provided year.
func GetHolidaysByYear() http.HandlerFunc {
	validator := func(year string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"year": validation.ValidateYear(year),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		year := httprouter.ParamsFromContext(r.Context()).ByName("year")
		if err := validator(year).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// get holidays
		y, _ := strconv.Atoi(year)
		holidays, err := holidaySvc.GetHolidaysByDate(util.GetTime(y, 1, 1), constants.Year)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, holidays)
	}
}

// GetHolidaysByYearAndMonth returns the holidays in the provided year and month.
func GetHolidaysByYearAndMonth() http.HandlerFunc {
	validator := func(year, month string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"year":  validation.ValidateYear(year),
				"month": validation.ValidateMonth(month),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		year := params.ByName("year")
		month := params.ByName("month")
		if err := validator(year, month).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// get holidays
		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		holidays, err := holidaySvc.GetHolidaysByDate(util.GetTime(y, m, 1), constants.Month)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, holidays)
	}
}

// GetHolidaysByDay returns the holiday by the provided date.
func GetHolidaysByDay() http.HandlerFunc {
	validator := func(year, month, day string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"year":  validation.ValidateYear(year),
				"month": validation.ValidateMonth(month),
				"day":   validation.ValidateDay(day),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		year := params.ByName("year")
		month := params.ByName("month")
		day := params.ByName("day")
		if err := validator(year, month, day).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		d, _ := strconv.Atoi(day)
		holidays, err := holidaySvc.GetHolidaysByDate(util.GetTime(y, m, d), constants.DayOfMonth)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, holidays)
	}
}

// GetHolidayByID gets a holiday by id.
func GetHolidayByID() http.HandlerFunc {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateResourceID(id),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		if err := validator(id).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		holidayID, _ := strconv.Atoi(id)
		holiday, err := holidaySvc.GetHolidayByID(holidayID)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, holiday)
	}
}

// DeleteHolidayByID deletes a holiday by id.
func DeleteHolidayByID() http.HandlerFunc {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateResourceID(id),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		if err := validator(id).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// delete holiday
		holidayID, _ := strconv.Atoi(id)
		if err := holidaySvc.DeleteByID(holidayID); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusNoContent, nil)
	}
}

// UpdateHolidayByID updates a holiday.
func UpdateHolidayByID() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}

	validator := func(id, name, holidayType string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id":   validation.ValidateResourceID(id),
				"name": validation.ValidateResourceName(name),
				"type": validation.ValidateResourceName(holidayType),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		var rBody request
		if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
			return
		}

		// validate request
		if err := validator(id, rBody.Name, rBody.Type).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// update holiday
		holidayID, _ := strconv.Atoi(id)
		holiday := model.Holiday{HolidayType: rBody.Type, Name: rBody.Name}
		if err := holidaySvc.UpdateByID(holidayID, &holiday); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		holiday.ID = holidayID
		util.RespondWithJSON(w, http.StatusOK, holiday)
	}
}
