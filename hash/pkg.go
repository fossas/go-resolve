package hash

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/fossas/go-resolve/models"
)

// Package computes the package hash of the given import path.
func Package(importpath string) (models.Package, error) {
	// Run `go list -json <importpath>`.
	cmd := exec.Command("go", "list", "-json", importpath)
	out, err := cmd.Output()
	if err != nil {
		return models.Package{}, errors.Wrapf(err, "could not run `go list` for import path %#v", importpath)
	}

	// Unmarshal results.
	var pkg goPkg
	err = json.Unmarshal(out, &pkg)
	if err != nil {
		return models.Package{}, errors.Wrapf(err, "could not unmarshal `go list` output for import path %#v", importpath)
	}

	// Compute hash.
	hash, err := Hash(pkg.Dir, pkg.SourceFiles())
	if err != nil {
		return models.Package{}, err
	}

	return models.Package{
		ImportPath: pkg.ImportPath,
		Hash:       hash,

		SourceFiles: pkg.SourceFiles(),
		Imports:     pkg.Imports,
		Deps:        pkg.Deps,
	}, nil
}
