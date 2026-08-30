package main

import "testing"

func mustNewInMemoryMapStore(t *testing.T, initialState map[string][]byte) InMemoryMapStore {
	t.Helper()

	store, err := newInMemoryMapStore(initialState)
	if err != nil {
		t.Fatalf("failed to create inMemoryMapStore: %v", err)
	}
	return store
}
