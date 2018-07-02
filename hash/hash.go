// Package hash implements package hashing.
//
// Packages are hashed by running `go list -json <package>` to get their source
// files, then computing the SHA256 of the files in lexicographical order.
//
// Source files are the union of a package's:
//
//   - GoFiles
//   - CgoFiles
//   - CFiles
//   - CXXFiles
//   - MFiles
//   - HFiles
//   - FFiles
//   - SFiles
//   - SwigFiles
//   - SwigCXXFiles
//   - SysoFiles
package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ilikebits/go-core/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/fossas/go-resolve/models"
)

// goPkg is a struct of `go list -json` output.
type goPkg struct {
	ImportPath string
	Dir        string

	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	CXXFiles     []string
	MFiles       []string
	HFiles       []string
	FFiles       []string
	SFiles       []string
	SwigFiles    []string
	SwigCXXFiles []string
	SysoFiles    []string
}

func (pkg *goPkg) files() []string {
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.CFiles...)
	files = append(files, pkg.CXXFiles...)
	files = append(files, pkg.MFiles...)
	files = append(files, pkg.HFiles...)
	files = append(files, pkg.FFiles...)
	files = append(files, pkg.SFiles...)
	files = append(files, pkg.SwigFiles...)
	files = append(files, pkg.SwigCXXFiles...)
	files = append(files, pkg.SysoFiles...)
	return files
}

// Package computes the package hash of the given import path.
func Package(importpath string) (models.Package, error) {
	// Run `go list -json <importpath>`
	cmd := exec.Command("go", "list", "-json", importpath)
	out, err := cmd.Output()
	if err != nil {
		return models.Package{}, errors.Wrap(err, "could not run `go list` for single package")
	}

	// Unmarshal results.
	var pkg goPkg
	err = json.Unmarshal(out, &pkg)
	if err != nil {
		return models.Package{}, errors.Wrap(err, "could not unmarshal `go list` output for single package")
	}

	// Compute hash.
	hash, err := Hash(pkg.Dir, pkg.files())
	if err != nil {
		return models.Package{}, err
	}

	return models.Package{
		ImportPath: pkg.ImportPath,
		Hash:       hash,
	}, nil
}

// Dir computes the package hashes of all packages within a given
// directory.
func Dir(gopath, dirname string) ([]models.Package, error) {
	log.Debug().Str("dirname", dirname).Msg("hash.Dir")
	// Run `go list -json <dirname>/...`
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Dir = dirname
	log.Debug().Str("cmd", fmt.Sprintf("%#v", cmd)).Msg("raw command")
	out, err := cmd.Output()
	if err != nil {
		// Ignore exit errors: `go list` can partially succeed, and will still exit
		// 1 on warnings.
		log.Error().Err(err).Msg("go list error")
		if e, ok := err.(*exec.ExitError); !ok {
			log.Error().Str("stderr", string(e.Stderr)).Msg("go list stderr output")
			return nil, errors.Wrap(err, "could not run `go list` for multiple packages")
		}
	}
	log.Debug().Str("output", string(out)).Msg("raw go list output")

	// The output of `go list -json <packages>` with more than one package is
	// actually newline-delimited JSON objects, not one valid JSON object.
	outstr := string(out)
	withCommas := strings.Replace(outstr, "}\n{", "},{", -1)
	data := []byte("[" + withCommas + "]")

	// Unmarshal results.
	var pkgs []goPkg
	err = json.Unmarshal(data, &pkgs)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal `go list` output for multiple packages")
	}
	arr := zerolog.Arr()
	for _, pkg := range pkgs {
		arr = arr.Str(pkg.ImportPath)
	}
	log.Debug().Array("packages", arr).Msg("go list ./... output")

	// Compute hash for each package.
	log.Debug().Msg("computing hashes")
	var hashed []models.Package
	for _, pkg := range pkgs {
		hash, err := Hash(pkg.Dir, pkg.files())
		if err != nil {
			log.Error().Err(err).Msg("could not compute hash")
			return nil, err
		}
		log.Debug().Str("hash", hash).Str("package", pkg.ImportPath).Str("dir", pkg.Dir).Msg("computed hash")
		hashed = append(hashed, models.Package{
			ImportPath: pkg.ImportPath,
			Hash:       hash,
		})
	}

	return hashed, nil
}

// Hash computes a package hash from a list of source file names.
//
// NOTE: this will sort the original slice of file names.
func Hash(dirname string, filenames []string) (string, error) {
	// Sort file names.
	sort.Strings(filenames)

	// Hash files.
	h := sha256.New()
	for _, filename := range filenames {
		// We use the extra closure here so the defer runs as soon as the iteration
		// completes. See golang/go#3978.
		err := func() error {
			file, err := os.Open(filepath.Join(dirname, filename))
			if err != nil {
				return errors.Wrap(err, "could not open source file")
			}
			defer file.Close()
			_, err = io.Copy(h, file)
			if err != nil {
				return errors.Wrap(err, "could not hash file contents")
			}
			return nil
		}()
		if err != nil {
			return "", err
		}
	}

	return base64.StdEncoding.EncodeToString(h.Sum([]byte{})), nil
}
