package observability

import (
	"context"
	"fmt"
	"havamal-api/config"
)

// Initialize sets up observability with logging and metrics
func Initialize(cfg config.Config) error {
	// Initialize logging first
	if err := InitializeLogging(cfg); err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}

	// Initialize metrics with the logger
	if err := InitializeMetrics(cfg, Logger); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down observability components
func Shutdown(ctx context.Context) error {
	var errs []error

	// Shutdown metrics
	if err := ShutdownMetrics(ctx); err != nil {
		errs = append(errs, fmt.Errorf("metrics shutdown: %w", err))
	}

	// Shutdown logging
	ShutdownLogging()

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}