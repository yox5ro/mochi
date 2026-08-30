package main

import (
	"errors"
	"slices"
	"testing"
)

func TestInMemoryMapStore_Get(t *testing.T) {
	tests := []struct {
		name         string
		initialState map[string][]byte
		key          string
		want         []byte
		wantErr      error
	}{
		{
			name: "can get by valid key",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "hoge",
			want:    []byte("value"),
			wantErr: nil,
		},
		{
			name: "returns nil when no key found",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "foo",
			want:    nil,
			wantErr: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newInMemoryMapStore(tt.initialState)
			actual, err := s.get(tt.key)

			if !slices.Equal(actual, tt.want) {
				t.Errorf("wanted %s, but got %s", tt.want, actual)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted error %s, but got %s", tt.wantErr, err)
			}
		})
	}
}

func TestInMemoryMapStore_Put(t *testing.T) {
	tests := []struct {
		name         string
		initialState map[string][]byte
		key          string
		value        []byte
		wantErr      error
	}{
		{
			name: "can create new key",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "foo",
			value:   []byte("bar"),
			wantErr: nil,
		},
		{
			name: "can update existing key",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "hoge",
			value:   []byte("new value"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newInMemoryMapStore(tt.initialState)
			err := s.put(tt.key, tt.value)
			data, _ := s.get(tt.key)

			if !slices.Equal(data, tt.value) {
				t.Errorf("wanted %s, but got %s", tt.value, data)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted error %s, but got %s", tt.wantErr, err)
			}
		})
	}
}

func TestInMemoryMapStore_Delete(t *testing.T) {
	tests := []struct {
		name         string
		initialState map[string][]byte
		key          string
		wantErr      error
	}{
		{
			name: "can delete by valid key",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "hoge",
			wantErr: nil,
		},
		{
			name: "returns empty string when no key found",
			initialState: map[string][]byte{
				"hoge": []byte("value"),
				"fuga": []byte("value2"),
			},
			key:     "foo",
			wantErr: errNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newInMemoryMapStore(tt.initialState)
			err := s.delete(tt.key)

			if data, err := s.get(tt.key); !errors.Is(err, errNotFound) {
				t.Errorf("data still found on %s", data)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted error %s, but got %s", tt.wantErr, err)
			}
		})
	}
}

func TestInMemoryMapStore_ValueOwnership(t *testing.T) {
	t.Run("initialState change does not affect store", func(t *testing.T) {
		const initialValueSting = "initial value"
		initialValue := []byte(initialValueSting)
		initialState := map[string][]byte{
			"key": initialValue,
		}

		s := newInMemoryMapStore(initialState)

		initialState["key"][0] = 'X'

		v, err := s.get("key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slices.Equal(v, []byte(initialValueSting)) {
			t.Fatalf("want value %s, but got %s", initialValueSting, v)
		}
	})

	t.Run("get result change does not affect store", func(t *testing.T) {
		const initialValueString = "initial value"
		initialValue := []byte(initialValueString)
		s := newInMemoryMapStore(map[string][]byte{
			"key": initialValue,
		})

		first, err := s.get("key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		first[0] = 'T'

		second, err := s.get("key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !slices.Equal([]byte(initialValueString), second) {
			t.Fatalf("value changed, want %s, but got %s", initialValueString, second)
		}
	})

	t.Run("change after put does not affect store", func(t *testing.T) {
		s := newInMemoryMapStore(make(map[string][]byte))

		const putValueString = "put value"
		putValue := []byte(putValueString)

		if err := s.put("key", putValue); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		putValue[0] = 'T'

		got, err := s.get("key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !slices.Equal(got, []byte(putValueString)) {
			t.Fatalf("value changed, want %s, but got %s", putValueString, got)
		}
	})
}
