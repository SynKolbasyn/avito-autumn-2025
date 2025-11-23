package database

import (
	"autumn-2025/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Pool(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.DBConnectionString())
	if err != nil {
		log.Fatal().Err(err).Str("server_address", cfg.DBConnectionString()).Msg("failed to create database pool")
	}

	return pool
}
