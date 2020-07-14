package session

import (
	"net/http"
	"time"

	"adeia-api/internal/db"
	"adeia-api/internal/model"
	"adeia-api/internal/repo/session"
	"adeia-api/internal/util"
	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/crypto"
	"adeia-api/internal/util/log"
)

// Service contains all session-related business logic.
type Service interface {
	NewSession(id int, empID string) (accessToken, refreshToken string, err error)
	RefreshToken(id int, empID, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	AddRefreshTokenCookie(w http.ResponseWriter, refreshToken string)
	ReadRefreshTokenCookie(r *http.Request) (string, error)
	ParseAccessToken(jwt string) (id string, err error)
	Destroy(id int, refreshToken string) error

	newAccessToken(empID string) (string, error)
	newRefreshToken() (token, hash []byte, err error)
}

// Impl is a Service implementation.
type Impl struct {
	sessionRepo session.Repo
}

// New returns a new Service.
func New(d db.DB) Service {
	return &Impl{session.New(d)}
}

// NewSession creates a new session and returns the accessToken and refreshToken.
func (i *Impl) NewSession(id int, empID string) (accessToken, refreshToken string, err error) {
	// generate new access token
	accessToken, err = i.newAccessToken(empID)
	if err != nil {
		log.Errorf("cannot create new JWT: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// generate new refresh token
	token, hash, err := i.newRefreshToken()
	if err != nil {
		log.Errorf("cannot create new refresh token: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// store refreshToken's hash to db
	s := model.Session{
		UserID:              id,
		RefreshToken:        hash,
		RefreshTokenExpires: time.Now().UTC().Add(time.Second * constants.RefreshTokenExpiry),
	}
	_, err = i.sessionRepo.Insert(&s)
	if err != nil {
		log.Errorf("cannot insert session into db: %v", err)
		return "", "", util.ErrDatabaseError
	}

	return accessToken, crypto.EncodeBase64(token), nil
}

// RefreshToken refreshes the tokens (generates new access and refresh tokens) and
// updates the corresponding Session.
func (i *Impl) RefreshToken(id int, empID, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// decode token to get byte array
	r, err := crypto.DecodeBase64(refreshToken)
	if err != nil {
		log.Errorf("cannot decode base64 refreshToken: %v", err)
		return "", "", util.ErrUnauthorized
	}

	// get user
	s, err := i.sessionRepo.GetByUserIDAndRefreshToken(id, crypto.Hash(r))
	if err != nil {
		log.Errorf("cannot get session by refreshToken: %v", err)
		return "", "", util.ErrDatabaseError
	} else if s == nil {
		log.Errorf("no session associated with refreshToken: %v", err)
		return "", "", util.ErrUnauthorized
	}

	// check expiry
	if s.RefreshTokenExpires.Before(time.Now().UTC()) {
		// refresh token has expired, redirect to login
		return "", "", util.ErrMustLogin
	}

	// generate new refresh token
	token, hash, err := i.newRefreshToken()
	if err != nil {
		log.Errorf("cannot create new refresh token: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// update session
	if _, err := i.sessionRepo.UpdateRefreshToken(
		s.ID,
		hash,
		time.Now().UTC().Add(time.Second*constants.RefreshTokenExpiry),
	); err != nil {
		log.Errorf("cannot update session: %v", err)
		return "", "", util.ErrDatabaseError
	}

	// create new jwt
	newAccessToken, err = i.newAccessToken(empID)
	if err != nil {
		log.Errorf("cannot create new JWT: %v", err)
		return "", "", util.ErrInternalServerError
	}

	return newAccessToken, crypto.EncodeBase64(token), nil
}

// Destroy deletes a Session identified by the id and refresh token.
func (i *Impl) Destroy(id int, refreshToken string) error {
	// decode token to get byte array
	r, err := crypto.DecodeBase64(refreshToken)
	if err != nil {
		log.Errorf("cannot decode base64 refreshToken: %v", err)
		return util.ErrInternalServerError
	}

	// delete session
	rowsAffected, err := i.sessionRepo.DeleteByUserIDAndRefreshToken(id, crypto.Hash(r))
	if err != nil {
		return util.ErrBadRequest
	} else if rowsAffected == 0 {
		return util.ErrResourceNotFound
	}

	return nil
}

// ParseAccessToken parses the access token and returns the empID.
func (i *Impl) ParseAccessToken(t string) (id string, err error) {
	// get payload
	payload, err := crypto.ParseJWT(t)
	if err != nil {
		log.Infof("invalid jwt token: %v", err)
		return "", util.ErrUnauthorized
	}

	// get id from payload
	id, ok := payload["id"].(string)
	if !ok {
		log.Infof("invalid id in payload")
		return "", util.ErrUnauthorized
	}
	return id, nil
}

// AddRefreshTokenCookie adds the refreshToken to a new cookie.
func (i *Impl) AddRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	util.AddCookie(
		w,
		constants.RefreshTokenCookieName,
		refreshToken,
		"/v1/users/sessions",
		constants.RefreshTokenExpiry,
	)
}

// ReadRefreshTokenCookie reads the refreshToken from the cookie.
func (i *Impl) ReadRefreshTokenCookie(r *http.Request) (string, error) {
	return util.GetCookie(r, constants.RefreshTokenCookieName)
}

func (i *Impl) newAccessToken(empID string) (string, error) {
	payload := map[string]interface{}{"id": empID}
	return crypto.NewJWT(payload, time.Second*constants.AccessTokenExpiry)
}

func (i *Impl) newRefreshToken() (token, hash []byte, err error) {
	b, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, err
	}
	return b, crypto.Hash(b), nil
}
