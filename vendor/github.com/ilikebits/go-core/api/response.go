package api

import "net/http"

// Response wraps a successful API response.
type Response struct {
	StatusCode int         `json:"-"`
	Result     interface{} `json:"result"`
}

func (r *Response) HTTPStatusCode() int {
	return r.StatusCode
}

func (r *Response) Body() interface{} {
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
		StatusCode: http.StatusOK,
		Result:     result,
	}
}
