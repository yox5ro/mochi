package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type request struct {
	Op    string `json:"op"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// one of Value or Error_ will be omitted
type response struct {
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

type httpServer struct {
	store Store
}

func (s httpServer) buildMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", s.handleReq)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (s httpServer) serveHTTP(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), s.buildMux())
}

func (s httpServer) handleReq(w http.ResponseWriter, r *http.Request) {
	var value string
	var err error
	defer func() {
		resp := response{}
		if err != nil {
			resp.Error = err.Error()
		} else {
			resp.Value = value
		}
		if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
			log.Printf("failed to encode response: %v", resp)
		}
	}()

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errors.New("error: invalid request")
		return
	}

	switch req.Op {
	case "get":
		value, err = s.store.get(req.Key)
	case "put":
		err = s.store.put(req.Key, req.Value)
	case "delete":
		err = s.store.delete(req.Key)
	default:
		err = errors.New("error: invalid request")
	}
}
