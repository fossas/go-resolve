// Package go-resolve implements the go-resolve command line utility.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fossas/go-resolve/hash"
)

// Prints the usage string.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: %s [flags] <import path>\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	apiURL := flag.String("api", "http://localhost:8080", "The API URL")
	verbose := flag.Bool("v", false, "Use verbose logging")
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	h, err := hash.Package(flag.Arg(0))
	if err != nil {
		log.Fatalf("Could not compute package hash: %s", err.Error())
	}
	if *verbose {
		log.Printf("Computed import path: %s\n", h.ImportPath)
		log.Printf("Computed hash: %s\n", h.Hash)
	}

	res, err := http.Get(*apiURL + "/api/lookup/" + h.Hash)
	if err != nil {
		log.Fatalf("API error: %s", err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Could not read body: %s", err.Error())
	}
	fmt.Printf("%s\n", body)
}
