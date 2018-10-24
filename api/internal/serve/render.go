package api

import (
	"encoding/json"
	"net/http"
)

// Renderable is a full HTTP JSON response.
type Renderable interface {
	HTTPStatusCode() int
	Body() interface{} // This will be marshalled into JSON.
}

// render writes a renderable to an http.ResponseWriter, handling details like
// response code, marshalling, and error handling.
func render(w http.ResponseWriter, v Renderable) {
	res, err := json.Marshal(v.Body())
	if err != nil {
		res, err := json.Marshal(ErrorInternal(err))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
	}
	w.WriteHeader(v.HTTPStatusCode())
	w.Write(res)
}
