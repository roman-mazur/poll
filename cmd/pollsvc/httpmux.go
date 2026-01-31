package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
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

	buildInfo, buildInfoOk := debug.ReadBuildInfo()
	httpMux.Handle("GET /ping", relaxCORS(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		var data struct {
			Version string `json:"version"`
		}
		if buildInfoOk {
			data.Version = buildInfo.Main.Version
		}
		_ = json.NewEncoder(rw).Encode(&data)
	})))

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
