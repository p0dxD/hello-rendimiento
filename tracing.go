package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// tracer names the app's own moments in its traces.
var tracer = otel.Tracer("hello-rendimiento")

// startTracing sends traces to rendimiento's collector, for the Universe
// page. With observe: true the platform sets OTEL_EXPORTER_OTLP_ENDPOINT,
// OTEL_SERVICE_NAME and the sampler; without them nothing is sent, and the
// spans the app makes cost next to nothing.
func startTracing(ctx context.Context) func(context.Context) error {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return func(context.Context) error { return nil }
	}
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return func(context.Context) error { return nil }
	}
	res, _ := resource.New(ctx, resource.WithFromEnv(), resource.WithTelemetrySDK())
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp.Shutdown
}

// traced times every request, named by its route ("GET /api/info").
func traced(h http.Handler) http.Handler {
	return otelhttp.NewHandler(h, "http", otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
		if r.Pattern != "" {
			return strings.TrimSuffix(r.Pattern, "{$}")
		}
		return r.Method
	}))
}
