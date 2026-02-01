package votes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cuego"
	"go.opentelemetry.io/otel/metric"
	"rmazur.io/poll/internal/telemetry"
)

// HTTPHandler provides a REST API handler for the votes data management.
func HTTPHandler(repo *Repository) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /votes", adapt(process(repo.Vote), voteCounter))
	mux.Handle("POST /labels", adapt(process(repo.Label), labelCounter))

	mux.Handle("GET /talk-data/{talk_id}", adapt(func(r *http.Request) (any, error) {
		return repo.Aggregate(r.PathValue("talk_id")), nil
	}, fetchCounter))
	return mux
}

// Metrics.
var (
	meter        = telemetry.Meter("votes")
	voteCounter  = must(meter.Int64Counter("operation.vote_total"))
	labelCounter = must(meter.Int64Counter("operation.label_total"))
	fetchCounter = must(meter.Int64Counter("operation.fetch_total"))
)

func process[T any](f func(T) error) handler {
	return func(r *http.Request) (any, error) {
		var data T
		if err := parse(r, &data); err != nil {
			log.Printf("parse error: %s", err)
			return nil, err
		}
		return nil, f(data)
	}
}

func parse(r *http.Request, out any) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		return &clientError{Msg: err.Error()}
	}
	if err := cuego.Validate(out); err != nil {
		return &clientError{Msg: fmt.Sprintf("%s", errors.Errors(err))}
	}
	return nil
}

type handler func(*http.Request) (any, error)

func adapt(f handler, counter metric.Int64Counter) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		counter.Add(context.Background(), 1)

		data, err := f(r)
		rw.Header().Add("content-type", "application/json")
		if err != nil {
			if ce, ok := err.(*clientError); ok {
				s := http.StatusBadRequest
				if ce.status != 0 {
					s = ce.status
				}
				rw.WriteHeader(s)
				_ = json.NewEncoder(rw).Encode(ce)
			} else {
				log.Printf("error processing the request: %s", err)
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		_ = json.NewEncoder(rw).Encode(data)
	})
}

type clientError struct {
	status int

	Msg string `json:"msg"`
}

func (ce *clientError) Error() string {
	return ce.Msg
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
