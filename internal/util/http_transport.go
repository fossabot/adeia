package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"adeia/internal/util/constants"
	"adeia/internal/util/log"

	"github.com/golang/gddo/httputil/header"
)

type dataResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Error ResponseError `json:"error"`
}

// RespondWithError writes the provided ResponseError as a JSON response.
func RespondWithError(w http.ResponseWriter, err ResponseError) {
	o := &errorResponse{err}
	respond(w, err.StatusCode, o)
}

// RespondWithJSON writes the given payload as a JSON response.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	o := &dataResponse{payload}
	respond(w, statusCode, o)
}

func respond(w http.ResponseWriter, statusCode int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("failed to marshal JSON response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

// DecodeJSONBody decodes a JSON request.Body into the provided interface
// (usually a struct). Adapted from Alex Edwards's blog, which is released under
// the MIT license.
// See: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			// request body not JSON
			return ErrInvalidRequestBody
		}
	}

	// set max body size
	r.Body = http.MaxBytesReader(w, r.Body, constants.MaxReqBodySize)

	dec := json.NewDecoder(r.Body)
	// TODO: decide if we need to disallow unknown fields
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case
			errors.As(err, &syntaxError),
			errors.Is(err, io.ErrUnexpectedEOF):
			// badly formed JSON
			return ErrInvalidJSON

		case errors.As(err, &unmarshalTypeError):
			// invalid value for field
			return ErrValidationFailed.AddValidationErr(
				unmarshalTypeError.Field,
				fmt.Sprintf("Please enter a valid %v", unmarshalTypeError.Field),
			)

		// There is an open issue regarding turning this into a sentinel error
		// at https://github.com/golang/go/issues/29035.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			// unknown field
			field, _ := strconv.Unquote(strings.TrimPrefix(err.Error(), "json: unknown field "))
			return ErrUnknownField.Msgf("Unknown field: %v", field)

		case errors.Is(err, io.EOF):
			// request body empty
			return ErrInvalidJSON

		// There is an open issue regarding turning this into a sentinel error
		// at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			// request body too large
			return ErrRequestBodyTooLarge

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		// contains multiple JSON objects
		return ErrInvalidJSON
	}

	return nil
}

func isReqMalformedErr(err error) bool {
	switch err.Error() {
	case
		ErrInvalidJSON.Error(),
		ErrInvalidRequestBody.Error(),
		ErrValidationFailed.Error(),
		ErrRequestBodyTooLarge.Error(),
		ErrUnknownField.Error():
		return true
	}

	return false
}

// DecodeBodyAndRespond calls the DecodeJSONBody on the provided params. It
// responds with the appropriate error when the request body is malformed and returns
// the error for the caller (usually the controller) to stop further processing
// the request.
func DecodeBodyAndRespond(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if err := DecodeJSONBody(w, r, dst); err != nil {
		if isReqMalformedErr(err) {
			// request malformed
			log.Debugf("received malformed request body: %v", err)
			RespondWithError(w, err.(ResponseError))
			return err
		}

		// we cannot parse due to some other error
		log.Debugf("error parsing request body: %v", err)
		RespondWithError(w, ErrParseReqBodyFailed)
		return err
	}

	return nil
}

// AddCookie adds a new cookie.
func AddCookie(w http.ResponseWriter, name, value, path string, maxAge int) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		// TODO: make this secure when HTTPS is set-up
		// Secure:     false,
	}
	http.SetCookie(w, &cookie)
}

// GetCookie returns the value in the cookie identified by the name.
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
