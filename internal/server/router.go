package server

import (
	"autumn-2025/internal/handlers"
	"autumn-2025/internal/repositories/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers HTTP routes on the provided Echo instance.
func RegisterRoutes(server *echo.Echo, pool *pgxpool.Pool) {
	server.GET("/health", handlers.GetHealth)
	registerTeamRoutes(server, pool)
	registerPullRequestRoutes(server, pool)
}

func registerTeamRoutes(server *echo.Echo, pool *pgxpool.Pool) {
	teamRepository := postgres.NewTeamRepository(pool)
	teamHandler := handlers.NewTeamHandler(teamRepository)
	group := server.Group("/team")
	group.POST("/add", teamHandler.Add)
	group.GET("/get", teamHandler.Get)
}

func registerPullRequestRoutes(server *echo.Echo, pool *pgxpool.Pool) {
	pullRequestRepository := postgres.NewPullRequestRepository(pool)
	teamHandler := handlers.NewPullRequestHandler(pullRequestRepository)
	group := server.Group("/pullRequest")
	group.POST("/create", teamHandler.Create)
	group.POST("/merge", teamHandler.Merge)
	group.POST("/reassign", teamHandler.Reassign)
}
