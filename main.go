package main

import (
	"flag"
	"github.com/st3fan/moz-csp-collector/csp"
	"log"
)

// moz-csp-collector -db <dburl>

var database = flag.String("database", "postgres://csp:csp@localhost/csp", "the database to connect to")

func main() {
	flag.Parse()
	session, err := csp.NewSession(*database)
	if err != nil {
		log.Fatalf("Can't open database session: %s", err)
	}
	defer session.Close()

	server := csp.NewServer(session)
	server.Run()
}
