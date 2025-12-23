package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Auth struct {
		Secret string
		TTL    time.Duration
	}
	Email struct {
		Host     string
		Port     string
		User     string
		Password string
		From     string
	}
	Migration struct {
		Path string
	}
	Observability struct {
		// Loki (logs)
		LokiURL      string
		LokiUser     string
		LokiPassword string
		
		// Prometheus (metrics)
		PrometheusURL      string
		PrometheusUser     string
		PrometheusPassword string
		
		// Common
		ServiceName string
		Environment string
	}
}

func LoadConfig() (Config, error) {
	// Carrega el .env
	_ = godotenv.Load()
	
	var cfg Config
	cfg.App.Port = getenvDefault("APP_PORT", "8080")
	
	// Database config...
	cfg.Database.Host = getenvDefault("DATABASE_HOST", "localhost")
	if cfg.Database.Host == "" {
		return Config{}, errors.New("DATABASE_HOST is required")
	}
	cfg.Database.Port = getenvDefault("DATABASE_PORT", "5432")
	if cfg.Database.Port == "" {
		return Config{}, errors.New("DATABASE_PORT is required")
	}
	cfg.Database.User = getenvDefault("DATABASE_USER", "postgres")
	if cfg.Database.User == "" {
		return Config{}, errors.New("DATABASE_USER is required")
	}
	cfg.Database.Password = getenvDefault("DATABASE_PASSWORD", "postgres")
	if cfg.Database.Password == "" {
		return Config{}, errors.New("DATABASE_PASSWORD is required")
	}
	cfg.Database.Name = getenvDefault("DATABASE_NAME", "postgres")
	if cfg.Database.Name == "" {
		return Config{}, errors.New("DATABASE_NAME is required")
	}
	
	// Auth config...
	cfg.Auth.Secret = getenvDefault("AUTH_SECRET", "secret")
	if cfg.Auth.Secret == "" {
		return Config{}, errors.New("AUTH_SECRET is required")
	}
	ttlString := getenvDefault("AUTH_TTL", "604800")
	ttlSeconds, err := strconv.Atoi(ttlString)
	if err != nil {
		return Config{}, errors.New("AUTH_TTL must be an integer representing seconds")
	}
	cfg.Auth.TTL = time.Duration(ttlSeconds) * time.Second
	
	// Email config...
	cfg.Email.Host = getenvDefault("EMAIL_HOST", "localhost")
	cfg.Email.Port = getenvDefault("EMAIL_PORT", "1025")
	cfg.Email.User = getenvDefault("EMAIL_USER", "")
	cfg.Email.Password = getenvDefault("EMAIL_PASSWORD", "")
	cfg.Email.From = getenvDefault("EMAIL_FROM", "")
	
	// Migration config...
	cfg.Migration.Path = getenvDefault("MIGRATION_PATH", "./migrations")
	
	// Observability configuration
	cfg.Observability.LokiURL = getenvDefault("LOKI_URL", "")
	cfg.Observability.LokiUser = getenvDefault("LOKI_USER", "")
	cfg.Observability.LokiPassword = getenvDefault("LOKI_PASSWORD", "")
	
	cfg.Observability.PrometheusURL = getenvDefault("PROMETHEUS_URL", "")
	cfg.Observability.PrometheusUser = getenvDefault("PROMETHEUS_USER", "")
	cfg.Observability.PrometheusPassword = getenvDefault("PROMETHEUS_PASSWORD", "")
	
	cfg.Observability.ServiceName = getenvDefault("SERVICE_NAME", "turniq-api")
	cfg.Observability.Environment = getenvDefault("ENVIRONMENT", "development")
	
	return cfg, nil
}

func getenvDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}