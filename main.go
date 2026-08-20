// main package to implements kvs server
package main

import (
	"flag"
	"fmt"
)

func printInitialMsg(port int) {
	str := `
    __  ___           __    _ 
   /  |/  /___  _____/ /_  (_)
  / /|_/ / __ \/ ___/ __ \/ / 
 / /  / / /_/ / /__/ / / / /  
/_/  /_/\____/\___/_/ /_/_/ 

Mochi server started: http://localhost:%d
`

	fmt.Printf(str, port)
}

func main() {
	port := flag.Int("port", 8080, "port to use")
	flag.Parse()

	store := InMemoryMapStore{}

	printInitialMsg(*port)
	s := newHTTPServer(store)
	s.serveHTTP(*port)
}
