package hash

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Dir computes a hash of a given directory. It does this by traversing the
// directory and hashing all contents except those at `./.git`.
func Dir(dirname string) (string, error) {
	h := sha256.New()

	// This relies on the fact that `filepath.Walk` has a deterministic traversal order.
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "could not read %s", path)
		}
		relPath, err := filepath.Rel(dirname, path)
		if err != nil {
			return errors.Wrapf(err, "could not get path of %s relative to %s", path, dirname)
		}
		if info.IsDir() && relPath == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return errors.Wrapf(err, "could not open %s", path)
			}
			defer file.Close()
			if _, err := io.Copy(h, file); err != nil {
				return errors.Wrapf(err, "could not read contents of %s", path)
			}
		}
		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "could not walk directory")
	}
	return base64.StdEncoding.EncodeToString(h.Sum([]byte{})), nil
}
