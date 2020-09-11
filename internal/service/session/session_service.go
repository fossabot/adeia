package session

import (
	"net/http"
	"time"

	"adeia/internal/model"
	"adeia/internal/service"
	"adeia/internal/util"
	"adeia/internal/util/constants"
	"adeia/internal/util/crypto"
	"adeia/internal/util/log"
)

type Service struct {
	sessionRepo service.SessionRepo
	jwtSecret   string
}

// New returns a new Service.
func New(s service.SessionRepo, jwtSecret string) *Service {
	return &Service{
		sessionRepo: s,
		jwtSecret:   jwtSecret,
	}
}

// NewSession creates a new session and returns the accessToken and refreshToken.
func (s *Service) NewSession(id int, empID string) (accessToken, refreshToken string, err error) {
	// generate new access token
	accessToken, err = s.newAccessToken(empID)
	if err != nil {
		log.Errorf("cannot create new JWT: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// generate new refresh token
	token, hash, err := s.newRefreshToken()
	if err != nil {
		log.Errorf("cannot create new refresh token: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// store refreshToken's hash to db
	session := model.Session{
		UserID:              id,
		RefreshToken:        hash,
		RefreshTokenExpires: time.Now().UTC().Add(time.Second * constants.RefreshTokenExpiry),
	}
	_, err = s.sessionRepo.Insert(&session)
	if err != nil {
		log.Errorf("cannot insert session into db: %v", err)
		return "", "", util.ErrDatabaseError
	}

	return accessToken, crypto.EncodeBase64(token), nil
}

// RefreshToken refreshes the tokens (generates new access and refresh tokens) and
// updates the corresponding Session.
func (s *Service) RefreshToken(id int, empID, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// decode token to get byte array
	r, err := crypto.DecodeBase64(refreshToken)
	if err != nil {
		log.Errorf("cannot decode base64 refreshToken: %v", err)
		return "", "", util.ErrUnauthorized
	}

	// get user
	session, err := s.sessionRepo.GetByUserIDAndRefreshToken(id, crypto.Hash(r))
	if err != nil {
		log.Errorf("cannot get session by refreshToken: %v", err)
		return "", "", util.ErrDatabaseError
	} else if session == nil {
		log.Errorf("no session associated with refreshToken: %v", err)
		return "", "", util.ErrUnauthorized
	}

	// check expiry
	if session.RefreshTokenExpires.Before(time.Now().UTC()) {
		// refresh token has expired, redirect to login
		return "", "", util.ErrMustLogin
	}

	// generate new refresh token
	token, hash, err := s.newRefreshToken()
	if err != nil {
		log.Errorf("cannot create new refresh token: %v", err)
		return "", "", util.ErrInternalServerError
	}

	// update session
	if _, err := s.sessionRepo.UpdateRefreshToken(
		session.ID,
		hash,
		time.Now().UTC().Add(time.Second*constants.RefreshTokenExpiry),
	); err != nil {
		log.Errorf("cannot update session: %v", err)
		return "", "", util.ErrDatabaseError
	}

	// create new jwt
	newAccessToken, err = s.newAccessToken(empID)
	if err != nil {
		log.Errorf("cannot create new JWT: %v", err)
		return "", "", util.ErrInternalServerError
	}

	return newAccessToken, crypto.EncodeBase64(token), nil
}

// Destroy deletes a Session identified by the id and refresh token.
func (s *Service) Destroy(id int, refreshToken string) error {
	// decode token to get byte array
	r, err := crypto.DecodeBase64(refreshToken)
	if err != nil {
		log.Errorf("cannot decode base64 refreshToken: %v", err)
		return util.ErrInternalServerError
	}

	// delete session
	rowsAffected, err := s.sessionRepo.DeleteByUserIDAndRefreshToken(id, crypto.Hash(r))
	if err != nil {
		return util.ErrBadRequest
	} else if rowsAffected == 0 {
		return util.ErrResourceNotFound
	}

	return nil
}

// ParseAccessToken parses the access token and returns the empID.
func (s *Service) ParseAccessToken(t string) (id string, err error) {
	// get payload
	payload, err := crypto.ParseJWT(s.jwtSecret, t)
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
func (s *Service) AddRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	util.AddCookie(
		w,
		constants.RefreshTokenCookieName,
		refreshToken,
		"/v1/users/sessions",
		constants.RefreshTokenExpiry,
	)
}

// ReadRefreshTokenCookie reads the refreshToken from the cookie.
func (s *Service) ReadRefreshTokenCookie(r *http.Request) (string, error) {
	return util.GetCookie(r, constants.RefreshTokenCookieName)
}

func (s *Service) newAccessToken(empID string) (string, error) {
	payload := map[string]interface{}{"id": empID}
	return crypto.NewJWT(s.jwtSecret, payload, time.Second*constants.AccessTokenExpiry)
}

func (s *Service) newRefreshToken() (token, hash []byte, err error) {
	b, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, err
	}
	return b, crypto.Hash(b), nil
}
