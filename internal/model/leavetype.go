package model

// Holiday represents the Holiday model.
type LeaveType struct {
	// ID represents the leave type ID. It is an auto-incremented db field.
	ID int `db:"id" json:"id"`

	// Name represents the type of the leave.
	Name string `db:"name" json:"name"`

	// Days represents the number of days of the holiday.
	Days int `db:"type" json:"type"`
}
