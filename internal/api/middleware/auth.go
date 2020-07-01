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

// IsAuthenticated is a middleware that allows only authenticated users through.
func IsAuthenticated(sessionSvc session.Service, usrSvc user.Service) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user id
			userID, err := sessionSvc.GetAndRefresh(r)
			if err != nil {
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

// IsNotAuthenticated is a middleware that allows only users that are not authenticated, through.
func IsNotAuthenticated(sessionSvc session.Service, usrSvc user.Service) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userID, err := sessionSvc.GetAndRefresh(r); err == nil || userID == "" {
				log.Debugf("user is authenticated: %v", err)
				util.RespondWithError(w, util.ErrBadRequest.Msg("already authenticated"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
