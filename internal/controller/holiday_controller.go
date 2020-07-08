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

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// HolidayRoutes returns a slice containing all holiday-related routes.
func HolidayRoutes() []*route.Route {
	return []*route.Route{
		route.New(http.MethodPost, "/holidays", CreateHoliday(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/{year}/month/{month}", GetHolidaysByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/{year}", GetHolidaysByYear(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/{year}/month/{month}/date/{date}", GetHolidaysByDate(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/{id}", GetHolidayByID(), middleware.Nil),
		route.New(http.MethodPut, "/holidays/{id}", UpdateHolidayByID(), middleware.Nil),
		route.New(http.MethodDelete, "/holidays/{id}", DeleteHolidayByID(), middleware.Nil),
	}
}

// CreateHoliday creates a new holiday.
func CreateHoliday() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Date string `json:"date"`
	}

	validator := func(r request) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"name": validation.Validate(r.Name,
					validation.Required,
					validation.RuneLength(4, 128),
					is.Alphanumeric,
				),
				"type": validation.Validate(r.Type,
					validation.Required,
					validation.RuneLength(4, 128),
					is.UTFLetterNumeric,
				),
				"date": validation.Validate(r.Date,
					validation.Required,
					validation.Date(time.RFC3339),
				),
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
	validator := func(year string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"year": validation.Validate(year,
					validation.Required,
					is.UTFNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		year := chi.URLParam(r, "year")
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
	validator := func(year, month string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"year": validation.Validate(year,
					validation.Required,
					is.UTFNumeric,
				),
				"month": validation.Validate(month,
					validation.Required,
					is.UTFNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
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

// GetHolidaysByDate returns the holiday by the provided date.
func GetHolidaysByDate() http.HandlerFunc {
	validator := func(year, month, date string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"year": validation.Validate(year,
					validation.Required,
					is.UTFNumeric,
				),
				"month": validation.Validate(year,
					validation.Required,
					is.UTFNumeric,
				),
				"date": validation.Validate(year,
					validation.Required,
					is.UTFNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
		date := chi.URLParam(r, "date")
		if err := validator(year, month, date).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		y, _ := strconv.Atoi(year)
		m, _ := strconv.Atoi(month)
		d, _ := strconv.Atoi(date)
		holidays, err := holidaySvc.GetHolidaysByDate(util.GetTime(y, m, d), constants.DateOfMonth)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, holidays)
	}
}

// GetHolidayByID gets a holiday by id.
func GetHolidayByID() http.HandlerFunc {
	validator := func(id string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					is.UTFNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
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
	validator := func(id string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					is.UTFNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
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

	validator := func(id, name, holidayType string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					is.UTFNumeric,
				),
				"name": validation.Validate(name,
					validation.Required,
					validation.RuneLength(4, 128),
					is.Alphanumeric,
				),
				"type": validation.Validate(holidayType,
					validation.Required,
					validation.RuneLength(4, 128),
					is.UTFLetterNumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
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
