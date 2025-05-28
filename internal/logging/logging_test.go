package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	testMessage             = "test message"
	testNotificationMessage = "test notification message"
	unexpectedErrorFormat   = "Unexpected error: %v"
	notificationMethod      = "notifications/message"
	lunoMCPLogger           = "luno-mcp"
)

func TestNewMultiHandler(t *testing.T) {
	// Create mock handlers
	var buf1, buf2 bytes.Buffer
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := NewMultiHandler(handler1, handler2)

	if len(multiHandler.handlers) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(multiHandler.handlers))
	}
}

func TestMultiHandlerEnabled(t *testing.T) {
	tests := []struct {
		name            string
		handler1Level   slog.Level
		handler2Level   slog.Level
		testLevel       slog.Level
		expectedEnabled bool
	}{
		{
			name:            "enabled when one handler supports level",
			handler1Level:   slog.LevelInfo,
			handler2Level:   slog.LevelError,
			testLevel:       slog.LevelInfo,
			expectedEnabled: true,
		},
		{
			name:            "enabled when both handlers support level",
			handler1Level:   slog.LevelDebug,
			handler2Level:   slog.LevelInfo,
			testLevel:       slog.LevelInfo,
			expectedEnabled: true,
		},
		{
			name:            "disabled when no handlers support level",
			handler1Level:   slog.LevelWarn,
			handler2Level:   slog.LevelError,
			testLevel:       slog.LevelDebug,
			expectedEnabled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf1, buf2 bytes.Buffer
			handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: tc.handler1Level})
			handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: tc.handler2Level})

			multiHandler := NewMultiHandler(handler1, handler2)
			enabled := multiHandler.Enabled(context.Background(), tc.testLevel)

			if enabled != tc.expectedEnabled {
				t.Errorf("Expected enabled = %v, got %v", tc.expectedEnabled, enabled)
			}
		})
	}
}

func TestMultiHandlerHandle(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelInfo})
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelDebug})

	multiHandler := NewMultiHandler(handler1, handler2)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, testMessage, 0)
	err := multiHandler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf(unexpectedErrorFormat, err)
	}

	// Both handlers should have received the message since they're both enabled for info level
	if !strings.Contains(buf1.String(), testMessage) {
		t.Error("Handler1 should have received the message")
	}
	if !strings.Contains(buf2.String(), testMessage) {
		t.Error("Handler2 should have received the message")
	}
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{})
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{})

	multiHandler := NewMultiHandler(handler1, handler2)
	newHandler := multiHandler.WithAttrs([]slog.Attr{slog.String("key", "value")})

	// Verify it returns a MultiHandler
	if _, ok := newHandler.(*MultiHandler); !ok {
		t.Error("WithAttrs should return a MultiHandler")
	}
}

func TestMultiHandlerWithGroup(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{})
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{})

	multiHandler := NewMultiHandler(handler1, handler2)
	newHandler := multiHandler.WithGroup("testgroup")

	// Verify it returns a MultiHandler
	if _, ok := newHandler.(*MultiHandler); !ok {
		t.Error("WithGroup should return a MultiHandler")
	}
}

func TestNewMCPNotificationHandler(t *testing.T) {
	mockSender := &MockNotificationSender{}
	level := slog.LevelInfo

	handler := NewMCPNotificationHandler(mockSender, level)

	if handler.s != mockSender {
		t.Error("NotificationSender should be set correctly")
	}
	if handler.level != level {
		t.Error("Level should be set correctly")
	}
}

func TestMCPNotificationHandlerEnabled(t *testing.T) {
	tests := []struct {
		name            string
		handlerLevel    slog.Level
		testLevel       slog.Level
		expectedEnabled bool
	}{
		{
			name:            "level above threshold",
			handlerLevel:    slog.LevelInfo,
			testLevel:       slog.LevelError,
			expectedEnabled: true,
		},
		{
			name:            "level at threshold",
			handlerLevel:    slog.LevelInfo,
			testLevel:       slog.LevelInfo,
			expectedEnabled: true,
		},
		{
			name:            "level below threshold",
			handlerLevel:    slog.LevelInfo,
			testLevel:       slog.LevelDebug,
			expectedEnabled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSender := &MockNotificationSender{}
			handler := NewMCPNotificationHandler(mockSender, tc.handlerLevel)

			enabled := handler.Enabled(context.Background(), tc.testLevel)
			if enabled != tc.expectedEnabled {
				t.Errorf("Expected enabled = %v, got %v", tc.expectedEnabled, enabled)
			}
		})
	}
}

func TestMCPNotificationHandlerWithAttrs(t *testing.T) {
	mockSender := &MockNotificationSender{}
	handler := NewMCPNotificationHandler(mockSender, slog.LevelInfo)

	attrs := []slog.Attr{slog.String("key", "value")}
	newHandler := handler.WithAttrs(attrs)

	// Should return the same handler (simplified implementation)
	if newHandler != handler {
		t.Error("WithAttrs should return the same handler")
	}
}

func TestMCPNotificationHandlerWithGroup(t *testing.T) {
	mockSender := &MockNotificationSender{}
	handler := NewMCPNotificationHandler(mockSender, slog.LevelInfo)

	newHandler := handler.WithGroup("testgroup")

	// Should return the same handler (simplified implementation)
	if newHandler != handler {
		t.Error("WithGroup should return the same handler")
	}
}

func TestMCPNotificationHandlerHandle(t *testing.T) {
	mockSender := &MockNotificationSender{}

	// Set up expectation for SendNotificationToAllClients being called
	mockSender.On("SendNotificationToAllClients",
		notificationMethod,
		map[string]interface{}{
			"level":  "info",
			"logger": lunoMCPLogger,
			"data":   testNotificationMessage,
		},
	).Return().Once()

	handler := NewMCPNotificationHandler(mockSender, slog.LevelInfo)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, testNotificationMessage, 0)
	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf(unexpectedErrorFormat, err)
	}

	// Verify mock expectations
	mockSender.AssertExpectations(t)
}

func TestMCPNotificationHandlerLevelConversion(t *testing.T) {
	tests := []struct {
		name        string
		slogLevel   slog.Level
		expectedMCP string
	}{
		{
			name:        "debug level conversion",
			slogLevel:   slog.LevelDebug,
			expectedMCP: "debug",
		},
		{
			name:        "info level conversion",
			slogLevel:   slog.LevelInfo,
			expectedMCP: "info",
		},
		{
			name:        "warn level conversion",
			slogLevel:   slog.LevelWarn,
			expectedMCP: "warning",
		},
		{
			name:        "error level conversion",
			slogLevel:   slog.LevelError,
			expectedMCP: "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSender := &MockNotificationSender{}

			// Set up expectation for the correct MCP level
			mockSender.On("SendNotificationToAllClients",
				notificationMethod,
				map[string]interface{}{
					"level":  tc.expectedMCP,
					"logger": lunoMCPLogger,
					"data":   testMessage,
				},
			).Return().Once()

			handler := NewMCPNotificationHandler(mockSender, slog.LevelDebug) // Allow all levels

			record := slog.NewRecord(time.Now(), tc.slogLevel, testMessage, 0)
			err := handler.Handle(context.Background(), record)
			if err != nil {
				t.Errorf(unexpectedErrorFormat, err)
			}

			mockSender.AssertExpectations(t)
		})
	}
}

func TestMCPNotificationHandlerFiltering(t *testing.T) {
	mockSender := &MockNotificationSender{}

	// No expectations set - the mock should not receive any calls
	// because debug messages should be filtered when handler level is Info

	handler := NewMCPNotificationHandler(mockSender, slog.LevelInfo)

	// Try to handle a debug message (should be filtered)
	debugRecord := slog.NewRecord(time.Now(), slog.LevelDebug, "debug message", 0)
	err := handler.Handle(context.Background(), debugRecord)
	if err != nil {
		t.Errorf(unexpectedErrorFormat, err)
	}

	// Verify no calls were made to the mock
	mockSender.AssertExpectations(t)
}

func TestSlogLevelToMCPLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		expected mcp.LoggingLevel
	}{
		{
			name:     "debug level",
			level:    slog.LevelDebug,
			expected: mcp.LoggingLevelDebug,
		},
		{
			name:     "info level",
			level:    slog.LevelInfo,
			expected: mcp.LoggingLevelInfo,
		},
		{
			name:     "warn level",
			level:    slog.LevelWarn,
			expected: mcp.LoggingLevelWarning,
		},
		{
			name:     "error level",
			level:    slog.LevelError,
			expected: mcp.LoggingLevelError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := slogLevelToMCPLevel(tc.level)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestMCPHooks(t *testing.T) {
	hooks := MCPHooks()

	if hooks == nil {
		t.Error("MCPHooks should return non-nil hooks")
	}

	// Test that hooks don't panic when called
	// Note: We can't easily test the actual hook behavior without complex setup
	// but we can at least verify the hooks object is created properly
}

func TestMultiHandlerIntegration(t *testing.T) {
	var consoleBuffer bytes.Buffer
	consoleHandler := slog.NewTextHandler(&consoleBuffer, &slog.HandlerOptions{})

	mockSender := &MockNotificationSender{}
	mockSender.On("SendNotificationToAllClients",
		notificationMethod,
		map[string]interface{}{
			"level":  "info",
			"logger": lunoMCPLogger,
			"data":   "integration test message",
		},
	).Return().Once()

	mcpHandler := NewMCPNotificationHandler(mockSender, slog.LevelInfo)
	multiHandler := NewMultiHandler(consoleHandler, mcpHandler)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "integration test message", 0)
	err := multiHandler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf(unexpectedErrorFormat, err)
	}

	// Verify console handler also received the message
	if !strings.Contains(consoleBuffer.String(), "integration test message") {
		t.Error("Console handler should have received the message")
	}

	// Verify MCP handler received the message
	mockSender.AssertExpectations(t)
}
