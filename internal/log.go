package internal

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ducminhgd/lusca/config"
	"github.com/go-chi/httplog/v2"
	gormslog "github.com/onrik/gorm-slog"
)

var (
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
)

func ConvertSlogLevel(l string) slog.Level {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func InitDBLogger(cfg *config.DBConfig) *gormslog.Logger {
	return gormslog.New(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     ConvertSlogLevel(cfg.LogLevel),
	})))
}

func InitAPILogger(cfg *config.Config) *httplog.Logger {
	logger := httplog.NewLogger("http-api", httplog.Options{
		JSON:             true,
		LogLevel:         ConvertSlogLevel(cfg.API.LogLevel),
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     cfg.Environemt,
		},
		QuietDownRoutes: []string{
			"/",
		},
	})
	return logger
}
