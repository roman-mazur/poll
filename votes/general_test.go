package votes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func submit(url string, data any) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(js))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", resp.StatusCode)
	}
	return nil
}

func TestHTTPHandler(t *testing.T) {
	repo := NewRepository()
	h := HTTPHandler(repo)

	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	ts, err := time.Parse(time.RFC3339, "2022-02-24T04:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("submit data", func(t *testing.T) {
		vote := Vote{
			TalkName:  "test-talk",
			Timestamp: ts,
			VoterId:   "some-voter-id",
			Value:     10,
		}
		if err := submit(srv.URL+"/votes", vote); err != nil {
			t.Error(err)
		}
		label := Label{
			TalkName:  "test-talk",
			Timestamp: ts,
			Name:      "some-slide-1",
		}
		if err := submit(srv.URL+"/labels", label); err != nil {
			t.Error(err)
		}

		agg := repo.Aggregate("test-talk")
		if len(agg.Labels) == 0 {
			t.Error("no labels stored")
		}
		if len(agg.Votes) == 0 {
			t.Error("no votes stored")
		}
	})

}
