// Package index implements package revision indexing.
package index

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/fossas/go-resolve/hash"
	"github.com/fossas/go-resolve/models"
)

// Repository indexes all revisions of all packages in a single repository that
// contains a Go package.
func Repository(importpath string, handler func(pkgs []models.Package) error) error {
	// Pick a random temporary directory to download into. This prevents two
	// different threads from clobbering each others' $GOPATHs.
	gopath, err := ioutil.TempDir("", "go-resolve-")
	if err != nil {
		return errors.Wrap(err, "could not get temp dir")
	}

	// Download the repository.
	cmd := exec.Command("go", "get", importpath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath))
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "could not run `go get`")
	}

	// Find repository containing the package.
	packagepath := filepath.Join(gopath, "src", importpath)
	repopath, err := FindRepository(packagepath)
	if err != nil {
		return err
	}

	// Open repository.
	repo, err := git.PlainOpen(repopath)
	if err != nil {
		return errors.Wrap(err, "could not open repository")
	}
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return errors.Wrap(err, "could not open log")
	}
	defer iter.Close()

	// Compute hashes for all revisions.
	err = iter.ForEach(func(commit *object.Commit) error {
		// Compute hashes for all packages within repository.
		pkgs, err := hash.Dir(repopath)
		if err != nil {
			return err
		}

		// Set package revisions.
		// TODO: set package versions from tags.
		for i := range pkgs {
			pkgs[i].Revision = commit.Hash.String()
		}

		// Upload hashes for this revision.
		err = handler(pkgs)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "could not iterate over revisions")
	}
	return nil
}

// FindRepository finds the repository containing a directory.
func FindRepository(dirname string) (string, error) {
	log.Printf("FindRepository(%#v)", dirname)
	for curr := dirname; curr != "." && curr != "/"; curr = filepath.Dir(curr) {
		log.Printf("os.Stat(%#v)", filepath.Join(curr, ".git"))
		info, err := os.Stat(filepath.Join(curr, ".git"))
		log.Printf("info = %#v", info)
		log.Printf("err = %#v", err)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "could not stat")
		}
		if info.IsDir() {
			return curr, nil
		}
	}
	return "", errors.New("could not find repository")
}
