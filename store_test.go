package main

import (
	"errors"
	"testing"
)

func TestInMemoryMapStore_Get(t *testing.T) {
	tests := []struct {
		name         string
		initialState map[string]string
		key          string
		want         string
		wantErr      error
	}{
		{
			name: "can get by valid key",
			initialState: map[string]string{
				"hoge": "key",
			},
			key:     "hoge",
			want:    "key",
			wantErr: nil,
		},
		{
			name: "returns empty string when no key found",
			initialState: map[string]string{
				"hoge": "key",
			},
			key:     "",
			want:    "",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newInMemoryMapStore(tt.initialState)
			actual, err := s.get(tt.key)

			if actual != tt.want {
				t.Errorf("wanted %s, but got %s", tt.want, actual)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted error %s, but got %s", tt.wantErr, err)
			}
		})
	}
}
