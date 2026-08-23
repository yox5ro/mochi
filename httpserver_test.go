package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestHTTPServer_BuildMux(t *testing.T) {
	server := httpServer{store: newInMemoryMapStore(nil)}

	tests := []struct {
		name           string
		path           string
		method         string
		wantStatusCode int
		wantEmptyBody  bool
		wantHeader     http.Header
	}{
		{
			name:           "root accepts POST",
			path:           "/",
			method:         http.MethodPost,
			wantStatusCode: http.StatusOK,
			wantEmptyBody:  false,
		},
		{
			name:           "non-root returns 404",
			path:           "/foo",
			method:         http.MethodPost,
			wantStatusCode: http.StatusNotFound,
			wantEmptyBody:  true,
		},
		{
			name:           "roots does not accept other than POST",
			path:           "/",
			method:         http.MethodGet,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantEmptyBody:  true,
			wantHeader:     map[string][]string{"Allow": []string{"POST"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

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
