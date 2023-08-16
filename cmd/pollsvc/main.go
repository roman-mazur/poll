// Command svc can be used to start an HTTP server that exposes REST API for the poll data management.
// Its implementation is based on the rmazur.io/poll/votes package.
//
// Example usage:
//
//	./pollsvc --addr=localhost:8080
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"rmazur.io/poll/votes"
)

var (
	addr = flag.String("addr", "127.0.0.1:17000", "Address to listen on")

	//go:embed www
	www embed.FS
)

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), `%s

Start the poll application service.
Possible flags are below.

`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	repo := votes.NewRepository()
	api := votes.HTTPHandler(repo)

	appFS, err := fs.Sub(www, "www")
	if err != nil {
		panic(err)
	}

	http.Handle("/v1/", http.StripPrefix("/v1", api))
	http.Handle("/", http.FileServer(http.FS(appFS)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
