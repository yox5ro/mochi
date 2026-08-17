package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type request struct {
	Op    string `json:"op"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// one of Value or Error_ will be omitted
type response struct {
	Value  string `json:"value,omitempty"`
	Error_ string `json:"error,omitempty"`
}

type httpServer struct {
	store Store
}

func newHttpServer(store Store) httpServer {
	return httpServer{store: store}
}

func (s httpServer) serveHttp(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", s.handleReq)
	log.Fatal(http.ListenAndServe(port, mux))
}

func (s httpServer) handleReq(w http.ResponseWriter, r *http.Request) {
	var value string
	var err error = nil
	defer func() {
		resp := response{}
		if err != nil {
			resp.Error_ = err.Error()
		} else {
			resp.Value = value
		}
		json.NewEncoder(w).Encode(resp)
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
