package main

import (
	"flag"
)

func main() {
	port := flag.String("port", ":8080", "port to use when protocol is HTTP")
	flag.Parse()

	store := InMemoryMapStore{}

	s := newHttpServer(store)
	s.serveHttp(*port)
}
