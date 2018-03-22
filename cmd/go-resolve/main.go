package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Prints the usage string.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: %s [flags] <package>\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	// TODO(@ilikebits): implement these flags when the API is implemented.
	flag.String("package", "", "The package name to look up")
	flag.String("hash", "", "The git tree-hash to look up")
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	target, err := parseTarget(flag.Arg(0))
	if err != nil {
		log.Fatalf("Invalid input: %s", err.Error())
	}

	hash, err := getTreeHash(target)
	if err != nil {
		log.Fatalf("Could not compute tree hash: %s", err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %s", err.Error())
	}
	gopath := os.Getenv("GOPATH")
	importPath, err := filepath.Rel(filepath.Join(gopath, "src"), cwd)
	if err != nil {
		log.Fatalf("Could not compute package import path: %s", err.Error())
	}
	fmt.Printf("%s %s\n", importPath, hash)
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
		targetPath = filepath.Join(gopath, "src", target)
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
