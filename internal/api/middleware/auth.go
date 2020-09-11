package middleware

import (
	"context"
	"net/http"
	"strings"

	"adeia/internal/controller"
	"adeia/internal/util"
	"adeia/internal/util/constants"
	"adeia/internal/util/log"
)

// AllowAuthenticated is a middleware that allows only authenticated users.
func AllowAuthenticated(sessionSvc controller.SessionService, userSvc controller.UserService) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get jwt from header
			t := r.Header.Get("Authorization")
			if t == "" {
				log.Debug("no authorization header present")
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// validate jwt
			token := strings.TrimPrefix(t, "Bearer ")
			empID, err := sessionSvc.ParseAccessToken(token)
			if err != nil {
				util.RespondWithError(w, err.(util.ResponseError))
				return
			}

			// get user
			usr, err := userSvc.GetUserByEmpID(empID)
			if err != nil {
				log.Debugf("cannot get user for associated session: %v", err)
				util.RespondWithError(w, util.ErrUnauthorized)
				return
			}

			// store user in context
			ctx := context.WithValue(r.Context(), constants.ContextUserKey, usr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
