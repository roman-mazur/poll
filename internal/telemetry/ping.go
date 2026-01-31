package telemetry

import (
	"encoding/json"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	internalMeter = sync.OnceValue(func() metric.Meter { return otel.Meter("internal/telemetry") })

	pingCounter = sync.OnceValue(func() metric.Int64Counter {
		res, err := internalMeter().Int64Counter("ping")
		if err != nil {
			panic(err)
		}
		return res
	})
)

// Ping constructs an HTTP handler that can be used for a liveness check returning the version info.
func Ping() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		pingCounter().Add(r.Context(), 1)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		var data struct {
			Version string `json:"version"`
		}
		if buildInfoOk {
			data.Version = buildInfo.Main.Version
		}
		_ = json.NewEncoder(rw).Encode(&data)
	})
}
