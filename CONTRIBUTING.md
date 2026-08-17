# Contributing Mochi

## Prerequisites

- Go compiler. For version information, see `go.mod`.
- golangci-lint.

## Development

Build
```sh
$ go build
```

Running test

```sh
$ go test ./...
```

Running linter

```
$ golangci-lint run
```

Running formatter

```
$ golangci-lint fmt
```