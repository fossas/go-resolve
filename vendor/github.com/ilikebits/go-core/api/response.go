package api

import "net/http"

// Response wraps a successful API response.
type Response struct {
	HTTPStatusCode int         `json:"-"`
	Result         interface{} `json:"result"`
}

// Response implements renderable.
func (r *Response) httpStatusCode() int {
	return r.HTTPStatusCode
}

func (r *Response) body() interface{} {
	return renderedResponse{
		Ok:     true,
		Result: r.Result,
	}
}

type renderedResponse struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result,omitempty"`
}

// OK is a helper for creating Responses with a 200 status code.
func OK(result interface{}) *Response {
	return &Response{
		HTTPStatusCode: http.StatusOK,
		Result:         result,
	}
}
