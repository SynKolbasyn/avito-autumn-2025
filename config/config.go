// Package config provides utilities for application configuration management.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

var logMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
}

type Config struct {
	ServerPort int    `env:"SERVER_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST"     envDefault:"postgres"`
	DBPort     int    `env:"DB_PORT"     envDefault:"5432"`
	DBUser     string `env:"DB_USER"     envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME"     envDefault:"autumn-2025"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	LogLevel   zerolog.Level
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error("failed to parse environment variables", "error", err)

		return nil, fmt.Errorf("parse env: %w", err)
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if level, ok := logMap[strings.ToLower(logLevel)]; ok {
		cfg.LogLevel = level
	} else {
		cfg.LogLevel = zerolog.InfoLevel
	}

	slog.Info("configuration loaded")
	slog.Info("server_port", cfg.ServerPort)
	slog.Info("db_host", cfg.DBHost)
	slog.Info("db_port", cfg.DBPort)
	slog.Info("db_name", cfg.DBName)

	return cfg, nil
}
