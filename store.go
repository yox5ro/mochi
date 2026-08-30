package main

import "errors"

const (
	MaxKeyLength   = 1024      // 1KiB
	MaxValueLength = 64 * 1024 // 64KiB
)

var (
	errNotFound      = errors.New("key not found")
	errKeyTooLarge   = errors.New("key too large")
	errValueTooLarge = errors.New("value too large")
)

type Store interface {
	get(key string) ([]byte, error)
	put(key string, value []byte) error
	delete(key string) error
}

func validateKeyLength(key string) error {
	if len(key) > MaxKeyLength {
		return errKeyTooLarge
	}
	return nil
}

func validateValueLength(value []byte) error {
	if len(value) > MaxValueLength {
		return errValueTooLarge
	}
	return nil
}
