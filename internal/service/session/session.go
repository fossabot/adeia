package session

import (
	"net/http"

	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/crypto"
	"adeia-api/internal/util/log"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Store represents the interface for the session store.
type Store interface {
	ExpireSession(key string) error
	GetAndRefreshSession(key string, expiry int) (string, error)
	SetSession(sessionID, value string, expiry int) error
}

// Service is the interface containing all methods for the Session service.
type Service interface {
	Create(w http.ResponseWriter, value string) error
	Get(r *http.Request) (value string, err error)

	addCookie(w http.ResponseWriter, s *Session)
	newSessionID() (string, error)
	retrieveSession(sessionID string) (Session, error)
	readFromCookie(r *http.Request) (sessionID string, err error)
	storeSession(s Session) error
}

// Session represents the session object, containing a sessionID:value pair. The
// value currently is only a string, but in the future it may be a struct,
// containing multiple fields that can tie a user to a sessionID.
type Session struct {
	sessionID string
	value     string
}

// Impl implements Service.
type Impl struct {
	store        Store
	cookieName   string
	maxCookieAge int
}

// NewService returns a new session service.
func NewService(store Store, cookieName string, cookieAge int) Service {
	return &Impl{
		store:        store,
		cookieName:   cookieName,
		maxCookieAge: cookieAge,
	}
}

// Create creates a new session object, adds to the store and adds the sessionID
// as a cookie.
func (i *Impl) Create(w http.ResponseWriter, value string) error {
	// create sessionID
	sessionID, err := i.newSessionID()
	if err != nil {
		return err
	}

	// add session to store
	s := Session{
		sessionID: sessionID,
		value:     value,
	}
	if err := i.storeSession(s); err != nil {
		return err
	}

	// add cookie
	i.addCookie(w, &s)
	return nil
}

// Get gets the sessionID from the cookie, retrieves the corresponding value from
// the store.
func (i *Impl) Get(r *http.Request) (value string, err error) {
	sessionID, err := i.readFromCookie(r)
	if err != nil {
		return "", err
	}

	// validate sessionID
	if err := validation.Validate(sessionID,
		validation.Required,
		is.UUIDv4,
	); err != nil {
		log.Debugf("validation failed for sessionID: %v", err)
		return "", nil
	}

	session, err := i.retrieveSession(sessionID)
	if err != nil {
		return "", err
	}
	return session.value, nil
}

func (i *Impl) retrieveSession(sessionID string) (Session, error) {
	value, err := i.store.GetAndRefreshSession(sessionID, i.maxCookieAge)
	if err != nil {
		return Session{}, err
	}

	return Session{
		sessionID: sessionID,
		value:     value,
	}, nil
}

func (i *Impl) storeSession(s Session) error {
	return i.store.SetSession(s.sessionID, s.value, i.maxCookieAge)
}

func (i *Impl) newSessionID() (string, error) {
	return crypto.NewUUID()
}

func (i *Impl) addCookie(w http.ResponseWriter, s *Session) {
	util.AddCookie(w, i.cookieName, s.sessionID, "/", i.maxCookieAge)
}

func (i *Impl) readFromCookie(r *http.Request) (sessionID string, err error) {
	return util.GetCookie(r, constants.SessionCookieKey)
}
