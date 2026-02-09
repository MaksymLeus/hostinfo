package server

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	custom_mw "hostinfo/internal/api/middleware"
	v1 "hostinfo/internal/api/v1"
	"hostinfo/internal/custom"
	"hostinfo/internal/health"
)

func Run() {
	e := echo.New()

	// Core middleware
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(custom_mw.Logging())
	e.Use(custom_mw.CORS())

	// Health endpoints (no rate limit)
	e.GET("/healthz", health.Health)
	e.GET("/healthz/live", health.Live)
	e.GET("/healthz/ready", health.Ready)

	// API v1
	api := e.Group("/api/v1")
	api.Use(custom_mw.RateLimit())

	v1.Register(api)

	addr := fmt.Sprintf(
		"%s:%s",
		custom.GetEnv("HOSTINFO_ADDR", "0.0.0.0"),
		custom.GetEnv("HOSTINFO_PORT", "8080"),
	)

	log.Printf("Server listening on %s", addr)
	log.Fatal(e.Start(addr))
}
