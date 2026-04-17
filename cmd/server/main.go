package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wieceslaw/chat-go/internal/server"
)

func main() {
	configFile := flag.String("config", "config/config.development.yaml", "Path to configuration file (required)")

	flag.Parse()

	// Check if config file is provided
	if *configFile == "" {
		fmt.Println("Error: config file parameter is required")
		fmt.Println("Usage: program -config <filename>")
		flag.PrintDefaults()
		log.Fatal("Error: config file parameter is required")
	}

	fmt.Printf("Using config file: %s\n", *configFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := server.New(ctx, *configFile)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	go func() {
		log.Printf("Server starting on %s", srv)
		if err := srv.Run(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	gracefulShutdown(srv, cancel)
}

func gracefulShutdown(srv *server.Server, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	cancel()

	ctx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shutdown gracefully: %v", err)
	}

	log.Println("Server exited properly")
}
