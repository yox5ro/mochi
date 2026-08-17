# Mochi

## Overview

Mochi is yet another KVS.  
Mochi is a hobby project.  
Mochi is not production ready.  

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

## Requests and Responses

Clients can do three kind of operations (get, put, delete). All these operations are designed as operations that specifies the desired final state and idempotent.  
When Mochi gets invalid `op`, `error` will be returned.

Request must be following JSON.

```json
{
  "op": "string", // ("get", "put", "delete")
  "key": "string | null", // required only when `op` is "get" or "delete"
  "value": "string | null" // required only when `op` is "put"
}
```

Response type will be below JSON.

```json
{
  "value": "string | null", // returned only when `op` is "get"
  "error": "string | null" // returned only when any error occurs
}
```

### HTTP status code

If request path is not `POST /`, status code will be 404.  
If request path is `/` but method is not `POST`, status code will be 405.  
If request body is not valid JSON, status code will be 400.  
Otherwise (i.e. the request is a valid `POST /` request with a well-formed JSON body), status code always will be 200, regardless of the `op` result — even `invalid op` or `key not found` cases are reported via the `error` field in the response body.  
If an unexpected server error occurs, status code will be 500.

### Endpoint

All requests must be directed to the root path of the HTTP server (`/`).

### Get request

When `key` is found, `value` will be returned.  
When `key` is not found, `error` will be returned, since the client is explicitly requesting a value that does not exist.

### Put request

Regardless of existence of `key`, the value will be replaced (created) as `value`.  
When `value` is `null`, `error` will be returned.
Response will be empty object unless error occurs.

### Delete request

When `key` is found, `key` and corresponding `value` will be removed.  
When `key` is not found, `error` will not be returned, since the desired final state (key absent) is already achieved.
Response will be empty object unless error occurs.

## Concurrency

Operation on a single key appear atomic.

## Contributing

See `CONTRIBUTING.md`
