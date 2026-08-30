package main

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errNotFound = errors.New("key not found")
)

type Store interface {
	get(key string) ([]byte, error)
	put(key string, value []byte) error
	delete(key string) error
}

type InMemoryMapStore struct {
	store map[string][]byte
}

func newInMemoryMapStore(initialState map[string][]byte) InMemoryMapStore {
	inMemoryMapStore := InMemoryMapStore{store: make(map[string][]byte)}

	for k, v := range initialState {
		inMemoryMapStore.store[k] = slices.Clone(v)
	}

	return inMemoryMapStore
}

func (s InMemoryMapStore) get(key string) ([]byte, error) {
	if v, ok := s.store[key]; ok {
		return slices.Clone(v), nil
	}
	return nil, fmt.Errorf("failed to get key %q: %w", key, errNotFound)
}

func (s InMemoryMapStore) put(key string, value []byte) error {
	s.store[key] = slices.Clone(value)
	return nil
}

func (s InMemoryMapStore) delete(key string) error {
	if _, ok := s.store[key]; !ok {
		return fmt.Errorf("failed to delete key %q: %w", key, errNotFound)
	}
	delete(s.store, key)
	return nil
}
