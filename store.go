package main

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
	inMemoryMapStore.store = initialState
	return inMemoryMapStore
}

func (s InMemoryMapStore) get(key string) (string, error) {
	return s.store[key], nil
}

func (s InMemoryMapStore) put(key, value string) error {
	s.store[key] = value
	return nil
}

func (s InMemoryMapStore) delete(key string) error {
	delete(s.store, key)
	return nil
}
