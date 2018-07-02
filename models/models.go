// Package models implements common application-level models.
package models

import "time"

type Package struct {
	ImportPath  string    `db:"import_path"`
	Revision    string    `db:"revision"`
	Hash        string    `db:"hash"`
	Version     string    `db:"version"`
	LastUpdated time.Time `db:"last_updated"`
}

func (p *Package) String() string {
	return p.ImportPath + " " + p.Revision + " " + p.Hash
}
