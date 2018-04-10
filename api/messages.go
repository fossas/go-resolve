package api

import "github.com/fossas/go-resolve/resolve"

// LookupRequest searches for a package given a hash.
type LookupRequest struct {
	// These fields should probably be pointers so we can tell when they're
	// missing and return an `INVALID_REQUEST` error.
	Hash string
}

// LookupResponse is an Either Err Result.
type LookupResponse struct {
	Ok     bool
	Err    string          `json:",omitempty"`
	Result resolve.Package `json:",omitempty"`
}

// VerifyRequest checks whether a hash is correct for a given package and
// revision.
type VerifyRequest struct {
	resolve.Package
}

// VerifyResponse is an Either Err Expected.
type VerifyResponse struct {
	Ok       bool
	Err      string          `json:",omitempty"`
	Expected resolve.Package `json:",omitempty"`
}

// CrawlRequest submits a package resolution job.
type CrawlRequest struct {
	resolve.Key
}

// CrawlResponse is an Either Err ().
type CrawlResponse struct {
	Ok  bool
	Err string `json:",omitempty"`
}

// ErrorResponse is a struct for holding errors.
type ErrorResponse struct {
	Ok  bool
	Err string `json:",omitempty"`
}
