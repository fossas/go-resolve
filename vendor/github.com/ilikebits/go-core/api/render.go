package api

import (
	"encoding/json"
	"net/http"
)

// renderable is a full HTTP JSON response.
type renderable interface {
	httpStatusCode() int

	// This will be marshalled into JSON.
	body() interface{}
}

// render writes a renderable to an http.ResponseWriter, handling details like
// response code, marshalling, and error handling.
func render(w http.ResponseWriter, v renderable) {
	res, err := json.Marshal(v.body())
	if err != nil {
		res, err := json.Marshal(ErrorInternal(err))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
	}
	w.WriteHeader(v.httpStatusCode())
	w.Write(res)
}
