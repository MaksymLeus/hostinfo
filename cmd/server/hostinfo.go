package main

import (
	"hostinfo/internal/server"
	"log"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)

	}
	log.Fatal(srv.Start())
}
