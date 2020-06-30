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

// Authenticate is a middleware that allows only authenticated users through.
func Authenticate(sessionSvc session.Service, usrSvc user.Service) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user id
			userID, err := sessionSvc.Get(r)
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
