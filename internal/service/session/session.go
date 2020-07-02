package session

import (
	"net/http"

	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/crypto"
	"adeia-api/internal/util/log"

	"github.com/go-ozzo/ozzo-validation/v4"
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
	GetAndRefresh(r *http.Request) (value string, err error)
	Destroy(w http.ResponseWriter, r *http.Request) error

	addCookie(w http.ResponseWriter, s *Session, maxAge int)
	newSessionID() (string, error)
	retrieveSession(sessionID string) (Session, error)
	readFromCookie(r *http.Request) (sessionID string, err error)
	storeSession(s Session) error
}

// Session represents the session object, containing a sessionID:value pair. The
// value currently is only a string, but in the future it may be a struct,
// containing multiple fields that can tie an user to a sessionID.
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

// New returns a new session service.
func New(store Store) Service {
	return &Impl{
		store:        store,
		cookieName:   constants.SessionCookieKey,
		maxCookieAge: constants.SessionExpiry,
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
	i.addCookie(w, &s, i.maxCookieAge)
	return nil
}

// GetAndRefresh gets the sessionID from the cookie, retrieves the corresponding value from
// the store.
func (i *Impl) GetAndRefresh(r *http.Request) (value string, err error) {
	sessionID, err := i.readFromCookie(r)
	if err != nil {
		return "", err
	}

	// validate sessionID
	if err := validation.Validate(sessionID,
		validation.Required,
	); err != nil {
		log.Debugf("validation failed for sessionID: %v", err)
		return "", err
	}

	session, err := i.retrieveSession(sessionID)
	if err != nil {
		return "", err
	}
	return session.value, nil
}

// Destroy destroys a session from the cookie and from the cache.
func (i *Impl) Destroy(w http.ResponseWriter, r *http.Request) error {
	sessionID, err := i.readFromCookie(r)
	if err != nil {
		return err
	}

	// remove from cache
	if err := i.store.ExpireSession(sessionID); err != nil {
		return err
	}

	i.addCookie(w, &Session{}, -1)
	return nil
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
	b, err := crypto.GenerateRandomBytes(128)
	if err != nil {
		return "", err
	}

	return crypto.EncodeBase64(b), nil
}

func (i *Impl) addCookie(w http.ResponseWriter, s *Session, maxAge int) {
	util.AddCookie(w, i.cookieName, s.sessionID, "/", maxAge)
}

func (i *Impl) readFromCookie(r *http.Request) (sessionID string, err error) {
	return util.GetCookie(r, constants.SessionCookieKey)
}
