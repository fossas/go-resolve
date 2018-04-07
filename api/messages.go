package api

import "github.com/fossas/go-resolve/resolve"

type LookupRequest struct {
	Hash string
}

type LookupResponse struct {
	Ok     bool
	Err    string          `json:",omitempty"`
	Result resolve.Package `json:",omitempty"`
}

type VerifyRequest struct {
	resolve.Package
}

type VerifyResponse struct {
	Ok       bool
	Err      string          `json:",omitempty"`
	Expected resolve.Package `json:",omitempty"`
}

type CrawlRequest struct {
	Name     string
	Revision string
}

type CrawlResponse struct {
	Ok  bool
	Err string `json:",omitempty"`
}
