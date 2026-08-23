// main package to implements kvs server
package main

import (
	"flag"
	"fmt"
	"log"
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

	store := newInMemoryMapStore(make(map[string]string))

	printInitialMsg(*port)
	s := httpServer{store: store}
	if err := s.serveHTTP(*port); err != nil {
		log.Fatal(err)
	}
}
