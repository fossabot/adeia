package controller

import (
	"errors"
	"net/http"
	"strconv"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/api/route"
	"adeia-api/internal/util"
	"adeia-api/internal/util/crypto"

	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// UserRoutes returns a slice containing all user-related routes.
func UserRoutes() []*route.Route {
	allowAuthenticated := middleware.NewChain(middleware.AllowAuthenticated(sessionSvc, usrSvc, true))
	allowUnauthenticated := middleware.NewChain(middleware.AllowAuthenticated(sessionSvc, usrSvc, false))

	return []*route.Route{
		// create new user
		route.New(http.MethodPost, "/users", CreateUser(), middleware.Nil),
		// get user
		route.New(http.MethodGet, "/users/{id}", GetUser(), allowAuthenticated),
		// delete user
		route.New(http.MethodDelete, "/users/{id}", DeleteUser(), allowAuthenticated),
		// activate user
		route.New(http.MethodPatch, "/users/{id}/activation", ActivateUser(), allowUnauthenticated),
		// login user
		route.New(http.MethodPost, "/users/{id}/sessions", LoginUser(), allowUnauthenticated),
		// logout user
		route.New(http.MethodPost, "/users/{id}/sessions/destroy", LogoutUser(), allowAuthenticated),
	}
}

// LoginUser logs in an user.
func LoginUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	validator := func(r request) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"email": validation.Validate(r.Email,
					validation.Required,
					validation.RuneLength(3, 120),
					is.EmailFormat,
				),
				"password": validation.Validate(r.Password,
					validation.Required,
					validation.RuneLength(12, 128),
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
func DeleteUser() http.HandlerFunc {
	validator := func(id string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					validation.RuneLength(5, 10),
					is.Alphanumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

// ActivateUser activates an user account.
func ActivateUser() http.HandlerFunc {
	type request struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	validator := func(r request, id string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					validation.RuneLength(5, 10),
					is.Alphanumeric,
				),
				"email": validation.Validate(r.Email,
					validation.Required,
					validation.RuneLength(3, 120),
					is.EmailFormat,
				),
				"password": validation.Validate(r.Password,
					validation.Required,
					validation.RuneLength(12, 128),
					validation.By(func(value interface{}) error {
						s, _ := value.(string)
						if crypto.PasswordStrength(s) < 3 {
							return errors.New("password is weak")
						}
						return nil
					}),
				),
				"confirm_password": validation.Validate(r.ConfirmPassword,
					validation.Required,
					validation.RuneLength(12, 128),
					validation.By(func(value interface{}) error {
						s, _ := value.(string)
						if s != r.Password {
							return errors.New("passwords do not match")
						}
						return nil
					}),
				),
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
				"name": validation.Validate(r.Name,
					validation.Required,
					validation.RuneLength(2, 255),
				),
				"employee_id": validation.Validate(r.EmployeeID,
					validation.Required.When(r.EmployeeID != ""), // employee_id is optional
					validation.RuneLength(5, 10),
					is.Alphanumeric,
				),
				"email": validation.Validate(r.Email,
					validation.Required,
					validation.RuneLength(3, 120),
					is.EmailFormat,
				),
				"designation": validation.Validate(r.Designation,
					validation.Required,
					validation.RuneLength(1, 255),
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
	}
}

// GetUser gets the user using the employee_id.
func GetUser() http.HandlerFunc {
	validator := func(id string) *util.Validation {
		return &util.Validation{
			Errors: validation.Errors{
				"id": validation.Validate(id,
					validation.Required,
					validation.RuneLength(5, 10),
					is.Alphanumeric,
				),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
