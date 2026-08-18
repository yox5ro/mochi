// main package to implements kvs server
package main

import (
	"flag"
)

func main() {
	port := flag.String("port", ":8080", "port to use when protocol is HTTP")
	flag.Parse()

	store := InMemoryMapStore{}

	s := newHTTPServer(store)
	s.serveHTTP(*port)
}
