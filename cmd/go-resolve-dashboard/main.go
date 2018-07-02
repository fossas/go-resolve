// Package go-resolve-dashboard implements a web dashboard UI for managing a
// go-resolve instance.
package main

import "flag"

func main() {
	flag.String("db", "", "database URL")
	flag.String("api", "", "API URL")
	flag.String("app", "", "location of static web app")

	// static server for web app

	// proxy go-resolve + indexing requests to API

	// provide own handlers for db introspection
}
