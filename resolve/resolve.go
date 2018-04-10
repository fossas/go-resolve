package resolve

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/fossas/go-resolve/hash"
)

// Key uniquely identifies a single package and revision.
type Key struct {
	Name     string `db:"package"`
	Revision string
}

// Package contains both a package key and its resolved hash.
type Package struct {
	Key
	Hash string
}

// All resolves every revision of a single package.
func All(name string) ([]Package, error) {
	// To make this fast, we'll need to make this thread-safe. Here's the game
	// plan:
	//
	// 1. Load the package into the filesystem.
	// 2. Read the log of commits.
	// 3. For each commit, spawn a new goroutine that:
	//    1. Copies the filesystem's contents into memory (to prevent goroutines
	//			 from clobbering each other)
	//    2. Checks out the specified revision and returns a (hash, error).
	// 4. Put those all together into a ([]Package, error).
	return nil, errors.New("not implemented")
}

// Single resolves a single package and revision.
// NOTE: this is not thread-safe because it uses the actual filesystem to check
// out and hash files instead of an in-memory copy.
func Single(name, revision string) (Package, error) {
	// Load the directory.
	err := exec.Command("go", "get", name).Run()
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not run `go get %s`", name)
	}
	gopath := os.Getenv("GOPATH")
	dir := filepath.Join(gopath, "src", name)

	// Load revision.
	// TODO: support VCSes that are not git.
	cmd := exec.Command("git", "checkout", revision)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not run `git checkout %s`", revision)
	}

	// Compute hash
	h, err := hash.Dir(dir)
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not calculate hash for %s %s", name, revision)
	}

	return Package{
		Key: Key{
			Name:     name,
			Revision: revision,
		},
		Hash: h,
	}, nil
}
