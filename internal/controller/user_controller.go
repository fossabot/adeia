package controller

import (
	"net/http"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/util"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julienschmidt/httprouter"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/users", CreateUser(), middleware.Nil), // create new user
		route.New(http.MethodGet, "/users/:id", GetUser(), middleware.Nil), // get user
	}
	return routes
}

// GetUser returns a new user using the employee_id.
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// from outside, the id appears to be the primary key.
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		usr, err := usrSvc.GetUserByEmpID(id)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		util.RespondWithJSON(w, http.StatusOK, usr)
	}
}

// CreateUser creates a new user.
func CreateUser() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		EmployeeID  string `json:"employee_id,omitempty"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
	}

	validator := func(r request) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"name": validation.Validate(
					r.Name,
					validation.Required,
					validation.Length(2, 255),
				),
				"employee_id": validation.Validate(
					r.EmployeeID,
					validation.Required.When(r.EmployeeID != ""), // employee_id is optional
					validation.Length(5, 10),
					is.Alphanumeric,
				),
				"email": validation.Validate(
					r.Email,
					validation.Required,
					validation.Length(5, 120),
					is.EmailFormat,
				),
				"designation": validation.Validate(
					r.Designation,
					validation.Required,
					validation.Length(1, 255),
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decode request body
		var rBody request
		if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
			return
		}

		// validate request
		if err := validator(rBody).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// create user
		empID, err := usrSvc.CreateUser(
			rBody.Name,
			rBody.Email,
			rBody.EmployeeID,
			rBody.Designation,
		)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// return response
		w.Header().Set("Location", "/v1/users/"+empID)
		w.WriteHeader(http.StatusCreated)
	}
}
