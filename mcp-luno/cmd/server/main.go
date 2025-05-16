package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/echarrod/mcp-luno/internal/config"
	"github.com/echarrod/mcp-luno/internal/server"
)

const (
	appName    = "mcp-luno"
	appVersion = "0.1.0"
)

func main() {
	// Parse command line flags
	transportType := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	sseAddr := flag.String("sse-address", "localhost:8080", "Address for SSE transport")
	lunoDomain := flag.String("domain", "", "Luno API domain (default: api.luno.com)")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Set up logger
	level := parseLogLevel(*logLevel)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load(*lunoDomain)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create MCP server
	mcpServer := server.NewMCPServer(appName, appVersion, cfg)

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		slog.Info("Received shutdown signal")
		cancel()
	}()

	// Start the server with the selected transport
	switch *transportType {
	case "stdio":
		slog.Info("Starting Luno MCP server using stdio transport")
		if err := server.ServeStdio(ctx, mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "sse":
		slog.Info("Starting Luno MCP server using SSE transport", slog.String("address", *sseAddr))
		if err := server.ServeSSE(ctx, mcpServer, *sseAddr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf("Invalid transport type: %s. Must be 'stdio' or 'sse'", *transportType)
	}
}

func parseLogLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.LevelInfo
	}
	return l
}
