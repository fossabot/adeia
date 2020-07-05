package controller

import (
	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/model"
	"adeia-api/internal/util"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// HolidayRoutes returns a slice containing all holiday-related routes.
func HolidayRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/holidays/", CreateHoliday(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month", GetHolidayByYearAndMonth(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year", GetHolidayByYear(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/year/:year/month/:month/date/:date", GetHolidayByDate(), middleware.Nil),
		route.New(http.MethodGet, "/holidays/id/:id", GetHolidayById(), middleware.Nil),
	}
	return routes
}

func CreateHoliday() http.HandlerFunc {
	validator := func(holiday model.Holiday) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"date": validation.Validate(holiday.HolidayDate.String(),
					validation.Required,
					//validation.Date(time.RFC3339),
				),
				"holiday_name": validation.Validate(holiday.Name,
					validation.Required,
					validation.RuneLength(4, 128),
					is.Alphanumeric,
				),
				"holiday_type": validation.Validate(holiday.HolidayType,
					validation.Required,
					validation.RuneLength(4, 128),
					is.UTFLetterNumeric,
				),
			},
		}
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		var holiday model.Holiday
		if err := util.DecodeBodyAndRespond(writer, request, &holiday); err != nil {
			return
		}
		// validate request
		if err := validator(holiday).Validate(); err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}

		// Call holiday service to create
		response , err := holidaySvc.CreateHoliday(holiday)
		if err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}

		writer.Header().Set("Location", "/holidays/"+string(holiday.ID))
		util.RespondWithJSON(writer, http.StatusCreated, response)
	}
}

func GetHolidayByYear() http.HandlerFunc {
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
	return func(writer http.ResponseWriter, request *http.Request) {
		param := httprouter.ParamsFromContext(request.Context()).ByName("year")
		if err := validator(param).Validate(); err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		year,_ := strconv.Atoi(param)
		holidays, err := holidaySvc.GetHolidayByDate(util.GetTime(year,1,1), model.Year)
		if err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(writer, http.StatusOK, holidays)
	}
}

func GetHolidayByYearAndMonth() http.HandlerFunc {
	validator := func(year, month string) *util.Validation {
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
			},
		}
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		if err := validator(params.ByName("year"),params.ByName("month")).Validate(); err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		year, _ := strconv.Atoi(params.ByName("year"))
		month, _ := strconv.Atoi(params.ByName("month"))


		holidays, err := holidaySvc.GetHolidayByDate(util.GetTime(year,month,1), model.Month)
		if err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(writer, http.StatusOK, holidays)
	}
}

func GetHolidayByDate() http.HandlerFunc {
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
	return func(writer http.ResponseWriter, request *http.Request) {
		params := httprouter.ParamsFromContext(request.Context())
		if err := validator(params.ByName("year"),params.ByName("month"),params.ByName("date")).Validate(); err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}

		year, _ := strconv.Atoi(params.ByName("year"))
		month, _ := strconv.Atoi(params.ByName("month"))
		date, _ := strconv.Atoi(params.ByName("date"))

		holidays, err := holidaySvc.GetHolidayByDate(util.GetTime(year,month,date), model.DateOfMonth)
		if err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(writer, http.StatusOK, holidays)
	}
}


func GetHolidayById() http.HandlerFunc {
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
	return func(writer http.ResponseWriter, request *http.Request) {
		id := httprouter.ParamsFromContext(request.Context()).ByName("id")
		if err := validator(id).Validate(); err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		holidayId , _ := strconv.Atoi(id)
		holiday, err := holidaySvc.GetHolidayById(holidayId)
		if err != nil {
			util.RespondWithError(writer, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(writer, http.StatusOK, holiday)
	}
}
