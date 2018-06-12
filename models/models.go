// Package models implements common application-level models.
package models

type Package struct {
	ImportPath string
	Revision   string
	Hash       string
	Version    string
}
