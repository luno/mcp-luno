package resources

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/luno/luno-mcp/internal/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

const (
	expectedMIMEType = "application/json"
	expectedNameFmt  = "Expected name %q, got %q"
)

func TestNewWalletResource(t *testing.T) {
	resource := NewWalletResource()

	assert.Equal(t, WalletResourceURI, resource.URI)
	assert.Equal(t, "Luno Wallets", resource.Name)
	assert.Equal(t, expectedMIMEType, resource.MIMEType)
}

func TestNewTransactionsResource(t *testing.T) {
	resource := NewTransactionsResource()

	assert.Equal(t, TransactionsResourceURI, resource.URI)
	assert.Equal(t, "Luno Transactions", resource.Name)
	assert.Equal(t, expectedMIMEType, resource.MIMEType)
}

func TestNewAccountTemplate(t *testing.T) {
	expectedJSON := `{
		"uriTemplate": "luno://accounts/{id}",
		"name": "Luno Account",
		"description": "Returns details for a specific Luno account"
	}`

	template := NewAccountTemplate()

	actualJSON, err := json.Marshal(template)
	assert.NoError(t, err)

	// Compare JSON structures directly. Can't create an expected object as the fields can only be set internally.
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

func TestExtractAccountID(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{"valid account URI", "luno://accounts/1234567890", "1234567890"},
		{"empty URI", "", ""},
		{"invalid format", "luno://accounts", ""},
		{"short URI", "luno://", ""},
		{"no account ID", "luno://accounts/", ""},
		{"different resource", "luno://wallets/123", "123"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractAccountID(tc.uri)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestHandleWalletResourceStructure tests that the wallet resource handler can be created
func TestHandleWalletResourceStructure(t *testing.T) {
	handler := HandleWalletResource(nil)
	assert.NotNil(t, handler, "HandleWalletResource should return a non-nil handler")
}

// TestHandleTransactionsResourceStructure tests the transactions resource handler structure
func TestHandleTransactionsResourceStructure(t *testing.T) {
	handler := HandleTransactionsResource(nil)
	assert.NotNil(t, handler, "HandleTransactionsResource should return a non-nil handler")
}

// TestHandleAccountTemplateStructure tests the account template handler structure
func TestHandleAccountTemplateStructure(t *testing.T) {
	handler := HandleAccountTemplate(nil)
	assert.NotNil(t, handler, "HandleAccountTemplate should return a non-nil handler")
}

// createTestConfig creates a configuration for testing
func createTestConfig() *config.Config {
	// For testing, we create a config with a nil client
	// In real integration tests, this would be a properly configured client
	return &config.Config{
		LunoClient: nil,
	}
}

// TestHandleWalletResourceIntegration tests the wallet resource handler structure and behavior
func TestHandleWalletResourceIntegration(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name:        "config with nil client",
			config:      createTestConfig(),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := HandleWalletResource(tc.config)
			assert.NotNil(t, handler, "HandleWalletResource should return a non-nil handler")

			req := mcp.ReadResourceRequest{
				Params: struct {
					URI       string         `json:"uri"`
					Arguments map[string]any `json:"arguments,omitempty"`
				}{
					URI: WalletResourceURI,
				},
			}

			result, err := handler(context.Background(), req)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestHandleTransactionsResourceIntegration tests the transactions resource handler structure and behavior
func TestHandleTransactionsResourceIntegration(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
		},
		{
			name:        "config with nil client",
			config:      createTestConfig(),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := HandleTransactionsResource(tc.config)
			assert.NotNil(t, handler, "HandleTransactionsResource should return a non-nil handler")

			req := mcp.ReadResourceRequest{
				Params: struct {
					URI       string         `json:"uri"`
					Arguments map[string]any `json:"arguments,omitempty"`
				}{
					URI: TransactionsResourceURI,
				},
			}

			result, err := handler(context.Background(), req)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestHandleAccountTemplateIntegration tests the account template handler structure and behavior
func TestHandleAccountTemplateIntegration(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		uri         string
		expectError bool
	}{
		{
			name:        "nil config",
			config:      nil,
			uri:         "luno://accounts/1234567890",
			expectError: true,
		},
		{
			name:        "config with nil client",
			config:      createTestConfig(),
			uri:         "luno://accounts/1234567890",
			expectError: true,
		},
		{
			name:        "invalid URI format",
			config:      createTestConfig(),
			uri:         "invalid://uri",
			expectError: true,
		},
		{
			name:        "empty account ID",
			config:      createTestConfig(),
			uri:         "luno://accounts/",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := HandleAccountTemplate(tc.config)
			assert.NotNil(t, handler, "HandleAccountTemplate should return a non-nil handler")

			req := mcp.ReadResourceRequest{
				Params: struct {
					URI       string         `json:"uri"`
					Arguments map[string]any `json:"arguments,omitempty"`
				}{
					URI: tc.uri,
				},
			}

			result, err := handler(context.Background(), req)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
