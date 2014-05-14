package main

import (
	"github.com/st3fan/moz-csp/csp"
	"log"
)

func main() {
	session, err := csp.NewSession("postgres://csp:csp@localhost/csp")
	if err != nil {
		log.Fatalf("Can't open database session: %s", err)
	}
	defer session.Close()

	server := csp.NewServer(session)
	server.Run()
}
