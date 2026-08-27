# Mochi

[![codecov](https://codecov.io/gh/yox5ro/mochi/graph/badge.svg?token=FNPYYR374B)](https://codecov.io/gh/yox5ro/mochi)

## Overview

Mochi is yet another KVS.  
Mochi is a hobby project.  
Mochi is not production ready.  

NOTE: Mochi currently serves as in-memory KVS. Data will not be persisted.

## Installation

Clone this repository and build by yourself and put the binary into your `$PATH`. You will need compatible Go compiler. See `CONTRIBUTING.md` for build details.

## Usage

Mochi runs as CLI tool.

```sh
$ mochi -help
Usage of mochi:
  -port string
        port to use when protocol is HTTP (default ":8080")
```

### Starting server

Currently Mochi serves as HTTP server.

```sh
# start HTTP server on port :8080 (default)
$ mochi
```

### Request and response semantics

All request must be directed to `/` with single `key` query parameter. `key` can accept empty string.  
If any error happens, `Mochi-Error-Code` response header will be returned.  
`Mochi-Error-Code` will be one of followings:
- `key-invalid`: `key` is invalid
- `key-not-found`: value of corresponding `key` is not found
- `path-not-found`: request path is not `/`
- `request-method-invalid`: request path is `/`, but request method is not one of the `GET`, `PUT`, `DELETE`. In this case, `Allow: GET, PUT, DELETE` header will also be returned
- `internal`: any other internal error

Response body will be empty unless successful get request

### Get request

Request must be directed to `GET /`.
When `key` is found, corresponding byte sequence will be returned on response body.
When `key` is not found, error will be returned as `Mochi-Error-Code: key-not-found` and status code will be 404, since the client is explicitly requesting a value that does not exist.  

### Put request

Request must be directed to `PUT /` with optional request body.
If request body is empty, the value of corresponding key will be upserted to empty byte sequence, otherwise upserted to request body byte sequence.

### Delete request

Request must be directed to `DELETE /`.
When `key` is found, `key` and corresponding value will be removed.  
When `key` is not found, the status code will be 204, since the desired final state (key absent) is already achieved.  

### HTTP status code

404 if following condition is met:
- any request failed with `path-not-found` error or get request failed with `key-not-found` error

405 if following condition is met:
- any request failed with `request-method-invalid` error

400 if following conditions are all met:
- request path is `/` and HTTP method is one of `GET`, `PUT`, `DELETE`
- request failed with `key-invalid` error

500 if following conditions are all met:
- request path is `/` and HTTP method is one of `GET`, `PUT`, `DELETE`
- request failed with `internal` error

204 if following conditions are all met
- request path is `/` and HTTP method is one of `PUT`, `DELETE`
- request is successful

otherwise, 200.

## Concurrency

Mochi currently supports linearizability for operations on the same key.

## Contributing

See `CONTRIBUTING.md`
