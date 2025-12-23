package observability

import (
	"fmt"
	"havamal-api/config"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/loki-client-go/pkg/urlutil"
	promConfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

var (
	Logger     *slog.Logger
	lokiClient *loki.Client
)

// LokiWriter implements io.Writer to send logs to Loki
type LokiWriter struct {
	client      *loki.Client
	serviceName string
	environment string
}

func (w *LokiWriter) Write(p []byte) (n int, err error) {
	// Parse log level from JSON if possible (basic parsing)
	logLevel := "info"
	logStr := string(p)
	if strings.Contains(logStr, `"level":"error"`) || strings.Contains(logStr, `"level":"ERROR"`) {
		logLevel = "error"
	} else if strings.Contains(logStr, `"level":"warn"`) || strings.Contains(logStr, `"level":"WARN"`) {
		logLevel = "warn"
	} else if strings.Contains(logStr, `"level":"debug"`) || strings.Contains(logStr, `"level":"DEBUG"`) {
		logLevel = "debug"
	}

	labels := model.LabelSet{
		"service":     model.LabelValue(w.serviceName),
		"environment": model.LabelValue(w.environment),
		"level":       model.LabelValue(logLevel),
	}

	if err := w.client.Handle(labels, time.Now(), logStr); err != nil {
		// Log to stderr if Loki fails, but don't block
		fmt.Fprintf(os.Stderr, "Failed to send to Loki: %v\n", err)
	}

	return len(p), nil
}

// InitializeLogging sets up structured logging with Loki
func InitializeLogging(cfg config.Config) error {
	var writer io.Writer = os.Stdout

	// Si tenim configuració de Loki, enviem logs allà també
	if cfg.Observability.LokiURL != "" && 	   
	   cfg.Observability.LokiPassword != "" {
		
		lokiCfg := loki.Config{}
		
		u, err := url.Parse(cfg.Observability.LokiURL)
		if err != nil {
			return fmt.Errorf("failed to parse loki url: %w", err)
		}
		lokiCfg.URL = urlutil.URLValue{URL: u}
		
		lokiCfg.TenantID = ""
		lokiCfg.BatchWait = 1 * time.Second
		lokiCfg.BatchSize = 1024 * 1024 // 1MB
		lokiCfg.Timeout = 10 * time.Second
		lokiCfg.BackoffConfig.MinBackoff = 500 * time.Millisecond
		lokiCfg.BackoffConfig.MaxBackoff = 5 * time.Minute
		lokiCfg.BackoffConfig.MaxRetries = 10
		
		username := cfg.Observability.LokiUser
		if username == "" {
			username = "turniq"
		}

		lokiCfg.Client.BasicAuth = &promConfig.BasicAuth{
			Username: username,
			Password: promConfig.Secret(cfg.Observability.LokiPassword),
		}

		client, err := loki.New(lokiCfg)
		if err != nil {
			return fmt.Errorf("failed to create loki client: %w", err)
		}

		lokiClient = client

		// MultiWriter per enviar logs tant a stdout com a Loki
		lokiWriter := &LokiWriter{
			client:      client,
			serviceName: cfg.Observability.ServiceName,
			environment: cfg.Observability.Environment,
		}
		writer = io.MultiWriter(os.Stdout, lokiWriter)
		
		fmt.Printf("Loki configured: %s\n", cfg.Observability.LokiURL)
	} else {
		fmt.Println("Loki not configured, logging to stdout only")
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Logger = slog.New(handler).With(
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
	)

	slog.SetDefault(Logger)
	Logger.Info("Logger initialized",
		slog.Bool("loki_enabled", cfg.Observability.LokiURL != ""),
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
	)

	return nil
}

// ShutdownLogging gracefully shuts down logging components
func ShutdownLogging() {
	if lokiClient != nil {
		Logger.Info("Stopping Loki client")
		lokiClient.Stop()
	}
}
