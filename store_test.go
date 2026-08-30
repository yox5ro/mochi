package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func Test_ValidateKeyLength(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr error
	}{
		{
			name:    "empty key is valid",
			key:     "",
			wantErr: nil,
		},
		{
			name:    "max size key is valid",
			key:     strings.Repeat("a", MaxKeyLength),
			wantErr: nil,
		},
		{
			name:    "max size + 1 key is invalid",
			key:     strings.Repeat("a", MaxKeyLength+1),
			wantErr: errKeyTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := validateKeyLength(tt.key); !errors.Is(err, tt.wantErr) {
				t.Fatalf("want error %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func Test_ValidateValueLength(t *testing.T) {
	tests := []struct {
		name    string
		value   []byte
		wantErr error
	}{
		{
			name:    "empty value is valid",
			value:   []byte(""),
			wantErr: nil,
		},
		{
			name:    "max size value is valid",
			value:   bytes.Repeat([]byte{'a'}, MaxValueLength),
			wantErr: nil,
		},
		{
			name:    "max size + 1 value is invalid",
			value:   bytes.Repeat([]byte{'a'}, MaxValueLength+1),
			wantErr: errValueTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := validateValueLength(tt.value); !errors.Is(err, tt.wantErr) {
				t.Fatalf("want error %v, but got %v", tt.wantErr, err)
			}
		})
	}
}
