package controller

import (
	"fmt"
	"net/http"

	"adeia-api/internal/api/middleware"
	"adeia-api/internal/util"
	"adeia-api/internal/util/validation"

	"github.com/go-chi/chi"
)

// RoleRoutes returns a slice containing all role-related routes.
func RoleRoutes() (string, chi.Router) {
	r := chi.NewRouter()

	// only authenticated users can create/edit roles
	r.Use(middleware.AllowAuthenticated(sessionSvc, usrSvc))
	r.Method(http.MethodPost, "/", CreateRole())

	return "/roles", r
}

// CreateRole creates a new role.
func CreateRole() *ProtectedHandler {
	type request struct {
		Name string `json:"name"`
	}

	validator := func(r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"name": validation.ValidateResourceName(r.Name),
			},
		}
	}

	return &ProtectedHandler{
		PermissionName: "CREATE_ROLES",
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

			// create role
			role, err := roleSvc.CreateRole(rBody.Name)
			if err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			// return response
			w.Header().Set("Location", fmt.Sprintf("/v1/roles/%d", role.ID))
			util.RespondWithJSON(w, http.StatusCreated, role)
		},
	}
}
