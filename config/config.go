package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

type Config struct {
	ServerPort int    `env:"SERVER_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST"     envDefault:"postgres"`
	DBPort     int    `env:"DB_PORT"     envDefault:"5432"`
	DBUser     string `env:"DB_USER"     envDefault:"user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME"     envDefault:"database"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	LogLevel   zerolog.Level
}

func parseLogLevel(logLevelStr string) zerolog.Level {
	var logMap = map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
	}

	if level, ok := logMap[strings.ToLower(logLevelStr)]; ok {
		return level
	}

	return zerolog.InfoLevel
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		slog.Error("failed to parse environment variables", "error", err)

		return nil, fmt.Errorf("parse env: %w", err)
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	cfg.LogLevel = parseLogLevel(strings.TrimSpace(logLevelStr))

	slog.Info("configuration loaded",
		"server_port", cfg.ServerPort,
		"db_host", cfg.DBHost,
		"db_port", cfg.DBPort,
		"db_name", cfg.DBName,
	)

	return cfg, nil
}

func (cfg *Config) ServerAddress() string {
	return fmt.Sprintf(":%d", cfg.ServerPort)
}

func (cfg *Config) DBConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)
}
