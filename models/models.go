// Package models implements application-level models.
package models

import "fmt"

type Package struct {
	ImportPath string
	Hash       string

	SourceFiles []string
	Imports     []string
	Deps        []string

	VCS        VCS
	Repository string
	Revision   string
	Version    string
}

func (p *Package) String() string {
	return fmt.Sprintf("%#v", p)
}

//go:generate stringer -type=VCS

type VCS int

const (
	_ VCS = iota
	Git
	Subversion
	Mercurial
	Bazaar
)
