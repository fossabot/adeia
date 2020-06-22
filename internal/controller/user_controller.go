package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/util"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() []*route.Route {
	routes := []*route.Route{
		route.New(http.MethodPost, "/users/", CreateUser(), middleware.Nil), // create new user
		route.New(http.MethodGet, "/users/:id", GetUser(), middleware.Nil), // get user
	}
	return routes
}

func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := httprouter.ParamsFromContext(r.Context()).ByName("id")

		usr, err := usrSvc.GetUserByID(id)
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
		EmployeeID  string `json:"employee_id"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
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
		w.WriteHeader(http.StatusCreated)
	}
}
