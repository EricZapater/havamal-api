package main

import (
	"flag"
	"havamal-api/config"
	"havamal-api/internal/db"
	"havamal-api/internal/migrations"
	"log/slog"
	"os"
)

func main() {
	action := flag.String("action", "up", "migration action: up or down")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Unable to load config", slog.Any("error", err))
		os.Exit(1)
	}

	database, err := db.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("Unable to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer database.Close()

	switch *action {
	case "reset":
		slog.Info("Resetting database...")
		if _, err := database.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
			slog.Error("Failed to reset database", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Database reset. Running migrations...")
		if err := migrations.RunMigrations(database, cfg.Migration.Path); err != nil {
			slog.Error("Failed to run migrations", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Reset and migrations complete")
	case "down":
		slog.Info("Rolling back migrations...")
		if err := migrations.RollbackMigration(database, cfg.Migration.Path); err != nil {
			slog.Error("Failed to rollback migrations", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Rollback complete")
	default:
		slog.Info("Running migrations...")
		if err := migrations.RunMigrations(database, cfg.Migration.Path); err != nil {
			slog.Error("Failed to run migrations", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Migrations complete")
	}
}
