package util

import (
	"fmt"
	"net/http"
)

const (
	// ValidationMissingField is the ValidationCode for a missing field.
	ValidationMissingField = "missing_field"

	// ValidationInvalid is the ValidationCode for an invalid field.
	ValidationInvalid = "invalid"

	// ValidationAlreadyExists is the ValidationCode for a field that already exists.
	ValidationAlreadyExists = "already_exists"

	// ValidationUnprocessable is the ValidationCode for a field that is unprocessable.
	ValidationUnprocessable = "unprocessable"

	// ValidationCustom is a custom ValidationCode.
	ValidationCustom = "custom"
)

// ValidationError represents a validation error on a field.
type ValidationError struct {
	Field          string `json:"field"`
	ValidationCode string `json:"code"`
	Message        string `json:"message"`
}

// ResponseError represents an error that is sent as a response to the client.
type ResponseError struct {
	StatusCode       int               `json:"-"`
	ErrorCode        string            `json:"code"`
	Message          string            `json:"message,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// Msg sets ResponseError's Message.
func (re ResponseError) Msg(m string) ResponseError {
	re.Message = m
	return re
}

// Msgf sets ResponseError's Message.
func (re ResponseError) Msgf(format string, a ...interface{}) ResponseError {
	re.Message = fmt.Sprintf(format, a...)
	return re
}

// Error returns the ErrorCode of ResponseError.
func (re ResponseError) Error() string {
	return re.ErrorCode
}

// ValidationErr appends a new ValidationError to the ResponseError.
func (re ResponseError) ValidationErr(f, v, m string) ResponseError {
	re.ValidationErrors = append(re.ValidationErrors, ValidationError{
		Field:          f,
		ValidationCode: v,
		Message:        m,
	})
	return re
}

// outputError is a small wrapper used when marshalling to JSON.
type outputError struct {
	Error ResponseError `json:"error"`
}

var (
	// ErrInvalidJSON is the error returned when the request body contains invalid
	// JSON.
	ErrInvalidJSON = ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "INVALID_JSON",
		Message:    "Problems parsing JSON",
	}

	// ErrInvalidRequestBody is the error returned when the request body is not
	// JSON.
	ErrInvalidRequestBody = ResponseError{
		StatusCode: http.StatusUnsupportedMediaType,
		ErrorCode:  "INVALID_REQUEST_BODY",
		Message:    "Body must be a JSON object",
	}

	// ErrValidationFailed is the error returned when some of the fields do not
	// conform to the validation rules. Add a list of ValidationErrors to this to
	// specify the fields that are failing.
	ErrValidationFailed = ResponseError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorCode:  "VALIDATION_FAILED",
		Message:    "Validation failed for some fields",
	}

	// ErrRequestBodyTooLarge is the error returned when the request body is too large.
	ErrRequestBodyTooLarge = ResponseError{
		StatusCode: http.StatusRequestEntityTooLarge,
		ErrorCode:  "REQUEST_BODY_TOO_LARGE",
		Message:    "Request body exceeds 1MiB",
	}

	// ErrParseReqBodyFailed is the error returned when the request body is okay,
	// but something else happened and we cannot parse it.
	ErrParseReqBodyFailed = ResponseError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "PARSE_REQUEST_BODY_FAILED",
		Message:    "An error occurred while parsing the request body",
	}

	// ErrUnknownField is the error returned when the request body contains an
	// unknown field.
	ErrUnknownField = ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "UNKNOWN_FIELD",
		Message:    "An unknown field is present in the request body",
	}

	// ErrResourceAlreadyExists is the error returned when a resource already
	// exists with the specified fields. This is usually returned from POST routes.
	ErrResourceAlreadyExists = ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "RESOURCE_ALREADY_EXISTS",
		Message:    "A resource already exists with the specified fields.",
	}

	// ErrInternalServerError is a generic error returned when something fails
	// we don't want to send too much information to the client. Add a message to
	// this to specify additional details.
	ErrInternalServerError = ResponseError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "INTERNAL_SERVER_ERROR",
	}

	ErrResourceNotFound = ResponseError{
		StatusCode: http.StatusNotFound,
		ErrorCode: "RESOURCE_NOT_FOUND",
	}
)
