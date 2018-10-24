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
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

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
