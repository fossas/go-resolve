package api

import (
	"errors"

	"github.com/fossas/go-resolve/resolve"
	"github.com/go-redis/redis"
)

var (
	errPackageNotFound          = errors.New("could not find package")
	errPackageIncorrectChecksum = errors.New("package checksum is different")
)

// Redis schema: string(Hash) => string(Name + " " + Revision)

func lookup(cache *redis.Client, hash string) (resolve.Package, error) {
	// TODO: implement this.
	return resolve.Package{}, nil
}

func verify(cache *redis.Client, actual resolve.Package) (resolve.Package, error) {
	// TODO: implement this.
	return resolve.Package{}, nil
}
