package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORS(allowOrigins []string) echo.MiddlewareFunc {
	if len(allowOrigins) == 0 {
		// Default to localhost if no origins provided
		allowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"GET", "POST", "PUT"},
	})
}
