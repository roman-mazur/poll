package votes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"

	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cuego"
)

func HTTPHandler(repo *Repository) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/votes", adapt(method("POST", process(repo.Add))))
	mux.Handle("/labels", adapt(method("POST", process(repo.Label))))

	mux.Handle("/talk-data/", adapt(method("GET", func(r *http.Request) (any, error) {
		talkId := path.Base(r.URL.Path)
		if talkId == "" {
			return nil, &clientError{Msg: "no talk ID"}
		}
		return repo.Aggregate(talkId), nil
	})))
	return mux
}

type touch interface {
	touch()
}

func process[T touch](f func(T) error) handler {
	return func(r *http.Request) (any, error) {
		var data T
		if err := parse(r, &data); err != nil {
			log.Printf("parse error: %s", err)
			return nil, err
		}
		data.touch()
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

func adapt(f handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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

func method(m string, f handler) handler {
	return func(r *http.Request) (any, error) {
		if r.Method != m {
			return nil, &clientError{status: http.StatusMethodNotAllowed, Msg: "wrong verb"}
		}
		return f(r)
	}
}

type clientError struct {
	status int
	Msg    string `json:"msg"`
}

func (ce *clientError) Error() string {
	return ce.Msg
}
