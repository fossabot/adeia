package controller

import (
	"net/http"
	"strconv"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/util"
	"adeia-api/internal/util/validation"

	"github.com/go-chi/chi"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() (string, chi.Router) {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", CreateUser())
	r.Method(http.MethodPatch, "/activation", ActivateUser())
	r.Method(http.MethodPost, "/sessions", LoginUser())

	r.Route("/{id}", func(r chi.Router) {
		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowAuthenticated(sessionSvc, usrSvc, true))
			r.Method(http.MethodGet, "/", GetUser())
			r.Method(http.MethodDelete, "/", DeleteUser())
			r.Method(http.MethodPost, "/sessions/destroy", LogoutUser())
		})
	})

	return "/users", r
}

// LoginUser logs in an user.
func LoginUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"email":    validation.ValidateEmail(r.Email),
				"password": validation.ValidateLoginPwd(r.Password),
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

		// check credentials
		usr, err := usrSvc.LoginUser(rBody.Email, rBody.Password)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// create new session and add to cookie
		if err := sessionSvc.Create(w, strconv.Itoa(usr.ID)); err != nil {
			util.RespondWithError(w, util.ErrInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// LogoutUser logs-out an user.
func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := sessionSvc.Destroy(w, r); err != nil {
			util.RespondWithError(w, util.ErrInternalServerError)
			return
		}
		util.RespondWithJSON(w, http.StatusNoContent, nil)
	}
}

// DeleteUser deletes an user account.
func DeleteUser() *ProtectedHandler {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateEmpID(id),
			},
		}
	}

	return &ProtectedHandler{
		PermissionName: "DELETE_USERS",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// from outside, the id appears to be the primary key.
			id := chi.URLParam(r, "id")

			// validate request
			if err := validator(id).Validate(); err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			// delete user
			if err := usrSvc.DeleteUser(id); err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			w.WriteHeader(http.StatusNoContent)
		},
	}
}

// ActivateUser activates an user account.
func ActivateUser() http.HandlerFunc {
	type request struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	validator := func(r request, id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id":               validation.ValidateEmpID(id),
				"email":            validation.ValidateEmail(r.Email),
				"password":         validation.ValidateNewPwd(r.Password),
				"confirm_password": validation.ValidateConfirmPwd(r.ConfirmPassword, r.Password),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		empID := chi.URLParam(r, "id")
		// decode request body
		var rBody request
		if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
			return
		}

		// validate request
		if err := validator(rBody, empID).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// activate user
		usr, err := usrSvc.ActivateUser(empID, rBody.Email, rBody.Password)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// return response
		util.RespondWithJSON(w, http.StatusOK, usr)
	}
}

// CreateUser creates a new user.
func CreateUser() *ProtectedHandler {
	type request struct {
		Name        string `json:"name"`
		EmployeeID  string `json:"employee_id,omitempty"`
		Email       string `json:"email"`
		Designation string `json:"designation"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"name":        validation.ValidateName(r.Name),
				"employee_id": validation.ValidateEmpIDOptional(r.EmployeeID), // employee_id is optional
				"email":       validation.ValidateEmail(r.Email),
				"designation": validation.ValidateDesignation(r.Designation),
			},
		}
	}

	return &ProtectedHandler{
		PermissionName: "CREATE_USERS",
		Handler: func(w http.ResponseWriter, r *http.Request) {
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
			usr, err := usrSvc.CreateUser(
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
			w.Header().Set("Location", "/v1/users/"+usr.EmployeeID)
			util.RespondWithJSON(w, http.StatusCreated, usr)
		},
	}
}

// GetUser gets the user using the employee_id.
func GetUser() *ProtectedHandler {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateEmpID(id),
			},
		}
	}

	return &ProtectedHandler{
		PermissionName: "GET_USERS",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// from outside, the id appears to be the primary key.
			id := chi.URLParam(r, "id")

			// validate request
			if err := validator(id).Validate(); err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			// get user
			usr, err := usrSvc.GetUserByEmpID(id)
			if err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}
			util.RespondWithJSON(w, http.StatusOK, usr)
		},
	}
}
