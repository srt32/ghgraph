package main

import (
	"log"
	"os"
	"strconv"

	"github.com/srt32/ghgraph/internal/server"
)

func main() {
	// Get port from environment variable or use default
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	// Create and start server
	srv := server.NewServer(port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
