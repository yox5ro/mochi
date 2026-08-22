package main

import (
	"errors"
	"fmt"
	"maps"
)

var (
	errNotFound = errors.New("key not found")
)

type Store interface {
	get(key string) (value string, err error)
	put(key, value string) error
	delete(key string) error
}

type InMemoryMapStore struct {
	store map[string]string
}

func newInMemoryMapStore(initialState map[string]string) InMemoryMapStore {
	inMemoryMapStore := InMemoryMapStore{}

	s := maps.Clone(initialState)
	if s == nil {
		s = make(map[string]string)
	}

	inMemoryMapStore.store = s
	return inMemoryMapStore
}

func (s InMemoryMapStore) get(key string) (string, error) {
	if v, ok := s.store[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("failed to get key %q: %w", key, errNotFound)
}

func (s InMemoryMapStore) put(key, value string) error {
	s.store[key] = value
	return nil
}

func (s InMemoryMapStore) delete(key string) error {
	if _, ok := s.store[key]; !ok {
		return fmt.Errorf("failed to delete key %q: %w", key, errNotFound)
	}
	delete(s.store, key)
	return nil
}
