// Package telemetry helps with observability for the poll applications.
package telemetry

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	om "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var (
	globalTeardown tearDownFunc

	buildInfo   *debug.BuildInfo
	buildInfoOk bool
)

func init() {
	buildInfo, buildInfoOk = debug.ReadBuildInfo()

	f, err := initializeOpenTelemetry()
	if err != nil {
		panic(err)
	}
	globalTeardown = f
}

// Shutdown performs a graceful teardown on the telemetry components initialized when this package is imported.
func Shutdown(ctx context.Context) error { return globalTeardown(ctx) }

// Meter constructs an Open Telemetry meter makng sure that the SDK is initialized first.
func Meter(name string) om.Meter { return otel.Meter(name) }

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

	attrs := []attribute.KeyValue{semconv.ServiceName("pollsvc")}
	if buildInfoOk {
		attrs = append(attrs, semconv.ServiceVersion(buildInfo.Main.Version))
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, attrs...),
	)
	if err != nil {
		return teardown, err
	}

	// Meter.
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return teardown, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(30*time.Second))),
	)
	allTearDown = append(allTearDown, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Emit runtime metrics.
	if err := runtime.Start(); err != nil {
		_ = teardown(context.Background())
		return teardown, err
	}

	return teardown, nil
}
