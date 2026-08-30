package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
)

const (
	headerKeyMochiErrorCode = "Mochi-Error-Code"

	headerValueKeyInvalid           = "key-invalid"
	headerValueKeyNotFound          = "key-not-found"
	headerValuePathNotFound         = "path-not-found"
	headerValueRequestMethodInvalid = "request-method-invalid"
	headerValueKeyTooLarge          = "key-too-large"
	headerValueValueTooLarge        = "value-too-large"
	headerValueInternal             = "internal"
)

var errKeyInvalid = errors.New("key invalid")

type httpServer struct {
	store Store
}

func (s httpServer) serveHTTP(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), s.buildMux())
}

func (s httpServer) buildMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleGetReq)
	mux.HandleFunc("PUT /{$}", s.handlePutReq)
	mux.HandleFunc("DELETE /{$}", s.handleDeleteReq)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.Header().Set(headerKeyMochiErrorCode, headerValuePathNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !slices.Contains([]string{http.MethodGet, http.MethodPut, http.MethodDelete}, r.Method) {
			w.Header().Set("Allow", "GET, PUT, DELETE")
			w.Header().Set(headerKeyMochiErrorCode, headerValueRequestMethodInvalid)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (s httpServer) handleGetReq(w http.ResponseWriter, r *http.Request) {
	key, err := extractKey(r)
	if err != nil {
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyInvalid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := s.store.get(key)
	switch {
	case err == nil:
		if _, writeErr := w.Write(value); writeErr != nil {
			// TODO log the error
			fmt.Println("response write failed")
		}
		return
	case errors.Is(err, errNotFound):
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	case errors.Is(err, errKeyTooLarge):
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyTooLarge)
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		w.Header().Set(headerKeyMochiErrorCode, headerValueInternal)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s httpServer) handlePutReq(w http.ResponseWriter, r *http.Request) {
	key, err := extractKey(r)
	if err != nil {
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyInvalid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := io.LimitReader(r.Body, MaxValueLength+1)
	value, err := io.ReadAll(body)
	if err != nil {
		w.Header().Set(headerKeyMochiErrorCode, headerValueInternal)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.store.put(key, value)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
		return
	case errors.Is(err, errKeyTooLarge):
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyTooLarge)
		w.WriteHeader(http.StatusBadRequest)
		return
	case errors.Is(err, errValueTooLarge):
		w.Header().Set(headerKeyMochiErrorCode, headerValueValueTooLarge)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	default:
		w.Header().Set(headerKeyMochiErrorCode, headerValueInternal)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s httpServer) handleDeleteReq(w http.ResponseWriter, r *http.Request) {
	key, err := extractKey(r)
	if err != nil {
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyInvalid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.store.delete(key)
	switch {
	case err == nil || errors.Is(err, errNotFound):
		w.WriteHeader(http.StatusNoContent)
		return
	case errors.Is(err, errKeyTooLarge):
		w.Header().Set(headerKeyMochiErrorCode, headerValueKeyTooLarge)
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		w.Header().Set(headerKeyMochiErrorCode, headerValueInternal)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func extractKey(r *http.Request) (string, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if values, ok := query["key"]; !ok || len(values) != 1 || err != nil {
		return "", errKeyInvalid
	}
	return query["key"][0], nil
}
