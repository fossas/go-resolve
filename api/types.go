package api

import "github.com/fossas/go-resolve/common"

type LookupRequest struct {
	Hash string
}

type LookupResponse struct {
	Ok     bool
	Err    string         `json:",omitempty"`
	Result common.Package `json:",omitempty"`
}

type VerifyRequest struct {
	common.Package
}

type VerifyResponse struct {
	Ok       bool
	Err      string         `json:",omitempty"`
	Expected common.Package `json:",omitempty"`
}
