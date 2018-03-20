package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func getTreeHash(dirname string) (string, error) {
	r, err := git.PlainOpen(dirname)
	if err == git.ErrRepositoryNotExists {
		dir, err := os.Open(dirname)
		if err != nil {
			return "", errors.Wrap(err, "could not open package directory")
		}
		files, err := dir.Readdirnames(0)
		if err != nil {
			return "", errors.Wrap(err, "could not read package file names")
		}

		// Package is not in a git repository, so we need to calculate the tree hash ourselves.
		fs := osfs.New(dirname)
		dot, err := fs.Chroot(".git")
		if err != nil {
			return "", errors.Wrap(err, "could not initialize filesystem")
		}
		storage, err := filesystem.NewStorage(dot)
		if err != nil {
			return "", errors.Wrap(err, "could not initialize filestore")
		}
		r, err = git.Init(storage, fs)
		if err != nil {
			return "", errors.Wrap(err, "could not init repository")
		}
		w, err := r.Worktree()
		if err != nil {
			return "", errors.Wrap(err, "could not get worktree")
		}
		// By default, `go-git` will add the `.git` folder if you try to `w.Add(".")`. It also doesn't fully support
		// `git rm --cached .git`, `git reset HEAD .git`, or `git checkout .git` so there's no way to add all and then
		// remove. This works around `go-git` by explicitly only adding non-`.git` files.
		for _, file := range files {
			_, err = w.Add(file)
			if err != nil {
				return "", errors.Wrap(err, "could not add package")
			}
		}
		h, err := w.Commit("go-revision commit", &git.CommitOptions{
			Author: &object.Signature{
				Name: "go-revision",
				When: time.Now(),
			},
		})
		if err != nil {
			return "", errors.Wrap(err, "could not commit package")
		}
		commit, err := r.CommitObject(h)
		if err != nil {
			return "", errors.Wrapf(err, "could not get commit object from new hash (%s)", h)
		}
		err = os.RemoveAll(filepath.Join(dirname, ".git"))
		if err != nil {
			return "", errors.Wrap(err, "could not clean up git repository")
		}
		return commit.TreeHash.String(), nil
	} else if err == nil {
		// Package is in a git repository, so we can look up the tree hash.
		head, err := r.Head()
		if err != nil {
			return "", errors.Wrap(err, "could not get HEAD")
		}
		commit, err := r.CommitObject(head.Hash())
		if err != nil {
			return "", errors.Wrapf(err, "could not get commit object from HEAD hash (%s)", head.Hash())
		}
		return commit.TreeHash.String(), nil
	} else {
		return "", errors.Wrap(err, "could not open git repository")
	}
}
