package main

import (
	"context"
	"havamal-api/config"
	"havamal-api/internal/db"
	"havamal-api/internal/migrations"
	"havamal-api/internal/observability"
	"havamal-api/server"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Unable to load config: ", err)
	}

	// Initialize observability
	if err := observability.Initialize(cfg); err != nil {
		log.Fatal("Unable to initialize observability: ", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := observability.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down observability", slog.Any("error", err))
		}
	}()

	slog.Info("Starting Havamal API", slog.String("version", "1.0.0"))

	database, err := db.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("Unable to connect to database", slog.Any("error", err))
		log.Fatal(err)
	}

	// Run database migrations
	slog.Info("Running database migrations...")
	if err := migrations.RunMigrations(database, cfg.Migration.Path); err != nil {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		log.Fatal(err)
	}

	server := server.NewServer(cfg, database)
	if err := server.Setup(); err != nil {
		slog.Error("Unable to setup server", slog.Any("error", err))
		log.Fatal(err)
	}

	// Start server in a goroutine
	go func() {
		if err := server.Run(); err != nil {
			slog.Error("Server error", slog.Any("error", err))
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
}