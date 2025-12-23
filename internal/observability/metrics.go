package observability

import (
	"context"
	"fmt"
	"havamal-api/config"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	meterProvider *metric.MeterProvider
	logger        *slog.Logger
)

// InitializeMetrics sets up Prometheus metrics with optional remote write
func InitializeMetrics(cfg config.Config, log *slog.Logger) error {
	logger = log

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.Observability.ServiceName),
			semconv.DeploymentEnvironment(cfg.Observability.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create local Prometheus exporter for /metrics endpoint
	promExporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	// Create meter provider
	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExporter),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	logger.Info("Prometheus metrics initialized (local only)",
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
	)

	return nil
}

// ShutdownMetrics gracefully shuts down metrics components
func ShutdownMetrics(ctx context.Context) error {


	if meterProvider != nil {
		logger.Info("Shutting down meter provider")
		if err := meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("meter provider shutdown: %w", err)
		}
	}

	return nil
}


