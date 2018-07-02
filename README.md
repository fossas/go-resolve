# go-resolve

## Quickstart

```bash
go get -u github.com/fossas/go-resolve/cmd/go-resolve

go-resolve ./vendor/github.com/user/project/package
```

## Overview

`go-resolve` is a set of three components: `go-resolve`, `go-resolve-api`, and
`go-resolve-worker`, that together compose a Go package revision resolver.

- `go-resolve` is a command-line utility for computing and looking up a package
  hash on a `go-resolve-api` server.
- `go-resolve-api` is an API server for queuing and querying hashes.
- `go-resolve-worker` is an asynchronous worker for computing package hashes.

## Usage

### Hash lookup

Hash lookups attempt to resolve a package's revision given its package hash.

```bash
$ go-resolve github.com/project/foo 
{
  "packages": {
    "github.com/foo/bar/baz": {
      "found": true,
      "repository": "github.com/foo/bar",
      "ambiguous": false, # may be true if best revision is ambiguous
      "revision": {
        # the best revision
        "hash": "abcd",
        "timestamp": "January 1st, 2018"
      },
      "candidates": [
        # revision candidates, from newest to oldest
        {
          "hash": "abcd",
          "timestamp": "January 1st, 2018",
          "matches": [
            "github.com/foo/bar",
            "github.com/foo/bar/baz"
          ]
        }
      ]
    }
  }
}
```

### Multi-hash lookup

Multi-hash lookups attempt to resolve the revisions of packages at multiple Go
import paths. When package hashes are ambiguous, a multi-hash lookup will
attempt to disambiguate multiple candidate revisions from the same repository
by preferring the revision which matches the most package hashes. See [package
hash ambiguity]() for details.

```bash
$ go-resolve github.com/project/foo github.com/baz/quux
{"ok": false, ""}
```

### Vendor folder lookup

Vendor lookups attempt to resolve the revisions of all Go packages within a
folder. This is equivalent to a multi-hash lookup of all Go packages within the
folder.

```bash
$ go-resolve -vendor ./vendor
{"ok": false, ""}
```

### Local modification lookup

Local modification lookups check whether a package's actual, on-disk hash
matches its expected hash given an expected revision.

```bash
$ go-resolve -revision expected-revision-commit-hash github.com/project/foo 
{"ok": true, ""}

$ go-resolve -revision expected-revision-commit-hash github.com/project/foo 
{"ok": false, ""}

$ go-resolve -version expected-version-tag github.com/project/foo 
{"ok": true, ""}
```

## Details

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

Package hashes may be _ambiguous_. That is, a package hash may be present in
multiple revisions, because new revisions do not necessarily change all packages
within the repository. For single hash lookups, this means that a package hash
lookup may return multiple possible revisions.

For multi-hash lookups, `go-resolve` attempts to disambiguate package hashes by
preferring revisions which match all package hashes of packages within the same
repository.

```sql
SELECT * FROM packages GROUP BY packages.repository_id;
(SELECT * FROM revisions JOIN packages ON revisions.revision_id = packages.revision_id WHERE packages.hash = $1)
INTERSECT
(SELECT * FROM revisions JOIN packages ON revisions.revision_id = packages.revision_id WHERE packages.hash = $2);
```