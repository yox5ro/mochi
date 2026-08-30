package main

import "errors"

var (
	errNotFound = errors.New("key not found")
)

type Store interface {
	get(key string) ([]byte, error)
	put(key string, value []byte) error
	delete(key string) error
}
