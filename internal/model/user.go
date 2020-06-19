package model

import "github.com/dchest/uniuri"

// declarations to generate employee ID
var employeeIDChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const employeeIDLength = 6

// User represents the User model.
type User struct {
	// Id is a surrogate primary key that is auto-incremented by the database.
	// This has no meaning outside the database, except that it only identifies an account.
	// This will never change for an account and all foreign keys must use this field.
	// This should not be exposed outside.
	ID int `db:"id"`

	// EmployeeId is a natural key used extensively throughout the system (in URIs, etc.).
	// It must be unique, short and user-rememberable (preferably 6-8 chars long).
	// It can change for an account, if needed.
	EmployeeID string `db:"employee_id"`

	// Name represents the name of the user.
	Name string `db:"name"`

	// Email represents the name of the user.
	Email string `db:"email"`

	// Password represents the hashed password of the user.
	Password string `db:"password"`
}

// New creates a new User.
func New(name, email string) *User {
	return &User{
		EmployeeID: newEmpID(),
		Name:       name,
		Email:      email,
	}
}

// newEmpID generates a user and URL-friendly alpha-numeric (upper-case letters
// and numbers) employee ID.
func newEmpID() string {
	return uniuri.NewLenChars(employeeIDLength, employeeIDChars)
}
