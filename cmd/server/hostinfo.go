package main

import (
	"hostinfo/assets"
	"hostinfo/internal/server"
	"log"
)

func main() {
	// Set the embedded frontend files from the assets package
	server.EmbeddedFrontend = assets.FrontendFS

	// Initialize and start the server
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)

	}
	log.Fatal(srv.Start())
}
