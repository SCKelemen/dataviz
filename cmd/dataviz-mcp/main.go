package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SCKelemen/dataviz/mcp/mcp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping server...")
		cancel()
	}()

	// Create and run server
	server, err := mcp.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Println("DataViz MCP server starting...")
	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
