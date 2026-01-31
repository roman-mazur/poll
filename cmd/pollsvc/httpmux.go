package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"rmazur.io/poll/internal/telemetry"
)

//go:embed www
var www embed.FS

func buildMux(votesRestApi http.Handler) http.Handler {
	appFS, err := fs.Sub(www, "www")
	if err != nil {
		panic(err)
	}

	httpMux := http.NewServeMux()

	var tc talkConfig

	httpMux.HandleFunc("POST /config/new", func(rw http.ResponseWriter, r *http.Request) {
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

	httpMux.Handle("GET /config/current", relaxCORS(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte(tc.CurrentId()))
	})))

	httpMux.Handle("GET /ping", relaxCORS(telemetry.Ping()))

	httpMux.Handle("/v1/", http.StripPrefix("/v1", relaxCORS(votesRestApi)))
	httpMux.Handle("/", http.FileServer(http.FS(appFS)))

	return httpMux
}

func relaxCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(rw, r)
	})
}
