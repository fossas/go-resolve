package resolve

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fossas/go-resolve/hash"
	"github.com/pkg/errors"
)

type Package struct {
	Name     string
	Hash     string
	Revision string
}

func Resolve(name, revision string) (Package, error) {
	err := exec.Command("go", "get", name).Run()
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not run go get %s", name)
	}
	gopath := os.Getenv("GOPATH")
	h, err := hash.DirHash(filepath.Join(gopath, name))
	if err != nil {
		return Package{}, errors.Wrapf(err, "could not calculate hash for %s %s", name, revision)
	}
	return Package{
		Name:     name,
		Revision: revision,
		Hash:     h,
	}, nil
}
