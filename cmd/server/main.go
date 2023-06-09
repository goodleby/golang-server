package main

import (
	"context"
	"log"
	"os"

	"github.com/goodleby/golang-server/config"
	"github.com/goodleby/golang-server/server"
)

func main() {
	ctx := context.Background()

	config, err := config.Load(ctx, ".env")
	if err != nil {
		log.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	s, err := server.New(ctx, config)
	if err != nil {
		log.Printf("Error creating server: %v", err)
		os.Exit(1)
	}

	if err := s.Start(ctx); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
