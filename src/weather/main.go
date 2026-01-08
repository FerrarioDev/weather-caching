package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FerrarioDev/weathercache/config"
	"github.com/FerrarioDev/weathercache/internal/service"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Get server port from config
	port := config.ServerPort.GetValue()
	if port == "" {
		port = ":8084" // default value
	}

	// Create the server
	appServer := service.NewServer(port)

	// Start server in a goroutine
	go func() {
		if err := appServer.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down server...")
	if err := appServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	fmt.Println("Server stopped")
}
