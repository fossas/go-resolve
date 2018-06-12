package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Request wraps a raw request with useful API utilities.
type Request struct {
	Raw *http.Request
}

// Validatable is implemented by requests that should be validated.
type Validatable interface {
	Validate() error
}

// TODO: other ideas for request-scoped utilities:
//
//   - Database transaction to be automatically committed or aborted.
//   - OpenTracing using request context to create and finish spans.
//

// newRequest sets up a Request.
func newRequest(raw *http.Request) *Request {
	// Set up request-scoped logger.
	ctx := raw.Context()
	requestID := middleware.GetReqID(ctx)
	logger := log.With().Str("RequestID", requestID).Logger()

	return &Request{
		Raw: raw.WithContext(logger.WithContext(ctx)),
	}
}

// JSON unmarshals the request body into JSON and validates the result,
// converting errors into API errors.
func (r *Request) JSON(v interface{}) *Error {
	req := r.Raw
	defer req.Body.Close()

	// Read request body.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ErrorInternal(err)
	}

	// Decode request body as JSON into struct.
	err = json.Unmarshal(body, v)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalFieldError:
			return ErrorInvalidArgs(err)
		case *json.UnmarshalTypeError:
			return ErrorInvalidArgs(err)
		default:
			return ErrorMalformedJSON(err)
		}
	}

	// Validate request struct.
	switch v := v.(type) {
	case Validatable:
		err := v.Validate()
		if err != nil {
			return ErrorInvalidArgs(err)
		}
	default:
	}

	return nil
}

// TODO: use go:generate or something to parse and verify IO annotations?
// Example annotations:
//
//   - @IO(Database)
//   - @IO(Database:slow)
//   - @IO(HTTP)
//   - @IO(ResponseWriter)
//   - @IO(Database, HTTP, *)
//
// We can annotate some base functions (e.g. something that writes to a
// http.ResponseWriter) and then check to make sure none of our code uses raw
// http.ResponseWriters without explicit exception annotations.

// Context returns the underlying request's context.Context.
func (r *Request) Context() context.Context {
	return r.Raw.Context()
}

// Logger returns the request's scoped logger. It adds a RequestID key to the
// structured logging fields.
func (r *Request) Logger() *zerolog.Logger {
	return zerolog.Ctx(r.Context())
}
