package role

import (
	"fmt"
	"net/http"
	"strconv"

	"adeia/internal/api/middleware"
	"adeia/internal/controller"
	"adeia/internal/model"
	"adeia/internal/util"
	"adeia/internal/util/validation"

	"github.com/go-chi/chi"
)

// Routes returns a slice containing all role-related routes.
func Routes() (string, chi.Router) {
	r := chi.NewRouter()

	// only authenticated users can create/edit roles
	r.Use(middleware.AllowAuthenticated(controller.SessionSvc, controller.UserSvc))
	r.Method(http.MethodPost, "/", CreateRole())
	r.Method(http.MethodPut, "/{id}", UpdateRole())

	return "/roles", r
}

// CreateRole creates a new role.
func CreateRole() *controller.ProtectedHandler {
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

	return &controller.ProtectedHandler{
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
			role, err := controller.RoleSvc.CreateRole(rBody.Name)
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

// UpdateRole updates a role.
func UpdateRole() *controller.ProtectedHandler {
	type request struct {
		Name string `json:"name"`
	}

	validator := func(id string, r request) *validation.Validation {
		return &validation.Validation{
			Errors: validation.Errors{
				"id":   validation.ValidateResourceID(id),
				"name": validation.ValidateResourceName(r.Name),
			},
		}
	}

	return &controller.ProtectedHandler{
		PermissionName: "UPDATE_ROLES",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// decode request body
			var rBody request
			if err := util.DecodeBodyAndRespond(w, r, &rBody); err != nil {
				return
			}

			// validate request
			id := chi.URLParam(r, "id")
			if err := validator(id, rBody).Validate(); err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			roleID, _ := strconv.Atoi(id)
			role := model.Role{Name: rBody.Name}
			if err := controller.RoleSvc.UpdateByID(roleID, &role); err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			role.ID = roleID
			util.RespondWithJSON(w, http.StatusOK, role)
		},
	}
}
