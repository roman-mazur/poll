// Command svc can be used to start an HTTP server that exposes REST API for the poll data management.
// Its implementation is based on the rmazur.io/poll/votes package.
//
// Example usage:
//
//	./pollsvc --addr=:443 --tls=/opt/pollsvc-certs
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"rmazur.io/poll/internal/telemetry"
	"rmazur.io/poll/votes"
)

var (
	addr        = flag.String("addr", "127.0.0.1:17000", "Address to listen on")
	tlsPath     = flag.String("tls", "", "Path to the TLS cert.pem and pkey.pem files that should be used to configure the HTTP server")
	adminSecret = flag.String("admin-secret", "", "Admin secret to check for the config update")
)

func main() {
	shutdownCh := shutdownSignal()
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

	server := &http.Server{
		Addr:           *addr,
		Handler:        buildMux(api),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serverClosedCh := make(chan error)

	defer callShutdown("telemetry", telemetry.Shutdown)

	initReady := make(chan struct{})
	go func() {
		var err error
		log.Printf("Listening on %s", *addr)
		close(initReady)
		if *tlsPath != "" {
			err = server.ListenAndServeTLS(filepath.Join(*tlsPath, "cert.pem"), filepath.Join(*tlsPath, "pkey.pem"))
		} else {
			err = server.ListenAndServe()
		}
		serverClosedCh <- err
	}()

	<-initReady
	log.Println("Server initialized")
	select {
	case <-shutdownCh:
		log.Println("Shutdown requested")
		callShutdown("http", server.Shutdown)
	case <-serverClosedCh:
	}
}

func callShutdown(name string, f func(context.Context) error) {
	teardownCtx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	if err := f(teardownCtx); err != nil {
		log.Printf("Error shutting down %s: %s", name, err)
	}
	log.Printf("%s shutdown complete", name)
}

func shutdownSignal() chan os.Signal {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, shutdownSignals...)
	return signalCh
}
