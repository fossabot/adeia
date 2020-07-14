package model

import "time"

// Session represents the Session model.
//
// A new Session is created for every user during every login and it is deleted when the user
// logs out. Every user can have multiple sessions at a time (each session represents a
// different device). After successful login, we return an access token (a JWT) with a
// short expiry and a refresh token (random bytes; set as a cookie) with a longer expiry.
//
// The access token will be sent with all requests as a bearer token in the `Authorization`
// header. Before the access token expires, the client must refresh its access token by
// hitting the refresh route (along with refresh token in the cookie). During "refresh",
// we generate new access and refresh tokens, and reset their expiry.
type Session struct {
	// ID is a surrogate key without any meaning. Table is indexed on this field and
	// is an auto-increment. This field should not be exposed outside.
	ID int `db:"id"`

	// UserID is a foreign-key referencing the `id` field in model.User. Even this
	// field is a surrogate key and should not be exposed outside for security purposes.
	UserID int `db:"user_id"`

	// RefreshToken is a hashed version of the refresh token. It is a unique field and
	// is stored as a byte array.
	RefreshToken []byte `db:"refresh_token"`

	// RefreshTokenExpires represents the time of expiry of the refresh token that the
	// session is associated with. If the refresh token expires, "refreshing" the
	// access token will fail, and the user must login again. As all other time fields
	// in the database, this also must be stored in UTC, and read to local time by the
	// application.
	RefreshTokenExpires time.Time `db:"refresh_token_expires"`
}
