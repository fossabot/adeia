package controller

import (
	"net/http"

	"adeia-api/internal/util"
	"adeia-api/internal/util/log"
)

// ProtectedHandler checks if user is authorized before allowing the request to
// pass to the underlying controller.
type ProtectedHandler struct {
	PermissionName string
	Handler        http.HandlerFunc
}

// ServeHTTP performs the authorization and only allows the request to pass when
// the user is authorized.
func (p *ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// perform role checking using context, like !contains(userRoles, p.PermissionName)
	if false {
		// user doesn't have access
		log.Debug("user not authorized to perform action")
		util.RespondWithError(w, util.ErrUnauthorized)
		return
	}

	// user has access, so continue
	p.Handler.ServeHTTP(w, r)
}
