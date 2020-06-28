package middleware

import (
	"context"
	"net/http"

	"adeia-api/internal/service/user"
	"adeia-api/internal/util"
	"adeia-api/internal/util/log"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Authenticated is a middleware that allows only authenticated users.
func Authenticated(usrSvc user.Service) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get session id
			sessID, err := util.GetSessionFromCookie(r)
			if err != nil {
				log.Debugf("cannot get session cookie: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// validate sessionID
			if err := validation.Validate(sessID,
				validation.Required,
				is.UUIDv4,
			); err != nil {
				log.Debugf("validation failed for sessID: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// get user
			usr, err := usrSvc.GetAuthenticated(sessID)
			if err != nil {
				log.Debugf("cannot get user for associated session: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// store in context
			ctx := context.WithValue(r.Context(), util.ContextUserKey, usr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
