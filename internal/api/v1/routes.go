package v1

import (
	"hostinfo/internal/api/v1/handlers"

	"github.com/labstack/echo/v4"
)

// Register registers all v1 API routes
// @title Hostinfo API v1
// @version 1.0
// @description Hostinfo server API documentation.
// @BasePath /api/v1
func Register(g *echo.Group) {
	g.GET("/info", handlers.Info)
	g.GET("/cloud", handlers.Cloud)
	g.GET("/kubernetes", handlers.Kubernetes)

	g.POST("/ping", handlers.Ping)
	g.POST("/curl", handlers.Curl)
	g.POST("/dns", handlers.Dig)
	g.POST("/tcp", handlers.TCP)
}
