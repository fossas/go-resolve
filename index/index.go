// Package index implements package revision indexing.
package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/ilikebits/go-core/log"
	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/fossas/go-resolve/hash"
	"github.com/fossas/go-resolve/models"
)

// Repository indexes all revisions of all packages in a single repository that
// contains a Go package.
//
// This implementation is thread-safe. It sets up a temporary GOPATH on each
// invocation to run `go get`.
func Repository(importpath string) ([]models.Package, error) {
	// TODO: thread a context through all of these.
	// TODO: support more than just Git.

	// Pick a random temporary directory to download into. This prevents two
	// different threads from clobbering each others' $GOPATHs.
	gopath, err := ioutil.TempDir("", "go-resolve-")
	if err != nil {
		return nil, errors.Wrap(err, "could not get temp dir")
	}
	defer func() {
		if err != nil {
			err = os.RemoveAll(gopath)
		}
	}()

	// Download the repository.
	cmd := exec.Command("go", "get", importpath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath))
	err = cmd.Run()
	if err != nil {
		return nil, errors.Wrap(err, "could not run `go get`")
	}

	// Find repository containing the package.
	pkgpath := filepath.Join(gopath, "src", importpath)
	vcs, repopath, err := FindRepository(pkgpath)
	if err != nil {
		return nil, err
	}

	// Open repository.
	repo, err := git.PlainOpen(repopath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open repository")
	}
	logItr, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not open log")
	}
	defer logItr.Close()
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "could not open worktree")
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, errors.Wrap(err, "could not open remote `origin`")
	}
	tagItr, err := repo.Tags()
	if err != nil {
		return nil, errors.Wrap(err, "could not open tags")
	}

	// Get tag set.
	tags := make(map[string]string)
	tagItr.ForEach(func(tag *plumbing.Reference) error {
		name := tag.Name()
		if name.IsTag() {
			tags[tag.Hash().String()] = name.Short()
		}
		return nil
	})

	// Compute hashes for all revisions.
	// TODO: set package versions from tags.
	log.Debug().Msg("computing hashes")
	var output []models.Package
	err = logItr.ForEach(func(commit *object.Commit) error {
		h := commit.Hash

		// Check out revision.
		err := worktree.Checkout(&git.CheckoutOptions{
			Hash: h,
		})
		if err != nil {
			return errors.Wrap(err, "unable to checkout revision during iteration")
		}

		// Compute hashes for all packages within repository.
		pkgs, err := hash.Dir(gopath, repopath)
		if err != nil {
			return err
		}

		// Set revision metadata.
		repoURL := remote.Config().URLs[0]
		for _, pkg := range pkgs {
			rev := h.String()
			pkg.VCS = vcs
			pkg.Repository = repoURL
			pkg.Revision = rev
			if tag, ok := tags[rev]; ok {
				pkg.Version = tag
			}
			output = append(output)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not iterate over revisions")
	}
	return output, nil
}

// FindRepository finds the repository containing a directory.
func FindRepository(dirname string) (models.VCS, string, error) {
	// TODO: support more than just Git.

	log.Debug().Str("dirname", dirname).Msg("FindRepository")
	abs, err := filepath.Abs(dirname)
	if err != nil {
		return -1, "", errors.Wrap(err, "could not get absolute path")
	}
	for curr := abs; curr != "." && curr != "/"; curr = filepath.Dir(curr) {
		vcs := filepath.Join(curr, ".git")

		log.Debug().Str("name", vcs).Msg("os.Stat")
		info, err := os.Stat(vcs)
		log.Debug().Str("info", fmt.Sprintf("%#v", info)).Err(err).Msg("os.Stat result")
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return -1, "", errors.Wrap(err, "could not stat")
		}

		if info.IsDir() {
			return models.Git, curr, nil
		}
	}
	return -1, "", errors.New("could not find repository")
}
