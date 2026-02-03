package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"hostinfo/internal/api"
	"hostinfo/internal/api/middleware"
	"hostinfo/internal/custom"
	"hostinfo/internal/health"
	"hostinfo/internal/host"
)

func main() {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	// Register REST API
	apiMux := http.NewServeMux()
	api.RegisterRoutes(apiMux) // pass mux instead of using http.HandleFunc directly
	// api.RegisterRoutes()

	// Main mux
	mainMux := http.NewServeMux()

	// UI endpoint
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_ = tmpl.Execute(w, host.Collect())
	})

	// Mount API under /api/v1/
	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", middleware.RateLimiter(apiMux)))

	// Health endpoint
	mainMux.HandleFunc("/healthz", health.HealthHandler)
	mainMux.HandleFunc("/healthz/live", health.LiveHandler)
	mainMux.HandleFunc("/healthz/ready", health.ReadyHandler)

	// Generate Server ADDRESS
	port := custom.GetEnv("HOSTINFO_PORT", "8080")
	hostAddr := custom.GetEnv("HOSTINFO_ADDR", "0.0.0.0")
	addr := fmt.Sprintf("%s:%s", hostAddr, port)

	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mainMux))

}
