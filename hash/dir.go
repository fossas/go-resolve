package hash

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/fossas/go-resolve/models"
)

// Dir computes the package hashes of all packages within a given
// directory.
func Dir(gopath, dirname string) ([]models.Package, error) {
	log.WithField("dirname", dirname).Debug("hashing packages within directory")

	// Run `go list -json <dirname>/...` to get source files
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Dir = dirname
	log.WithField("cmd", fmt.Sprintf("%#v", cmd)).Debug("go list command")

	out, err := cmd.Output()
	if err != nil {
		// Ignore exit errors: `go list` can partially succeed, and will still exit
		// 1 on warnings.
		log.WithError(err).Warn("go list error")
		if e, ok := err.(*exec.ExitError); !ok {
			log.WithError(e).WithField("stderr", string(e.Stderr)).Error("go list error output")
			return nil, errors.Wrapf(err, "could not run `go list` in directory %#v", dirname)
		}
	}
	log.WithField("output", string(out)).Debug("raw go list output")

	// The output of `go list -json <packages>` with more than one package is
	// actually newline-delimited JSON objects, not one valid JSON object.
	outstr := string(out)
	withCommas := strings.Replace(outstr, "}\n{", "},{", -1)
	data := []byte("[" + withCommas + "]")

	// Unmarshal results.
	var pkgs []goPkg
	err = json.Unmarshal(data, &pkgs)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal output of running `go list` in directory %#v", dirname)
	}
	log.WithField("pkgs", pkgs).Debug("unmarshalled go list output")

	// Compute hash for each package.
	log.Debug("computing hashes")
	var hashed []models.Package
	for _, pkg := range pkgs {
		hash, err := Hash(pkg.Dir, pkg.SourceFiles())
		if err != nil {
			log.WithError(err).WithField("pkg", pkg).Error("could not compute hash")
			return nil, err
		}
		log.WithFields(log.Fields{
			"hash":       hash,
			"importpath": pkg.ImportPath,
			"dir":        pkg.Dir,
		}).Debug("computed hash")
		hashed = append(hashed, models.Package{
			ImportPath: pkg.ImportPath,
			Hash:       hash,
		})
	}

	return hashed, nil
}
