package middleware

import (
	"context"
	"net/http"
	"strconv"

	"adeia-api/internal/service/session"
	"adeia-api/internal/service/user"
	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/log"
)

// AllowAuthenticated is a middleware that allows users based on their auth state.
// If `allowAuthenticated` is true, it will only allow authenticated users through and if
// `allowAuthenticated` is set to false, it will only allow unauthenticated users through.
func AllowAuthenticated(sessionSvc session.Service, usrSvc user.Service, allowAuthenticated bool) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user id
			userID, err := sessionSvc.GetAndRefresh(r)
			if !allowAuthenticated {
				if err == nil {
					// user is authenticated
					log.Debugf("user is authenticated: %v", err)
					util.RespondWithError(w, util.ErrBadRequest.Msg("already authenticated"))
					return
				}

				// user is not authenticated
				next.ServeHTTP(w, r)
				return
			}

			if err != nil {
				// user is not authenticated
				log.Debugf("cannot get session cookie: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// get user
			userIDStr, _ := strconv.Atoi(userID)
			usr, err := usrSvc.GetUserByID(userIDStr)
			if err != nil {
				log.Debugf("cannot get user for associated session: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// store in context
			ctx := context.WithValue(r.Context(), constants.ContextUserKey, usr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
