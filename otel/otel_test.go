package otel_pkg

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// mockExporter is a no-op SpanExporter for testing purposes.
type mockExporter struct{}

// ExportSpans implements sdktrace.SpanExporter.ExportSpans.
//
// It is a no-op, and always returns nil.
func (m *mockExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	return nil
}

// Shutdown does nothing and returns no error.
func (m *mockExporter) Shutdown(_ context.Context) error {
	return nil
}

// fakeExporterFactory is a fake ExporterFactory for testing purposes, which
// returns a no-op SpanExporter.
func fakeExporterFactory(_ context.Context) (sdktrace.SpanExporter, error) {
	return &mockExporter{}, nil
}

// TestInitTracerProvider_WithMockExporter tests that initTracerProviderWithExporter returns
// a valid shutdown function that can be used to shut down the TracerProvider.
//
// It also checks that the shutdown function returns no error when called.
func TestInitTracerProvider_WithMockExporter(t *testing.T) {
	ctx := context.Background()

	shutdown, err := initTracerProviderWithExporter("test-service", ctx, fakeExporterFactory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected shutdown function, got nil")
	}

	if err := shutdown(ctx); err != nil {
		t.Errorf("unexpected error on shutdown: %v", err)
	}
}
