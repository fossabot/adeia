package model

import "database/sql"

// User represents the User model.
type User struct {
	// Id is a surrogate primary key that is auto-incremented by the database.
	// This has no meaning outside the database, except that it only identifies an account.
	// This will never change for an account and all foreign keys must use this field.
	// This should not be exposed outside.
	ID int `db:"id" json:"-"`

	// EmployeeId is a natural key used extensively throughout the system (in URIs, etc.).
	// It must be unique, short and user-rememberable (preferably 6-8 chars long).
	// It is case-insensitive (internally managed by Postgres as `citext` type).
	// See: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#user-ids
	EmployeeID string `db:"employee_id" json:"employee_id"`

	// Name represents the name of the user.
	Name string `db:"name" json:"name"`

	// Email represents the name of the user.
	Email string `db:"email" json:"email"`

	// Password represents the hashed password of the user.
	Password string `db:"password" json:"-"`

	// Designation represents the designation of the user.
	Designation string `db:"designation" json:"designation"`

	// IsActivated represents whether the user account is activated or not.
	IsActivated bool `db:"is_activated" json:"is_activated"`

	// DeletedAt stores the time at which a user was deleted. A timestamp is used
	// instead of a bool, because we may need to run cleanups like "clear db records
	// of users that were deleted more than a month ago". This is a nullable field
	// and all queries should use `WHERE deleted_at IS NULL`.
	DeletedAt sql.NullTime `db:"deleted_at" json:"-"`
}
