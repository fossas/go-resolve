package api

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/fossas/go-resolve/resolve"
)

func lookupByHash(db *sqlx.DB, hash string) (resolve.Package, error) {
	var found resolve.Package
	err := db.Get(&found, "SELECT * FROM revisions WHERE hash = $1", hash)
	if err != nil {
		return resolve.Package{}, errors.Wrapf(err, "could not look up hash %#v", hash)
	}
	return found, nil
}

func lookupByKey(db *sqlx.DB, key resolve.Key) (resolve.Package, error) {
	var found resolve.Package
	err := db.Get(&found, "SELECT * FROM revisions WHERE package = $1 AND revision = $2", key.Name, key.Revision)
	if err != nil {
		return resolve.Package{}, errors.Wrapf(err, "could not look up package %#v %#v", key.Name, key.Revision)
	}
	return found, nil
}
