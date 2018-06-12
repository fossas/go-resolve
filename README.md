# go-resolve

## Quickstart

```bash
go get -u github.com/fossas/go-resolve/cmd/go-resolve

go-resolve ./vendor/github.com/user/project/package
```

## Overview

`go-resolve` is a set of three components: `go-resolve`, `go-resolve-api`, and
`go-resolve-worker`, that together compose a Go package revision resolver.

- `go-resolve` is a command-line tool for computing and looking up a package
  hash on a `go-resolve-api` server.
- `go-resolve-api` is an API server for queuing and querying hashes.
- `go-resolve-worker` is an asynchronous worker for computing package hashes.

## Design

### Package hashing

Packages are hashed by running `go list -json <package>` to get their source
files, then computing the SHA256 of the files in lexicographical order.

Source files are the union of a package's:

- `GoFiles`
- `CgoFiles`
- `CFiles`
- `CXXFiles`
- `MFiles`
- `HFiles`
- `FFiles`
- `SFiles`
- `SwigFiles`
- `SwigCXXFiles`
- `SysoFiles`
