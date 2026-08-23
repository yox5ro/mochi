# Mochi

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

### HTTP status code

404 if following condition is met:
- request path is not `/`
- response body will be empty

405 if following condition is met:
- request path is `/` but HTTP method is not `POST`
- response body will be empty

400 if following conditions are all met:
- request path is `/` and HTTP method is `POST`
- request body is invalid JSON, or empty, or does not consist of exactly one JSON value, or has duplicate object member names
- response body will be empty

500 if following condition is met:
- request path is `/` and HTTP method is `POST`
- unexpected server error occurs
- response body will be empty

otherwise, 200.

If valid JSON does not match any of get, put, delete request body format, 200 will be returned with following response body.

```json
{
  "error": "string"
}
```

### Endpoint

All requests must be directed to the root path of the HTTP server (`/`).

### Get request

Request must be below format. Other fields than the format will be ignored.
- `"op"` must be `"get"`
- `"key"` allows empty string

```json
{
  "op": "get",
  "key": "string"
}
```

When `key` is found, response will be below format

```json
{
  "value": "string"
}
```

When `key` is not found, `error` will be returned, since the client is explicitly requesting a value that does not exist.  
When expected error happens or request does not satisfies the format (e.g. type mismatch, null value, missing field), `error` will be returned.

```json
{
  "error": "string"
}
```

### Put request

Request must be below format. Other fields than the format will be ignored.
- `"op"` must be `"put"`
- `"key"` allows empty string
- `"value"` allows empty string

```json
{
  "op": "put",
  "key": "string",
  "value": "string"
}
```

Regardless of existence of `key`, the value will be replaced (created) as `value`. 
If the request is successful, empty object will be returned.

```json
{}
```
When expected error happens or request does not satisfies the format (e.g. type mismatch, null value, missing field), `error` will be returned.

```json
{
  "error": "string"
}
```

### Delete request

Request must be below format. Other fields than the format will be ignored.
- `"op"` must be `"delete"`
- `"key"` allows empty string

```json
{
  "op": "delete",
  "key": "string"
}
```

When `key` is found, `key` and corresponding `value` will be removed.  
When `key` is not found, `error` will not be returned, since the desired final state (key absent) is already achieved.  
In both case, request are successful and empty object will be returned.


```json
{}
```

When expected error happens or request does not satisfies the format (e.g. type mismatch, null value, missing field), `error` will be returned.

```json
{
  "error": "string"
}
```

## Concurrency

Mochi currently supports linearizability for operations on the same key.

## Contributing

See `CONTRIBUTING.md`
