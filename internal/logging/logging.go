// Package logging provides MCP-compatible logging functionality
//
// This package bridges Go's standard logging (slog) with the Model Context Protocol (MCP)
// notification system. It serves several important purposes:
//
//  1. Adapts standard Go logging to MCP's notification system, which allows using familiar
//     logging functions (slog.Info, slog.Debug, etc.) throughout the codebase.
//
//  2. Enables logging to multiple outputs simultaneously via MultiHandler - console logs
//     for local debugging and MCP notifications for client awareness.
//
//  3. Handles the conversion between Go's log levels and MCP log levels automatically.
//
//  4. Maintains a single, consistent logging approach across the application while still
//     delivering logs to both local console and MCP clients.
//
//  5. Makes codebase maintenance easier by centralizing logging logic, allowing future
//     changes to logging behavior without modifying every log statement.
//
// Without this package, we would need to have duplicate logging calls throughout
// the codebase - one for console output and another for MCP client notifications.
package logging

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MultiHandler is a handler that forwards records to multiple handlers
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a handler that forwards records to multiple handlers
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

// Enabled implements slog.Handler
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle implements slog.Handler
func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs implements slog.Handler
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

// WithGroup implements slog.Handler
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

// MCPNotificationHandler is a handler that sends logs as MCP notifications
type MCPNotificationHandler struct {
	s     *server.MCPServer
	level slog.Level
}

// NewMCPNotificationHandler creates a new handler that forwards logs to MCP clients
func NewMCPNotificationHandler(s *server.MCPServer, level slog.Level) *MCPNotificationHandler {
	return &MCPNotificationHandler{
		s:     s,
		level: level,
	}
}

// Enabled implements slog.Handler
func (h *MCPNotificationHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle implements slog.Handler
func (h *MCPNotificationHandler) Handle(ctx context.Context, record slog.Record) error {
	// Convert slog level to MCP logging level
	level := slogLevelToMCPLevel(record.Level)

	// Extract the message
	message := record.Message

	// Create a logging message notification using the MCP helper function
	notification := mcp.NewLoggingMessageNotification(level, "mcp-luno", message)

	// Send the notification to all clients - need to create a map to pass the params correctly
	h.s.SendNotificationToAllClients(notification.Method, map[string]any{
		"level":  string(level),
		"logger": "mcp-luno",
		"data":   message,
	})

	return nil
}

// WithAttrs implements slog.Handler
func (h *MCPNotificationHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we ignore attrs in this implementation
	return h
}

// WithGroup implements slog.Handler
func (h *MCPNotificationHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we ignore group in this implementation
	return h
}

// slogLevelToMCPLevel converts a slog.Level to an MCP LoggingLevel
func slogLevelToMCPLevel(level slog.Level) mcp.LoggingLevel {
	switch {
	case level <= slog.LevelDebug:
		return mcp.LoggingLevelDebug
	case level <= slog.LevelInfo:
		return mcp.LoggingLevelInfo
	case level <= slog.LevelWarn:
		return mcp.LoggingLevelWarning
	default:
		return mcp.LoggingLevelError
	}
}

// MCPHooks returns hooks for the MCP server that handle logging
func MCPHooks() *server.Hooks {
	// Create a new hooks instance
	hooks := &server.Hooks{}

	// Add request hook - using BeforeAny hook
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		// Log all incoming requests at debug level
		slog.DebugContext(ctx, "MCP request received",
			slog.String("method", string(method)),
			slog.Any("id", id))
	})

	// Add response hook - using OnSuccess hook
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		// Log all outgoing responses at debug level
		slog.DebugContext(ctx, "MCP response sent",
			slog.Any("id", id),
			slog.String("method", string(method)))
	})

	// Add error hook
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		// Log all errors at error level
		slog.ErrorContext(ctx, "MCP error occurred",
			slog.Any("error", err),
			slog.String("method", string(method)))
	})

	return hooks
}
