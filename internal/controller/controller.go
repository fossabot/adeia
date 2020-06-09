package controller

import (
	"adeia-api/internal/model"
	"adeia-api/internal/repository"
	log "adeia-api/internal/utils/logger"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index is a simple handler that writes a welcome message.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome\n")
}

// SampleRoute is a handler to test db access.
func SampleRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	repo := repository.NewUserRepo()

	user := model.NewUser("test", "example@example1.com")

	err := repo.Insert(user)
	if err != nil {
		log.Debugf("error inserting user: %v", err)
	}
	log.Debug("successfully inserted")

	s, err := repo.GetByEmpID(user.EmployeeID)
	if err != nil {
		log.Debugf("error selecting user: %v", err)
	}
	log.Debugf("retrieved user: %+v", s)
}
