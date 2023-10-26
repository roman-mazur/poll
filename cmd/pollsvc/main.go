// Command svc can be used to start an HTTP server that exposes REST API for the poll data management.
// Its implementation is based on the rmazur.io/poll/votes package.
//
// Example usage:
//
//	./pollsvc --addr=:443 --tls=/opt/pollsvc-certs
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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"rmazur.io/poll/votes"
)

var (
	addr        = flag.String("addr", "127.0.0.1:17000", "Address to listen on")
	tlsPath     = flag.String("tls", "", "Path to the TLS cert.pem and pkey.pem files that should be used to configure the HTTP server")
	adminSecret = flag.String("admin-secret", "", "Admin secret to check for the config update")

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

	http.HandleFunc("/config/new", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("authorization") != *adminSecret {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		key := strings.TrimSpace(r.URL.Query().Get("key"))
		tc.Setup(key)
		log.Println("New talk", tc.CurrentId())

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(tc.CurrentId()))
	})

	http.Handle("/config/current", cors(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(tc.CurrentId()))
	})))

	http.Handle("/v1/", http.StripPrefix("/v1", cors(api)))
	http.Handle("/", http.FileServer(http.FS(appFS)))

	server := &http.Server{
		Addr:           *addr,
		Handler:        nil,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on %s", *addr)
	if *tlsPath != "" {
		log.Fatal(server.ListenAndServeTLS(filepath.Join(*tlsPath, "cert.pem"), filepath.Join(*tlsPath, "pkey.pem")))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(rw, r)
	})
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
