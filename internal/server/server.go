package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	custom_mw "hostinfo/internal/api/middleware"
	v1 "hostinfo/internal/api/v1"
	"hostinfo/internal/custom"
	"hostinfo/internal/health"
)

// EmbeddedFrontend holds the embedded frontend files
var EmbeddedFrontend embed.FS

// Server represents the HTTP server
type Server struct {
	FrontendPath string
	e            *echo.Echo
	allowOrigins []string
}

// New creates a new Server instance
// Returns an error if critical configuration validation fails (e.g., auth misconfigured)
func New() (*Server, error) {
	s := &Server{
		e:            echo.New(),
		FrontendPath: custom.GetEnv("HOSTINFO_FRONTEND_PATH", "./frontend/dist"),
	}

	s.setupRoutes()

	return s, nil
}

func (s *Server) Start() error {
	// Get allowed origins from environment, default to localhost only
	s.allowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		// Split by comma for multiple origins
		s.allowOrigins = strings.Split(envOrigins, ",")
		for i := range s.allowOrigins {
			s.allowOrigins[i] = strings.TrimSpace(s.allowOrigins[i])
		}
	}
	s.setupMiddleware()

	addr := fmt.Sprintf(
		"%s:%s",
		custom.GetEnv("HOSTINFO_HOST", "0.0.0.0"),
		custom.GetEnv("HOSTINFO_PORT", "8080"),
	)

	log.Printf("Server listening on http://%s", addr)
	log.Printf("Frontend path: %s", s.FrontendPath)
	log.Printf("CORS allowed origins: %v", s.allowOrigins)

	s.serveFrontend()

	return s.e.Start(addr)
}

// ---------------- API Handlers ----------------
func (s *Server) setupRoutes() {
	// Health endpoints (no rate limit)
	s.e.GET("/healthz", health.Health)
	s.e.GET("/healthz/live", health.Live)
	s.e.GET("/healthz/ready", health.Ready)

	// API v1
	api := s.e.Group("/api/v1")
	api.Use(custom_mw.RateLimit())
	v1.Register(api)
}

// ---------------- Middleware ----------------
func (s *Server) setupMiddleware() {
	// Core middleware
	s.e.HideBanner = true
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.RequestID())
	s.e.Use(custom_mw.Logging())
	s.e.Use(custom_mw.CORS(s.allowOrigins))
}

// ---------------- Frontend ----------------

func (s *Server) serveFrontend() {
	// Try to use filesystem path first (for development)
	if _, err := os.Stat(s.FrontendPath); err == nil {
		log.Printf("Serving frontend from filesystem: %s", s.FrontendPath)
		s.serveFrontendFromFilesystem()
		return
	}
	if os.Getenv("HOSTINFO_RUN_API_ONLY") == "true" {
		log.Println("HOSTINFO_RUN_API_ONLY is true, skipping frontend serving. API will be available at /api/v1/*")
	} else {
		// Fall back to embedded frontend (for production binaries)
		log.Println("Serving frontend from embedded files")
		s.serveFrontendFromEmbedded()
	}
}

// DEV: serve frontend from filesystem
func (s *Server) serveFrontendFromFilesystem() {
	s.e.GET("/*", func(c echo.Context) error {
		reqPath := c.Request().URL.Path
		fullPath := filepath.Join(s.FrontendPath, reqPath)

		// Asset exists → serve it
		if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
			return c.File(fullPath)
		}

		// SPA fallback
		return c.File(filepath.Join(s.FrontendPath, "index.html"))
	})
}

// PROD: serve frontend from embedded FS
func (s *Server) serveFrontendFromEmbedded() {
	buildFS, err := fs.Sub(EmbeddedFrontend, "frontend/dist")
	if err != nil {
		log.Printf("Failed to access embedded frontend: %v", err)
		s.serveErrorPage()
		return
	}
	// check if folder exists in embedded FS
	if _, err := fs.Stat(buildFS, "."); err != nil {
		log.Printf("Embedded frontend not found: %v", err)
		s.serveErrorPage()
		return
	}
	log.Printf("Serving embedded frontend with CORS allowed origins: %v\n", s.allowOrigins)
	fileServer := http.FileServer(http.FS(buildFS))

	s.e.GET("/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean(r.URL.Path)
		if cleanPath == "/" {
			cleanPath = "/index.html"
		}

		// Try requested file
		if _, err := fs.Stat(buildFS, cleanPath[1:]); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback → index.html
		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	})))
}

// Fallback page if frontend missing
func (s *Server) serveErrorPage() {
	log.Println("Serving error page for missing frontend")
	s.e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
  <title>Server Running</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f5f5f5; padding: 40px; }
    .box { background: white; padding: 30px; max-width: 600px; margin: auto; border-radius: 8px; }
  </style>
</head>
<body>
  <div class="box">
    <h1>Backend is running</h1>
    <p>Frontend build was not found.</p>
    <p>API available at <code>/api/v1/*</code></p>
  </div>
</body>
</html>
`)
	})
}
