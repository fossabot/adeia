package util

import (
	"fmt"
	"net/http"
)

// ResponseError represents an error that is sent as a response to the client.
type ResponseError struct {
	StatusCode       int               `json:"-"`
	ErrorCode        string            `json:"code"`
	Message          string            `json:"message,omitempty"`
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
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

// AddValidationErr adds a new ValidationError to the ValidationErrors map.
func (re ResponseError) AddValidationErr(f, m string) ResponseError {
	if re.ValidationErrors == nil {
		re.ValidationErrors = make(map[string]string)
	}
	re.ValidationErrors[f] = m
	return re
}

// ValidationErr sets the entire ValidationErrors map to the provided map, which
// is useful when processing fields using third-party validation libraries.
func (re ResponseError) ValidationErr(m map[string]string) ResponseError {
	re.ValidationErrors = m
	return re
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
		Message:    "A resource already exists with the specified fields",
	}

	// ErrInternalServerError is a generic error returned when something fails
	// we don't want to send too much information to the client. Add a message to
	// this to specify additional details.
	ErrInternalServerError = ResponseError{
		StatusCode: http.StatusInternalServerError,
		ErrorCode:  "INTERNAL_SERVER_ERROR",
	}

	// ErrResourceNotFound is the error returned when a resource is not found.
	// This is usually returned from GET, PUT, DELETE routes.
	ErrResourceNotFound = ResponseError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "RESOURCE_NOT_FOUND",
	}

	// ErrBadRequest is a generic error returned when a request is bad.
	ErrBadRequest = ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "BAD_REQUEST",
	}

	// ErrUnauthorized is the error returned when the user is not authenticated/
	// authorized to perform an action.
	ErrUnauthorized = ResponseError{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  "UNAUTHORIZED",
		Message:    "You need to be logged-in to perform this action.",
	}

	// ErrAccountNotActivated is the error returned when the user account is not
	// activated.
	ErrAccountNotActivated = ResponseError{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  "ACCOUNT_NOT_ACTIVATED",
		Message:    "Your account is not activated",
	}
)
