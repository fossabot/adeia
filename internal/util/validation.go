package util

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

// Validation represents a map of validation rules.
type Validation struct {
	validation.Errors
}

// Validate calls the underlying library's Filter() and returns an
// ErrValidationFailed when validation fails.
func (v *Validation) Validate() error {
	err := v.Errors.Filter()
	if err == nil {
		return nil
	}

	e := make(map[string]string)
	for k, v := range err.(validation.Errors) {
		e[k] = v.Error()
	}
	return ErrValidationFailed.ValidationErr(e)
}
