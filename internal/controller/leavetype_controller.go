package controller

import (
	"adeia-api/internal/model"
	"adeia-api/internal/util"
	"adeia-api/internal/util/validation"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

// LeaveTypeRoutes returns a slice containing all leave-type-related routes.
func LeaveTypeRoutes() (string, chi.Router) {
	r := chi.NewRouter()

	r.Post("/", CreateLeaveType())

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", GetLeaveTypeById())
		r.Put("/", UpdateLeaveTypeByID())
		r.Delete("/", DeleteLeaveTypeByID())
	})

	return "/leavetype", r
}

func DeleteLeaveTypeByID() http.HandlerFunc {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateResourceID(id),
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
		if err := leaveTypeSvc.DeleteByID(holidayID); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusNoContent, nil)
	}
}

func UpdateLeaveTypeByID() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Days string `json:"days"`
	}

	validator := func(id, name, days string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id":   validation.ValidateResourceID(id),
				"name": validation.ValidateResourceName(name),
				"days": validation.ValidateResourceID(days),
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
		if err := validator(id, rBody.Name, rBody.Days).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		leaveTypeID, _ := strconv.Atoi(id)
		days, _ := strconv.Atoi(id)
		leaveType := model.LeaveType{Days: days, Name: rBody.Name}
		if err := leaveTypeSvc.UpdateByID(leaveTypeID, &leaveType); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		leaveType.ID = leaveTypeID
		util.RespondWithJSON(w, http.StatusOK, leaveType)
	}
}

func GetLeaveTypeById() http.HandlerFunc {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateResourceID(id),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := validator(id).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		leaveTypeId, _ := strconv.Atoi(id)
		leaveType, err := leaveTypeSvc.GetLeaveTypeByID(leaveTypeId)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}
		util.RespondWithJSON(w, http.StatusOK, leaveType)
	}
}

func CreateLeaveType() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Days string `json:"days"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"name": validation.ValidateResourceName(r.Name),
				"days": validation.ValidateResourceID(r.Days),
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

		days, _ := strconv.Atoi(rBody.Days)
		leaveType := model.LeaveType{Name: rBody.Name, Days: days}
		response, err := leaveTypeSvc.CreateLeaveType(leaveType)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		w.Header().Set("Location", "/leavetype/"+string(response.ID))
		util.RespondWithJSON(w, http.StatusCreated, response)
	}
}
