package user

import (
	"net/http"

	"adeia/internal/api/middleware"
	"adeia/internal/controller"
	"adeia/internal/model"
	"adeia/internal/util"
	"adeia/internal/util/constants"
	"adeia/internal/util/validation"

	"github.com/go-chi/chi"
)

// Routes returns a slice containing all user-related routes.
func Routes() (string, chi.Router) {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", CreateUser())
	r.Method(http.MethodPatch, "/activation", ActivateUser())
	r.Method(http.MethodPost, "/sessions", LoginUser())

	r.Group(func(r chi.Router) {
		// protected routes
		r.Use(middleware.AllowAuthenticated(controller.SessionSvc, controller.UserSvc))

		r.Method(http.MethodPost, "/sessions/refresh", RefreshToken())
		r.Method(http.MethodPost, "/sessions/destroy", LogoutUser())
		r.Route("/{id}", func(r chi.Router) {
			r.Method(http.MethodGet, "/", GetUser())
			r.Method(http.MethodDelete, "/", DeleteUser())
		})
	})

	return "/users", r
}

// RefreshToken refreshes the access and refresh token. The user must be authenticated
// before they can hit this route.
func RefreshToken() http.HandlerFunc {
	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	validator := func(refreshToken string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"refreshToken": validation.ValidateToken(refreshToken),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := controller.SessionSvc.ReadRefreshTokenCookie(r)
		if err != nil {
			util.RespondWithError(w, util.ErrUnauthorized)
			return
		}

		// validate
		if err := validator(refreshToken).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// get user from context
		user, ok := r.Context().Value(constants.ContextUserKey).(*model.User)
		if !ok {
			util.RespondWithError(w, util.ErrUnauthorized)
			return
		}

		// refresh the tokens
		newAccessToken, newRefreshToken, err := controller.SessionSvc.RefreshToken(user.ID, user.EmployeeID, refreshToken)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// set new refreshToken cookie
		controller.SessionSvc.AddRefreshTokenCookie(w, newRefreshToken)

		// send new access token
		resp := response{
			AccessToken: newAccessToken,
			ExpiresIn:   constants.AccessTokenExpiry,
		}
		util.RespondWithJSON(w, http.StatusOK, &resp)
	}
}

// LogoutUser logs out a user.
func LogoutUser() http.HandlerFunc {
	validator := func(refreshToken string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"refreshToken": validation.ValidateToken(refreshToken),
			},
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// get refresh token from cookie
		refreshToken, err := controller.SessionSvc.ReadRefreshTokenCookie(r)
		if err != nil {
			util.RespondWithError(w, util.ErrUnauthorized)
			return
		}

		// validate
		if err := validator(refreshToken).Validate(); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// get user from context
		user, ok := r.Context().Value(constants.ContextUserKey).(*model.User)
		if !ok {
			util.RespondWithError(w, util.ErrUnauthorized)
			return
		}

		// destroy session
		if err := controller.SessionSvc.Destroy(user.ID, refreshToken); err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		util.RespondWithJSON(w, http.StatusNoContent, nil)
	}
}

// LoginUser logs in an user.
func LoginUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
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
		usr, err := controller.UserSvc.LoginUser(rBody.Email, rBody.Password)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// create new session and add to cookie
		accessToken, refreshToken, err := controller.SessionSvc.NewSession(usr.ID, usr.EmployeeID)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// set refreshToken in cookie
		controller.SessionSvc.AddRefreshTokenCookie(w, refreshToken)

		// send access token
		resp := response{
			AccessToken: accessToken,
			ExpiresIn:   constants.AccessTokenExpiry,
		}
		util.RespondWithJSON(w, http.StatusOK, &resp)
	}
}

// DeleteUser deletes an user account.
func DeleteUser() *controller.ProtectedHandler {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateEmpID(id),
			},
		}
	}

	return &controller.ProtectedHandler{
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
			if err := controller.UserSvc.DeleteUser(id); err != nil {
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
		EmployeeID      string `json:"id"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id":               validation.ValidateEmpID(r.EmployeeID),
				"email":            validation.ValidateEmail(r.Email),
				"password":         validation.ValidateNewPwd(r.Password),
				"confirm_password": validation.ValidateConfirmPwd(r.ConfirmPassword, r.Password),
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

		// activate user
		usr, err := controller.UserSvc.ActivateUser(rBody.EmployeeID, rBody.Email, rBody.Password)
		if err != nil {
			util.RespondWithError(w, err.(util.ResponseError))
			return
		}

		// return response
		util.RespondWithJSON(w, http.StatusOK, usr)
	}
}

// CreateUser creates a new user.
func CreateUser() *controller.ProtectedHandler {
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

	return &controller.ProtectedHandler{
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
			usr, err := controller.UserSvc.CreateUser(
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
func GetUser() *controller.ProtectedHandler {
	validator := func(id string) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id": validation.ValidateEmpID(id),
			},
		}
	}

	return &controller.ProtectedHandler{
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
			usr, err := controller.UserSvc.GetUserByEmpID(id)
			if err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}
			util.RespondWithJSON(w, http.StatusOK, usr)
		},
	}
}
