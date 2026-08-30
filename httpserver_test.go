package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestHTTPServer_BuildMux(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		wantStatusCode int
		wantEmptyBody  bool
		wantHeader     http.Header
	}{
		{
			name:           "root accepts GET with key query param",
			path:           "/?key=foo",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
			wantEmptyBody:  false,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "root accepts PUT with key query param",
			path:           "/?key=foo",
			method:         http.MethodPut,
			wantStatusCode: http.StatusNoContent,
			wantEmptyBody:  true,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "root accepts DELETE with key query param",
			path:           "/?key=foo",
			method:         http.MethodDelete,
			wantStatusCode: http.StatusNoContent,
			wantEmptyBody:  true,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "non-root returns 404",
			path:           "/foo",
			method:         http.MethodGet,
			wantStatusCode: http.StatusNotFound,
			wantEmptyBody:  true,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValuePathNotFound},
			},
		},
		{
			name:           "roots does not accept other than GET, PUT, DELETE",
			path:           "/",
			method:         http.MethodPost,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantEmptyBody:  true,
			wantHeader: map[string][]string{
				"Allow":                 {"GET, PUT, DELETE"},
				headerKeyMochiErrorCode: {headerValueRequestMethodInvalid},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			server := httpServer{store: newInMemoryMapStore(map[string][]byte{
				"foo": []byte("bar"),
			})}
			server.buildMux().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, but got %d", tt.wantStatusCode, rec.Code)
			}
			if (rec.Body.Len() != 0) == tt.wantEmptyBody {
				t.Fatalf("want empty response body, but got %s", rec.Body.String())
			}
			for k, v := range tt.wantHeader {
				actual := rec.Header().Values(k)
				if !slices.Equal(actual, v) {
					t.Fatalf("want Header %s to be %v, but got %v", k, v, actual)
				}
			}
		})
	}
}

var errInternal = errors.New("internal error")

type failingStore struct{}

func (s failingStore) get(_ string) ([]byte, error) {
	return nil, errInternal
}

func (s failingStore) put(_ string, _ []byte) error {
	return errInternal
}

func (s failingStore) delete(_ string) error {
	return errInternal
}

type storeType = int

const (
	success storeType = iota
	fail
)

func TestHTTPServer_HandleGetRequest(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		storeType      storeType
		wantStatusCode int
		wantBody       []byte
		wantHeader     http.Header
	}{
		{
			name:           "valid key returns 200 and response body",
			path:           "/?key=foo",
			storeType:      success,
			wantStatusCode: http.StatusOK,
			wantBody:       []byte("bar"),
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "invalid key returns 400 and empty response body",
			path:           "/?key=foo&key=bar",
			storeType:      success,
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueKeyInvalid},
			},
		},
		{
			name:           "nonexistent key returns 404 and empty body",
			path:           "/?key=not-existing",
			storeType:      success,
			wantStatusCode: http.StatusNotFound,
			wantBody:       nil,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueKeyNotFound},
			},
		},
		{
			name:           "internal server error returns 500 and empty response body",
			path:           "/?key=foo",
			storeType:      fail,
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       nil,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueInternal},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			successStore := newInMemoryMapStore(map[string][]byte{
				"foo": []byte("bar"),
			})
			failingStore := failingStore{}
			server := httpServer{store: successStore}
			if tt.storeType == fail {
				server.store = failingStore
			}

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			server.handleGetReq(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, but got %d", tt.wantStatusCode, rec.Code)
			}

			actualBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("error reading response body: %s", err.Error())
			}
			if !slices.Equal(actualBody, tt.wantBody) {
				t.Fatalf("want response body %s, but got %s", string(tt.wantBody), rec.Body.String())
			}

			for k, v := range tt.wantHeader {
				actual := rec.Header().Values(k)
				if !slices.Equal(actual, v) {
					t.Fatalf("want Header %s to be %v, but got %v", k, v, actual)
				}
			}
		})
	}
}

func TestHTTPServer_HandlePutRequest(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           []byte
		storeType      storeType
		wantStatusCode int
		wantHeader     http.Header
	}{
		{
			name:           "valid key returns 204",
			path:           "/?key=foo",
			body:           bytes.Repeat([]byte("a"), MaxValueSize),
			storeType:      success,
			wantStatusCode: http.StatusNoContent,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "valid key with empty body returns 204",
			path:           "/?key=foo",
			body:           nil,
			storeType:      success,
			wantStatusCode: http.StatusNoContent,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "invalid key returns 400",
			path:           "/?key=foo&key=bar",
			body:           []byte("test"),
			storeType:      success,
			wantStatusCode: http.StatusBadRequest,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueKeyInvalid},
			},
		},
		{
			name:           "too large body returns 413",
			path:           "/?key=foo",
			body:           bytes.Repeat([]byte("a"), MaxValueSize+1),
			storeType:      success,
			wantStatusCode: http.StatusRequestEntityTooLarge,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueValueTooLarge},
			},
		},
		{
			name:           "internal server error returns 500",
			path:           "/?key=foo",
			body:           []byte("test"),
			storeType:      fail,
			wantStatusCode: http.StatusInternalServerError,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueInternal},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			successStore := newInMemoryMapStore(map[string][]byte{
				"foo": []byte("bar"),
			})
			failingStore := failingStore{}
			server := httpServer{store: successStore}
			if tt.storeType == fail {
				server.store = failingStore
			}

			req := httptest.NewRequest(http.MethodPut, tt.path, bytes.NewReader(tt.body))
			rec := httptest.NewRecorder()

			server.handlePutReq(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, but got %d", tt.wantStatusCode, rec.Code)
			}

			for k, v := range tt.wantHeader {
				actual := rec.Header().Values(k)
				if !slices.Equal(actual, v) {
					t.Fatalf("want Header %s to be %v, but got %v", k, v, actual)
				}
			}
		})
	}
}

func TestHTTPServer_HandleDeleteRequest(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		storeType      storeType
		wantStatusCode int
		wantHeader     http.Header
	}{
		{
			name:           "valid key returns 204",
			path:           "/?key=foo",
			storeType:      success,
			wantStatusCode: http.StatusNoContent,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "nonexistent key returns 204",
			path:           "/?key=not-existing",
			storeType:      success,
			wantStatusCode: http.StatusNoContent,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: nil,
			},
		},
		{
			name:           "invalid key returns 400",
			path:           "/?key=foo&key=bar",
			storeType:      success,
			wantStatusCode: http.StatusBadRequest,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueKeyInvalid},
			},
		},
		{
			name:           "internal server error returns 500",
			path:           "/?key=foo",
			storeType:      fail,
			wantStatusCode: http.StatusInternalServerError,
			wantHeader: map[string][]string{
				headerKeyMochiErrorCode: {headerValueInternal},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			successStore := newInMemoryMapStore(map[string][]byte{
				"foo": []byte("bar"),
			})
			failingStore := failingStore{}
			server := httpServer{store: successStore}
			if tt.storeType == fail {
				server.store = failingStore
			}

			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			rec := httptest.NewRecorder()

			server.handleDeleteReq(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("want status %d, but got %d", tt.wantStatusCode, rec.Code)
			}

			for k, v := range tt.wantHeader {
				actual := rec.Header().Values(k)
				if !slices.Equal(actual, v) {
					t.Fatalf("want Header %s to be %v, but got %v", k, v, actual)
				}
			}
		})
	}
}

func Test_ExtractKey(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr error
	}{
		{
			name: "single key can be extracted",
			path: "/?key=foo",
			want: "foo",
		},
		{
			name: "empty string is valid",
			path: "/?key=",
			want: "",
		},
		{
			name:    "multiple key is invalid",
			path:    "/?key=foo&key=bar",
			want:    "",
			wantErr: errKeyInvalid,
		},
		{
			name:    "invalid percent encoding is invalid",
			path:    "/?key=foo&key=%zz",
			want:    "",
			wantErr: errKeyInvalid,
		},
		{
			name:    "no query param is invalid",
			path:    "/",
			want:    "",
			wantErr: errKeyInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			actual, err := extractKey(req)

			if actual != tt.want {
				t.Fatalf("want return value %s, but got %s", tt.want, actual)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want error %v, but got %v", tt.wantErr, err)
			}
		})
	}
}
