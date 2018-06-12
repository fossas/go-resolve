// Package api provides mid-level primitives for writing JSON APIs for Ferox
// applications.
package api

import (
	"net/http"
)

// Handle transforms an API handler into an http.HandlerFunc.
func Handle(handler func(req *Request) (*Response, *Error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		res, err := handler(newRequest(r))
		if err != nil {
			render(w, err)
			return
		}
		render(w, res)
	}
}
