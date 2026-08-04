package telemetry

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracer configures OpenTelemetry from standard OTEL_* environment variables.
// When OTEL_TRACES_EXPORTER is empty or "none", a no-op tracer is used.
func InitTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exporter := os.Getenv("OTEL_TRACES_EXPORTER")
	if exporter == "" || exporter == "none" {
		tp := sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		return tp, nil
	}

	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "cv-backend"
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// Shutdown flushes and closes the tracer provider.
func Shutdown(ctx context.Context, tp *sdktrace.TracerProvider) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}

// TraceFilter skips high-cardinality infrastructure endpoints in Jaeger.
func TraceFilter(r *http.Request) bool {
	switch r.URL.Path {
	case "/metrics", "/healthz":
		return false
	default:
		return true
	}
}
