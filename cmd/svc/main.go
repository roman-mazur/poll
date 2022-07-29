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
