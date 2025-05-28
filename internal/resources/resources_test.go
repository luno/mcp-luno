package resources

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewWalletResource(t *testing.T) {
	resource := NewWalletResource()

	if resource.URI != WalletResourceURI {
		t.Errorf("Expected URI %q, got %q", WalletResourceURI, resource.URI)
	}

	if resource.Name != "Luno Wallets" {
		t.Errorf("Expected name %q, got %q", "Luno Wallets", resource.Name)
	}

	if resource.MIMEType != "application/json" {
		t.Errorf("Expected MIME type %q, got %q", "application/json", resource.MIMEType)
	}
}

func TestNewTransactionsResource(t *testing.T) {
	resource := NewTransactionsResource()

	if resource.URI != TransactionsResourceURI {
		t.Errorf("Expected URI %q, got %q", TransactionsResourceURI, resource.URI)
	}

	if resource.Name != "Luno Transactions" {
		t.Errorf("Expected name %q, got %q", "Luno Transactions", resource.Name)
	}

	if resource.MIMEType != "application/json" {
		t.Errorf("Expected MIME type %q, got %q", "application/json", resource.MIMEType)
	}
}

func TestNewAccountTemplate(t *testing.T) {
	template := NewAccountTemplate()

	if template.URITemplate.Raw() != AccountTemplateURI {
		t.Errorf("Expected URI template %q, got %q", AccountTemplateURI, template.URITemplate.Raw())
	}

	if template.Name != "Luno Account" {
		t.Errorf("Expected name %q, got %q", "Luno Account", template.Name)
	}
}

func TestExtractAccountID(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "valid account URI",
			uri:      "luno://accounts/12345",
			expected: "12345",
		},
		{
			name:     "account URI with longer ID",
			uri:      "luno://accounts/987654321",
			expected: "987654321",
		},
		{
			name:     "malformed URI with insufficient parts",
			uri:      "luno://accounts",
			expected: "",
		},
		{
			name:     "malformed URI with no parts",
			uri:      "",
			expected: "",
		},
		{
			name:     "malformed URI with single part",
			uri:      "accounts",
			expected: "",
		},
		{
			name:     "URI with extra path components",
			uri:      "luno://accounts/12345/extra/path",
			expected: "path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractAccountID(tc.uri)
			if result != tc.expected {
				t.Errorf("extractAccountID(%q) = %q, want %q", tc.uri, result, tc.expected)
			}
		})
	}
}

// TestHandleWalletResourceWithMockClient tests the wallet resource handler
// Note: This would require a mock Luno client to fully test without making real API calls
func TestHandleWalletResourceStructure(t *testing.T) {
	// Test that the handler can be created (structure test)
	// In a full implementation, you'd want to mock the config.Config and test the actual handler

	// This tests that the function signature is correct and can be called
	handler := HandleWalletResource(nil)
	if handler == nil {
		t.Error("HandleWalletResource should return a non-nil handler")
	}
}

// TestHandleTransactionsResourceStructure tests the transactions resource handler structure
func TestHandleTransactionsResourceStructure(t *testing.T) {
	handler := HandleTransactionsResource(nil)
	if handler == nil {
		t.Error("HandleTransactionsResource should return a non-nil handler")
	}
}

// TestHandleAccountTemplateStructure tests the account template handler structure
func TestHandleAccountTemplateStructure(t *testing.T) {
	handler := HandleAccountTemplate(nil)
	if handler == nil {
		t.Error("HandleAccountTemplate should return a non-nil handler")
	}
}

// Test that handlers return proper errors for invalid input
func TestHandleAccountTemplateWithInvalidInput(t *testing.T) {
	handler := HandleAccountTemplate(nil)

	// Test with empty URI
	request := mcp.ReadResourceRequest{}
	request.Params.URI = ""

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Error("Expected error for empty URI, got nil")
	}

	// Test with invalid URI format
	request.Params.URI = "invalid-uri-format"
	_, err = handler(context.Background(), request)
	if err == nil {
		t.Error("Expected error for invalid URI format, got nil")
	}
}
