package main

import (
	"fmt"
	"slices"
)

type InMemoryMapStore struct {
	store map[string][]byte
}

func newInMemoryMapStore(initialState map[string][]byte) (InMemoryMapStore, error) {
	inMemoryMapStore := InMemoryMapStore{store: make(map[string][]byte)}

	for k, v := range initialState {
		if err := validateKeyLength(k); err != nil {
			return InMemoryMapStore{}, err
		}
		if err := validateValueLength(v); err != nil {
			return InMemoryMapStore{}, err
		}
		inMemoryMapStore.store[k] = slices.Clone(v)
	}

	return inMemoryMapStore, nil
}

func (s InMemoryMapStore) get(key string) ([]byte, error) {
	if err := validateKeyLength(key); err != nil {
		return nil, err
	}
	if v, ok := s.store[key]; ok {
		return slices.Clone(v), nil
	}
	return nil, fmt.Errorf("failed to get key %q: %w", key, errNotFound)
}

func (s InMemoryMapStore) put(key string, value []byte) error {
	if err := validateKeyLength(key); err != nil {
		return err
	}
	if err := validateValueLength(value); err != nil {
		return err
	}
	s.store[key] = slices.Clone(value)
	return nil
}

func (s InMemoryMapStore) delete(key string) error {
	if err := validateKeyLength(key); err != nil {
		return err
	}
	if _, ok := s.store[key]; !ok {
		return fmt.Errorf("failed to delete key %q: %w", key, errNotFound)
	}
	delete(s.store, key)
	return nil
}
