package controller

import (
	"encoding/json"
	"net/http"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	log "adeia-api/internal/util/logger"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() []*route.Route {
	routes := []*route.Route{
		// create new user
		route.New(http.MethodPost, "/users/", CreateUser(), middleware.Nil),
	}
	return routes
}

// CreateUser creates a new user.
func CreateUser() http.HandlerFunc {
	type createUserRequest struct {
		Name        string `json:"name"`
		EmployeeID  string `json:"employee_id"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
	}

	type createUserResponse struct {
		Location string `json:"location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// decode request body
		var rBody createUserRequest
		if err := json.NewDecoder(r.Body).Decode(&rBody); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// TODO: perform validation
		// create user
		if err := usrSvc.CreateUser(rBody.Name, rBody.Email, rBody.EmployeeID, rBody.Designation); err != nil {
			log.Errorf("cannot create new user: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		resp := &createUserResponse{"http://api.example.com/v1/users/" + rBody.EmployeeID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
