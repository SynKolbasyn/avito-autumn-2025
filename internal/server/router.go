package server

import (
	"autumn-2025/internal/handlers"
	"autumn-2025/internal/repositories/postgres"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers HTTP routes on the provided Echo instance.
func RegisterRoutes(server *echo.Echo, ctx context.Context, pool *pgxpool.Pool) {
	server.GET("/health", handlers.GetHealth)
	registerTeamRoutes(server, ctx, pool)
}

func registerTeamRoutes(server *echo.Echo, ctx context.Context, pool *pgxpool.Pool) {
	teamRepository := postgres.NewTeamRepository(pool)
	teamHandler := handlers.NewTeamHandler(ctx, teamRepository)
	group := server.Group("/team")
	group.POST("/add", teamHandler.Add)
	group.GET("/get", teamHandler.Get)
}
