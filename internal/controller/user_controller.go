package controller

import (
	"net/http"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/util"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/users/", CreateUser(), middleware.Nil), // create new user
	}
	return routes
}

// CreateUser creates a new user.
func CreateUser() http.HandlerFunc {
	type request struct {
		Name        string `json:"name"`
		EmployeeID  string `json:"employee_id"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
	}

	type response struct {
		Location string `json:"location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decode request body
		var rBody request
		if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
			return
		}

		// TODO: perform validation
		// create user
		if err := usrSvc.CreateUser(rBody.Name, rBody.Email, rBody.EmployeeID, rBody.Designation); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// return response
		w.Header().Set("Location", "/v1/users/"+rBody.EmployeeID)
		util.RespondWithJSON(w, http.StatusCreated, &struct{}{})
	}
}
