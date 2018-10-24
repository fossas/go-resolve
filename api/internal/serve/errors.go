package api

import (
	"fmt"
	"net/http"
)

// Error wraps an error API response.
type Error struct {
	Raw        error `json:"-"` // For comparing against exported Go error values.
	StatusCode int   `json:"-"`

	ErrorCode string `json:"error"`
	Message   string `json:"message,omitempty"`
}

func (e Error) HTTPStatusCode() int {
	return e.StatusCode
}

func (e Error) Body() interface{} {
	return renderedError{
		Ok:        false,
		ErrorCode: e.ErrorCode,
		Message:   e.Message,
	}
}

type renderedError struct {
	Ok        bool   `json:"ok"`
	ErrorCode string `json:"error"`
	Message   string `json:"message,omitempty"`
}

// ErrorInternal wraps an unexpected internal error (something akin to a failed
// assertion or panic).
func ErrorInternal(err error) *Error {
	return &Error{
		Raw:        err,
		StatusCode: http.StatusInternalServerError,

		ErrorCode: "INTERNAL_ERROR",
		Message:   fmt.Sprintf("an unexpected error occurred (%s)", err.Error()),
	}
}

// ErrorMalformedJSON wraps JSON unmarshalling errors.
func ErrorMalformedJSON(err error) *Error {
	return &Error{
		Raw:        err,
		StatusCode: http.StatusBadRequest,

		ErrorCode: "MALFORMED_JSON",
		Message:   fmt.Sprintf("could not parse the request as JSON (%s)", err.Error()),
	}
}

// ErrorInvalidArgs wraps invalid argument errors (e.g. JSON validation errors).
func ErrorInvalidArgs(err error) *Error {
	return &Error{
		Raw:        err,
		StatusCode: http.StatusBadRequest,

		ErrorCode: "INVALID_ARGUMENT",
		Message:   fmt.Sprintf("invalid argument (%s)", err.Error()),
	}
}
