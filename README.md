# go-resolve

## Quickstart

`go-resolve` is a command line tool for finding the revisions of Go packages,
even in projects without build tool support. It works by looking up package
hashes against the `go-resolve` API.

```bash
$ go get -u github.com/fossas/go-resolve/cmd/go-resolve

$ go-resolve ./vendor/github.com/user/project/package
{
  "Hashes": {
    "notarealhash": "github.com/user/project/package"
  },
  "Packages": {
    "github.com/user/project/package": {
      "Ambiguous": true,
      "Revision": {
        "VCS": "git",
        "Repository": "github.com/user/project",
        "Hash": "notarealhash",
        "Message": "this is the commit message",
        "Timestamp": "January 1st, 2018"
      },
      "Candidates": [ /* Other Revisions */ ]
    }
  }
}
```

## Usage

### Package lookup

Package lookups attempt to resolve a package's revision given its import path.

```bash
$ go-resolve ./vendor/github.com/project/foo/baz
{
  "Hashes": {
    "notarealhash": "github.com/user/project/package"
  },
  "Packages": {
    "github.com/user/project/package": {
      "Ambiguous": true,
      "Revision": {
        "VCS": "git",
        "Repository": "github.com/user/project",
        "Hash": "notarealhash",
        "Message": "this is the commit message",
        "Timestamp": "January 1st, 2018"
      },
      "Candidates": [ /* Other Revisions */ ]
    }
  }
}
```

<!-- ### Multi-package lookup

Multi-package lookups attempt to resolve the revisions of multiple packages at
specified by Go import paths. When package hashes are ambiguous, a multi-hash
lookup will attempt to disambiguate multiple candidate revisions from the same
repository by preferring the revision which matches the most package hashes.
See [package hash ambiguity]() below for details.

```bash
$ go-resolve github.com/user/project/package ./vendor/foo/bar ./baz/quux/...
``` -->

<!-- ### Integrity lookup

Integrity lookups check whether a package's actual, on-disk hash matches its
expected hash given an expected revision.

```bash
$ go-resolve -revision expected-revision-commit-hash github.com/project/foo 
{"ok": true, ""}

$ go-resolve -revision expected-revision-commit-hash github.com/project/foo 
{"ok": false, ""}

$ go-resolve -version expected-version-tag github.com/project/foo 
{"ok": true, ""}
``` -->

## Running your own resolver

This project also provides a resolver API implementation:

- `github.com/fossas/go-resolve/cmd/go-resolve-api` is an API server for
  querying hashes.
- `github.com/fossas/go-resolve/cmd/go-resolve-worker` is an asynchronous worker
  that computes package hashes and indexes Go projects.

Instructions for running your own resolver are still WIP. See the Dockerfile and
Makefile to get started.

## Design

### Why `go-resolve`?

Idiomatic Go projects use dependencies by vendoring them. This often occurs with
build tool support, but sometimes there is not an easy way to look up a
package's revision. This can have many causes:

- No build tool was used.
- The build tool does not reproducibly specify revisions for some transitive
  dependencies.
- The build tool is not supported.
- The dependency manifest is missing or corrupted.

`go-resolve` was written to address these corner cases. Also, it's a fun project
tackling a well-defined and well-scoped technical problem.

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

### Package hash ambiguity

Package hashes may be _ambiguous_. A package hash may be present in multiple
revisions, because repositories contain many packages and new revisions
do not necessarily change all packages within the repository. For single-package
lookups, this means that a package hash lookup may return multiple possible
revisions.

For multi-hash lookups, `go-resolve` attempts to disambiguate package hashes by
choosing revisions which match more package hashes of packages within the same
revision.

### Revision selection

When multiple revisions may match a single package, a "best" revision is
selected by the following criteria:

1. If performing a multi-hash lookup, sort revisions in descending order by the
  number of packages they match.
2. To tiebreak, sort revisions in descending order by time.
3. Pick the first revision.
