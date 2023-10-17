// Command svc can be used to start an HTTP server that exposes REST API for the poll data management.
// Its implementation is based on the rmazur.io/poll/votes package.
//
// Example usage:
//
//	./pollsvc --addr=localhost:8080
package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

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

	var tc talkConfig

	http.HandleFunc("/config/new/", func(rw http.ResponseWriter, r *http.Request) {
		tc.Setup(path.Base(r.URL.Path))
		rw.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/config/current", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(tc.CurrentId()))
	})

	http.Handle("/v1/", http.StripPrefix("/v1", api))
	http.Handle("/", http.FileServer(http.FS(appFS)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type talkConfig struct {
	mu     sync.Mutex
	talkId string
}

func (tc *talkConfig) Setup(name string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	suffix := sha256.Sum256([]byte(time.Now().String()))
	tc.talkId = name + "-" + hex.EncodeToString(suffix[:])
}

func (tc *talkConfig) CurrentId() string {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.talkId
}
