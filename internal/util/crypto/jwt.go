package crypto

import (
	"time"

	"github.com/pascaldekloe/jwt"
	config "github.com/spf13/viper"
)

// NewJWT creates a jwt with the provided id and expiry.
func NewJWT(payload map[string]interface{}, expires time.Duration) (string, error) {
	var claims jwt.Claims

	// set claims & payload
	claims.Issued = jwt.NewNumericTime(time.Now().Round(time.Second))
	claims.Expires = jwt.NewNumericTime(time.Now().Add(expires).Round(time.Second))
	claims.Set = payload

	// sign token
	token, err := claims.HMACSign(jwt.HS256, []byte(config.GetString("server.jwt_secret")))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// ParseJWT validates a jwt and returns the payload.
func ParseJWT(token string) (payload map[string]interface{}, err error) {
	// check if signature is valid
	claims, err := jwt.HMACCheck([]byte(token), []byte(config.GetString("server.jwt_secret")))
	if err != nil {
		return nil, err
	}

	// check expiry
	if !claims.Valid(time.Now()) {
		return nil, err
	}

	// return payload
	return claims.Set, nil
}
