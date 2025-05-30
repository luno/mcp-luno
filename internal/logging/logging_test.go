package logging

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

const (
	loggerName         = "luno-mcp"
	logMsgMCPRequest   = "MCP request received"
	logMsgMCPResponse  = "MCP response sent"
	logMsgMCPError     = "MCP error occurred"
	logKeyMethod       = "method"
	logKeyID           = "id"
	logKeyError        = "error"
	testMessageDefault = "this is a test log message"
	testIntegrationMsg = "integration test message"

	// Corrected JSON string constants for assertions
	jsonLogLevelDebug     = `"level":"DEBUG"`
	jsonLogLevelError     = `"level":"ERROR"`
	jsonTestComponentAttr = `"component":"test"`
	jsonTestGroupAttrOpen = `"testGroup":{`
)

func TestSlogLevelToMCPLevel(t *testing.T) {
	testCases := []struct {
		name      string
		slogLevel slog.Level
		mcpLevel  mcp.LoggingLevel
	}{
		{"debug", slog.LevelDebug, mcp.LoggingLevelDebug},
		{"info", slog.LevelInfo, mcp.LoggingLevelInfo},
		{"warn", slog.LevelWarn, mcp.LoggingLevelWarning},
		{"error", slog.LevelError, mcp.LoggingLevelError},
		{"below debug", slog.LevelDebug - 4, mcp.LoggingLevelDebug},
		{"between info and warn", slog.LevelInfo + 2, mcp.LoggingLevelWarning},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.mcpLevel, slogLevelToMCPLevel(tc.slogLevel))
		})
	}
}

func TestMCPNotificationHandlerEnabled(t *testing.T) {
	handler := NewMCPNotificationHandler(&MockNotificationSender{}, slog.LevelInfo) // Only Info and above

	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo), "Info level should be enabled")
	assert.True(t, handler.Enabled(context.Background(), slog.LevelWarn), "Warn level should be enabled")
	assert.True(t, handler.Enabled(context.Background(), slog.LevelError), "Error level should be enabled")
	assert.False(t, handler.Enabled(context.Background(), slog.LevelDebug), "Debug level should be disabled")
}

func TestMCPNotificationHandlerHandleNotificationFormat(t *testing.T) {
	mockS := new(MockNotificationSender)
	handler := NewMCPNotificationHandler(mockS, slog.LevelDebug) // Enable all levels for this test

	level := slog.LevelInfo
	mcpLevel := slogLevelToMCPLevel(level)

	expectedParams := map[string]any{
		"level":  string(mcpLevel),
		"logger": loggerName,
		"data":   testMessageDefault,
	}
	expectedMethod := mcp.NewLoggingMessageNotification(mcpLevel, loggerName, testMessageDefault).Method

	mockS.On("SendNotificationToAllClients", expectedMethod, expectedParams).Return()

	record := slog.NewRecord(time.Now(), level, testMessageDefault, 0)
	err := handler.Handle(context.Background(), record)
	assert.NoError(t, err)

	mockS.AssertExpectations(t)
}

func TestMCPNotificationHandlerWithAttrsAndGroup(t *testing.T) {
	handler := NewMCPNotificationHandler(&MockNotificationSender{}, slog.LevelInfo)

	attrs := []slog.Attr{slog.String("key", "value")}
	handlerWithAttrs := handler.WithAttrs(attrs)
	assert.Equal(t, handler, handlerWithAttrs, "WithAttrs should return the same handler instance for simplicity")

	handlerWithGroup := handler.WithGroup("testGroup")
	assert.Equal(t, handler, handlerWithGroup, "WithGroup should return the same handler instance for simplicity")
}

func TestMultiHandlerEnabled(t *testing.T) {
	h1 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	h2 := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	multi := NewMultiHandler(h1, h2)

	assert.True(t, multi.Enabled(context.Background(), slog.LevelDebug), "Should be enabled if any handler is enabled for Debug")
	assert.True(t, multi.Enabled(context.Background(), slog.LevelInfo), "Should be enabled if any handler is enabled for Info")
	assert.True(t, multi.Enabled(context.Background(), slog.LevelWarn), "Should be enabled if any handler is enabled for Warn")

	multiOnlyH1 := NewMultiHandler(h1)
	assert.True(t, multiOnlyH1.Enabled(context.Background(), slog.LevelDebug))
	assert.False(t, NewMultiHandler().Enabled(context.Background(), slog.LevelDebug), "Should be false if no handlers")
}

func TestMultiHandlerHandle(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	h1 := slog.NewJSONHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelDebug})
	h2 := slog.NewJSONHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelInfo})
	multi := NewMultiHandler(h1, h2)

	debugRecord := slog.NewRecord(time.Now(), slog.LevelDebug, testMessageDefault, 0)
	infoRecord := slog.NewRecord(time.Now(), slog.LevelInfo, testMessageDefault, 0)

	buf1.Reset()
	buf2.Reset()
	err := multi.Handle(context.Background(), debugRecord)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(buf1.String(), testMessageDefault), "h1 should contain debug message")
	assert.Equal(t, "", buf2.String(), "h2 should not contain debug message")

	buf1.Reset()
	buf2.Reset()
	err = multi.Handle(context.Background(), infoRecord)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(buf1.String(), testMessageDefault), "h1 should contain info message")
	assert.True(t, strings.Contains(buf2.String(), testMessageDefault), "h2 should contain info message")
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	var buf1 bytes.Buffer
	h1 := slog.NewJSONHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelDebug})
	multi := NewMultiHandler(h1)

	attrs := []slog.Attr{slog.String("component", "test")}
	multiWithAttrs := multi.WithAttrs(attrs)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test with attrs", 0)
	err := multiWithAttrs.Handle(context.Background(), record)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(buf1.String(), jsonTestComponentAttr), "Log output should contain the attribute")
}

func TestMultiHandlerWithGroup(t *testing.T) {
	var buf1 bytes.Buffer
	h1 := slog.NewJSONHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelDebug})
	multi := NewMultiHandler(h1)

	multiWithGroup := multi.WithGroup("testGroup")

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test with group", 0)
	// Add an attribute to the record
	record.AddAttrs(slog.String("attrKey", "attrValue"))

	err := multiWithGroup.Handle(context.Background(), record)
	assert.NoError(t, err)

	// Add this log to see the exact output in the test logs
	t.Logf("TestMultiHandlerWithGroup actual output: %s", buf1.String())

	assert.True(t, strings.Contains(buf1.String(), jsonTestGroupAttrOpen), "Log output should contain the group opening")
	assert.True(t, strings.Contains(buf1.String(), "attrKey"), "Log output should contain the attribute key")
	assert.True(t, strings.Contains(buf1.String(), "attrValue"), "Log output should contain the attribute value")
}

func TestLoggingHooksExecution(t *testing.T) {
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	ctx := context.Background()

	t.Run("LogRequestHook logs request details", func(t *testing.T) {
		logOutput.Reset()
		id := "request-001"
		method := mcp.MCPMethod("system.listMethods")
		message := map[string]any{"service": "test"}

		LogRequestHook(ctx, id, method, message)

		output := logOutput.String()
		assert.Contains(t, output, logMsgMCPRequest)
		assert.Contains(t, output, `"`+logKeyMethod+`":"system.listMethods"`)
		assert.Contains(t, output, `"`+logKeyID+`":"request-001"`)
		assert.Contains(t, output, jsonLogLevelDebug)
	})

	t.Run("LogSuccessHook logs response details", func(t *testing.T) {
		logOutput.Reset()
		id := "request-002"
		method := mcp.MCPMethod("order.create")
		message := map[string]any{"pair": "XBTEUR", "price": "10000"}
		result := map[string]any{"order_id": "ord-123", "status": "pending"}

		LogSuccessHook(ctx, id, method, message, result)

		output := logOutput.String()
		assert.Contains(t, output, logMsgMCPResponse)
		assert.Contains(t, output, `"`+logKeyMethod+`":"order.create"`)
		assert.Contains(t, output, `"`+logKeyID+`":"request-002"`)
		assert.Contains(t, output, jsonLogLevelDebug)
	})

	t.Run("LogErrorHook logs error details", func(t *testing.T) {
		logOutput.Reset()
		id := "request-003"
		method := mcp.MCPMethod("account.getBalance")
		message := map[string]any{"asset": "XBT"}
		testErr := errors.New("simulated API failure")

		LogErrorHook(ctx, id, method, message, testErr)

		output := logOutput.String()
		assert.Contains(t, output, logMsgMCPError)
		assert.Contains(t, output, `"`+logKeyID+`":"request-003"`)
		assert.Contains(t, output, `"`+logKeyMethod+`":"account.getBalance"`)
		assert.Contains(t, output, `"`+logKeyError+`":"simulated API failure"`)
		assert.Contains(t, output, jsonLogLevelError)
	})
}

func TestIntegrationHooksWithNotificationHandler(t *testing.T) {
	var consoleBuffer bytes.Buffer
	mockNotifier := new(MockNotificationSender) // Create a new mock for each test run or sub-test for isolation

	consoleHandler := slog.NewJSONHandler(&consoleBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
	mcpNotificationHandler := NewMCPNotificationHandler(mockNotifier, slog.LevelDebug)
	multiHandler := NewMultiHandler(consoleHandler, mcpNotificationHandler)
	// Set this multiHandler as the default logger for the duration of this test
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(multiHandler))
	defer slog.SetDefault(originalLogger) // Restore original logger after test

	ctx := context.Background()

	t.Run("LogRequestHook interaction", func(t *testing.T) {
		consoleBuffer.Reset()
		// Reset mock expectations for this sub-test.
		// testify/mock doesn't have a direct Reset() for all expectations on the mock object itself.
		// Instead, we create a new mock or manage expectations per test.
		// For this structure, re-initialize mockNotifier for true isolation if needed, or ensure .On()...Once() is specific enough.
		mockNotifier = new(MockNotificationSender) // Re-initialize for this sub-test
		// Re-setup the handler with the new mock if it captures the mock instance
		mcpNotificationHandler := NewMCPNotificationHandler(mockNotifier, slog.LevelDebug)
		multiHandler := NewMultiHandler(consoleHandler, mcpNotificationHandler)
		slog.SetDefault(slog.New(multiHandler)) // Re-set default logger with new handler chain

		reqID := "req-integ-001"
		reqMethod := mcp.MCPMethod("test.integration")

		expectedNotificationParams := map[string]any{
			"level":  string(mcp.LoggingLevelDebug),
			"logger": loggerName,
			"data":   logMsgMCPRequest,
		}
		notification := mcp.NewLoggingMessageNotification(mcp.LoggingLevelDebug, loggerName, logMsgMCPRequest)
		mockNotifier.On("SendNotificationToAllClients", notification.Method, expectedNotificationParams).Once()

		LogRequestHook(ctx, reqID, reqMethod, nil)

		assert.True(t, strings.Contains(consoleBuffer.String(), logMsgMCPRequest), "Console log should contain BeforeAny message")
		assert.True(t, strings.Contains(consoleBuffer.String(), string(reqMethod)), "Console log should contain method")
		mockNotifier.AssertExpectations(t)
	})

	t.Run("LogErrorHook interaction", func(t *testing.T) {
		consoleBuffer.Reset()
		mockNotifier = new(MockNotificationSender) // Re-initialize for this sub-test
		mcpNotificationHandler := NewMCPNotificationHandler(mockNotifier, slog.LevelDebug)
		multiHandler := NewMultiHandler(consoleHandler, mcpNotificationHandler)
		slog.SetDefault(slog.New(multiHandler))

		errID := "err-integ-002"
		errMethod := mcp.MCPMethod("test.errorIntegration")
		testErr := errors.New("integration error")

		expectedErrorNotificationParams := map[string]any{
			"level":  string(mcp.LoggingLevelError),
			"logger": loggerName,
			"data":   logMsgMCPError,
		}
		errorNotification := mcp.NewLoggingMessageNotification(mcp.LoggingLevelError, loggerName, logMsgMCPError)
		mockNotifier.On("SendNotificationToAllClients", errorNotification.Method, expectedErrorNotificationParams).Once()

		LogErrorHook(ctx, errID, errMethod, nil, testErr)

		assert.True(t, strings.Contains(consoleBuffer.String(), logMsgMCPError))
		assert.True(t, strings.Contains(consoleBuffer.String(), string(errMethod)))
		assert.True(t, strings.Contains(consoleBuffer.String(), "integration error"))
		mockNotifier.AssertExpectations(t)
	})
}
