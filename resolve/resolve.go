package resolve

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

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
	dir, err := loadPackage(name)
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not load %s", name)
	}

	// Find the git repository (Go packages may be descendants of a git
	// repository).
	repoDir, err := findRepository(dir)
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not find git repository containing %s", dir)
	}

	// Open the git repository.
	// TODO: support VCSes that are not git.
	fs := osfs.New(repoDir)
	dotgit, err := fs.Chroot(".git")
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not open `.git` folder at %s", repoDir)
	}
	storage, err := filesystem.NewStorage(dotgit)
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not initialize filestore at %s", repoDir)
	}
	repo, err := git.Open(storage, dotgit)
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not open git repository at %s", repoDir)
	}

	// Load revision.
	worktree, err := repo.Worktree()
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not get git worktree at %s", repoDir)
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(revision),
	})
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not find revision %s in %s", revision, repoDir)
	}

	h, err := hash.Dir(repoDir)
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

// loadPackage runs `go get ${name}` and returns the absolute directory of the
// downloaded Go package, or an error.
func loadPackage(name string) (string, error) {
	err := exec.Command("go", "get", name).Run()
	if err != nil {
		return "", errors.Wrapf(err, "could not run go get %s", name)
	}
	gopath := os.Getenv("GOPATH")
	return filepath.Join(gopath, "src", name), nil
}

func findRepository(dir string) (string, error) {
	for dir != string(os.PathSeparator) {
		ok, err := hasRepo(dir)
		if err != nil {
			return "", errors.Wrapf(err, "could not find repo containing %s", dir)
		}
		if ok {
			return dir, nil
		}
		dir = filepath.Dir(dir)
	}
	ok, err := hasRepo(dir)
	if err != nil {
		return "", errors.Wrapf(err, "could not find repo containing %s", dir)
	}
	if ok {
		return dir, nil
	}
	return "", errors.Errorf("%s is not contained by a repo", dir)
}

func hasRepo(dir string) (bool, error) {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrapf(err, "could not check for repo in %s", dir)
	}
	return info.IsDir(), nil
}
