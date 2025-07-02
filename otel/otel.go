package otel_pkg

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// ExporterFactory is a function that creates a trace exporter.
type ExporterFactory func(ctx context.Context) (sdktrace.SpanExporter, error)

// DefaultExporterFactory creates an OTLP GRPC trace exporter based on the
// environment variable OTEL_ENDPOINT. The exporter is configured to be
// insecure, meaning it does not validate the server's certificate.
//
// This function should be used as a fallback when no other exporter factory
// is specified.
func DefaultExporterFactory(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(viper.GetString("OTEL_ENDPOINT")),
	)
}

// InitTracerProvider initializes an OpenTelemetry TracerProvider for a given service.
// It sets up a trace exporter using the OTLP gRPC protocol and configures the TracerProvider
// with resource attributes such as service name, version, and environment. It also sets up
// a batch span processor and a composite text map propagator for context propagation.
//
// Parameters:
//   - service: The name of the service for which the TracerProvider is being initialized.
//   - ctx: The context for managing the lifecycle of the TracerProvider and trace exporter.
//
// Returns:
//   - A function to shut down the TracerProvider, releasing any resources held.
//   - An error if there was a failure during the initialization of the trace exporter or resource.
func InitTracerProvider(service string, ctx context.Context) (func(context.Context) error, error) {
	return initTracerProviderWithExporter(service, ctx, DefaultExporterFactory)
}

// initTracerProviderWithExporter initializes an OpenTelemetry TracerProvider for a given service using the given trace exporter factory.
// It sets up a trace exporter using the OTLP gRPC protocol and configures the TracerProvider
// with resource attributes such as service name, version, and environment. It also sets up
// a batch span processor and a composite text map propagator for context propagation.
//
// Parameters:
//   - service: The name of the service for which the TracerProvider is being initialized.
//   - ctx: The context for managing the lifecycle of the TracerProvider and trace exporter.
//   - factory: A function that creates a trace exporter.
//
// Returns:
//   - A function to shut down the TracerProvider, releasing any resources held.
//   - An error if there was a failure during the initialization of the trace exporter or resource.
func initTracerProviderWithExporter(service string, ctx context.Context, factory ExporterFactory) (func(context.Context) error, error) {
	traceExporter, err := factory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := sdkresource.New(
		ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
