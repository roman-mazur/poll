// Package telemetry helps with observability for the poll applications.
package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var globalTeardown tearDownFunc

func init() {
	f, err := initializeOpenTelemetry()
	if err != nil {
		panic(err)
	}
	globalTeardown = f
}

// Shutdown performs a graceful teardown on the telemetry components initialized when this package is imported.
func Shutdown(ctx context.Context) error { return globalTeardown(ctx) }

type tearDownFunc func(context.Context) error

func initializeOpenTelemetry() (teardown tearDownFunc, err error) {
	var allTearDown []tearDownFunc
	teardown = func(ctx context.Context) error {
		var allErrors []error
		for _, f := range allTearDown {
			if err := f(ctx); err != nil {
				allErrors = append(allErrors, err)
			}
		}
		return errors.Join(allErrors...)
	}

	// Meter.
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return teardown, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(30*time.Second))),
	)
	allTearDown = append(allTearDown, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return teardown, nil
}
