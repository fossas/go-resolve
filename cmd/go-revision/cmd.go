package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// Prints a formatted error message, then exits with error code 1.
func fatal(format string, a ...interface{}) {
	fmt.Fprintf(flag.CommandLine.Output(), format, a...)
	os.Exit(1)
}

// Prints the usage string.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: go-revision <package>\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	target, err := parseTarget(flag.Arg(0))
	if err != nil {
		fatal("Invalid input: %s", err.Error())
	}

	r, err := git.PlainOpen(target)
	if err == git.ErrRepositoryNotExists {
		dir, err := os.Open(target)
		if err != nil {
			fatal("Could not open package directory: %s", err.Error())
		}
		files, err := dir.Readdirnames(0)
		if err != nil {
			fatal("Could not read package file names: %s", err.Error())
		}

		// Package is not in a git repository, so we need to calculate the tree hash ourselves.
		fs := osfs.New(target)
		dot, err := fs.Chroot(".git")
		if err != nil {
			fatal("Could not initialize filesystem: %s", err.Error())
		}
		storage, err := filesystem.NewStorage(dot)
		if err != nil {
			fatal("Could not initialize filestore: %s", err.Error())
		}
		r, err = git.Init(storage, fs)
		if err != nil {
			fatal("Could not init repository: %s", err.Error())
		}
		w, err := r.Worktree()
		if err != nil {
			fatal("Could not get worktree: %s", err.Error())
		}
		// By default, `go-git` will add the `.git` folder if you try to `w.Add(".")`. It also doesn't fully support
		// `git rm --cached .git`, `git reset HEAD .git`, or `git checkout .git` so there's no way to add all and then
		// remove. This works around `go-git` by explicitly only adding non-`.git` files.
		for _, file := range files {
			_, err = w.Add(file)
			if err != nil {
				fatal("Could not add package: %s", err.Error())
			}
		}
		h, err := w.Commit("go-revision commit", &git.CommitOptions{
			Author: &object.Signature{
				Name: "go-revision",
				When: time.Now(),
			},
		})
		if err != nil {
			fatal("Could not commit package: %s", err.Error())
		}
		commit, err := r.CommitObject(h)
		if err != nil {
			fatal("Could not get commit object from new hash (%s): %s", h, err.Error())
		}
		fmt.Println(commit.TreeHash)
		err = os.RemoveAll(filepath.Join(target, ".git"))
		if err != nil {
			fatal("Could not clean up git repository: %s", err.Error())
		}
	} else if err == nil {
		// Package is in a git repository, so we can look up the tree hash.
		head, err := r.Head()
		if err != nil {
			fatal("Could not get HEAD: %s", err.Error())
		}
		commit, err := r.CommitObject(head.Hash())
		if err != nil {
			fatal("Could not get commit object from HEAD hash (%s): %s", head.Hash(), err.Error())
		}
		fmt.Println(commit.TreeHash)
	} else {
		fatal("Could not open git repository: %s", err.Error())
	}
}

// Takes a target string (which can be either an absolute path, a relative path, or a Go import path) and return an
// absolute path to the target directory, or an error if the target is invalid.
func parseTarget(target string) (string, error) {
	if target == "" {
		return "", errors.New("no package specified")
	}

	var targetPath string
	switch target[0] {
	case '/':
		// Assume this is an absolute path
		targetPath = target
	case '.':
		// Assume this is a relative path
		cwd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "could not get working directory")
		}
		targetPath = filepath.Join(cwd, target)
	default:
		// Assume this is a Go import path
		gopath := os.Getenv("GOPATH")
		targetPath = filepath.Join(gopath, target)
	}

	stat, err := os.Stat(targetPath)
	if err != nil {
		return "", errors.Wrap(err, "could not stat package directory")
	}
	if !stat.IsDir() {
		return "", errors.New("specified path is not a package")
	}

	return targetPath, nil
}
