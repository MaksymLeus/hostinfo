package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"hostinfo/internal/api"
	"hostinfo/internal/custom"
	"hostinfo/internal/host"
)

func main() {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	// Register REST API
	api.RegisterRoutes()

	// Health endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// UI endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = tmpl.Execute(w, host.Collect())
	})

	// Generate Server ADDRESS
	port := custom.GetEnv("HOSTINFO_PORT", "8080")
	hostAddr := custom.GetEnv("HOSTINFO_ADDR", "0.0.0.0")
	addr := fmt.Sprintf("%s:%s", hostAddr, port)

	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
