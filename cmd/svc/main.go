// Command svc can be used to start an HTTP server that exposes REST API for the poll data management.
// Its implementation is based on the rmazur.io/poll/votes package.
//
// Example usage:
//		./svc --addr=localhost:8080
package main

import (
	"flag"
	"log"
	"net/http"

	"rmazur.io/poll/votes"
)

var (
	addr = flag.String("addr", "127.0.0.1:17000", "Address to listen on")
)

func main() {
	flag.Parse()

	repo := votes.NewRepository()
	api := votes.HTTPHandler(repo)

	http.Handle("/v1/", http.StripPrefix("/v1", api))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
