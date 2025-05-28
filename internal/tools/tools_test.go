package tools

import (
	"slices"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestNormalizeCurrencyPair runs tests on the normalizeCurrencyPair function
func TestNormalizeCurrencyPair(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple BTC to XBT", "BTC", "XBT"},
		{"BTC in pair", "BTCGBP", "XBTGBP"},
		{"BTC with hyphen separator", "BTC-GBP", "XBTGBP"},
		{"BTC with slash separator", "BTC/GBP", "XBTGBP"},
		{"BTC with underscore separator", "BTC_GBP", "XBTGBP"},
		{"Lowercase input", "btcgbp", "XBTGBP"},
		{"Mixed case input", "xbTGbP", "XBTGBP"},
		{"Non-BTC pair", "ETHZAR", "ETHZAR"},
		{"Non-BTC pair with separator", "ETH-ZAR", "ETHZAR"},
		{"BITCOIN text conversion", "BITCOIN", "XBT"},
		{"BITCOIN in pair", "BITCOINUSD", "XBTUSD"},
		{"Multiple separators", "BTC-_/GBP", "XBTGBP"},
		{"Combo of mappings", "BITCOIN/GBP", "XBTGBP"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeCurrencyPair(tc.input)
			if result != tc.expected {
				t.Errorf("normalizeCurrencyPair(%q) = %q, want %q",
					tc.input, result, tc.expected)
			}
		})
	}
}

// TestGetWorkingPairs tests the GetWorkingPairs function
func TestGetWorkingPairs(t *testing.T) {
	// This just tests that the function returns some pairs
	// Since the implementation may change to return dynamic results
	pairs := GetWorkingPairs()
	if len(pairs) == 0 {
		t.Error("GetWorkingPairs() returned empty list, expected some pairs")
	}

	// Check that known essential pairs are included
	essentialPairs := []string{"XBTZAR", "XBTGBP"}

	for _, essentialPair := range essentialPairs {
		if !slices.Contains(pairs, essentialPair) {
			t.Errorf("GetWorkingPairs() missing essential pair %s", essentialPair)
		}
	}
}

// TestFindSimilarPairs tests the findSimilarPairs function using table-driven tests
func TestFindSimilarPairs(t *testing.T) {
	tests := []struct {
		name           string
		inputPair      string
		expectedResult string
		expectResults  bool
	}{
		{
			name:           "BTC should suggest XBT alternatives",
			inputPair:      "BTCGBP",
			expectedResult: "XBTGBP",
			expectResults:  true,
		},
		{
			name:           "XBT prefix should return some suggestions",
			inputPair:      "XBTUSD",
			expectedResult: "", // We don't care about specific pairs, just that we get some
			expectResults:  true,
		},
		{
			name:           "Invalid pair should return fallback suggestions",
			inputPair:      "INVALIDPAIR",
			expectedResult: "", // We don't care about specific pairs, just that we get some
			expectResults:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			suggestions := findSimilarPairs(tc.inputPair)

			// Test that we get at least some suggestions when expected
			if tc.expectResults && len(suggestions) == 0 {
				t.Errorf("findSimilarPairs(%q) returned no suggestions, expected some",
					tc.inputPair)
			}

			// If we're looking for a specific result, check for it
			if tc.expectedResult != "" {
				if !slices.Contains(suggestions, tc.expectedResult) {
					t.Errorf("findSimilarPairs(%q) did not include expected suggestion %q",
						tc.inputPair, tc.expectedResult)
				}
			}
		})
	}
}

func TestToolCreation(t *testing.T) {
	tests := []struct {
		name     string
		toolFunc func() mcp.Tool
		toolName string
		params   []string
	}{
		{
			name:     "GetBalances tool",
			toolFunc: NewGetBalancesTool,
			toolName: GetBalancesToolID,
			params:   []string{},
		},
		{
			name:     "GetTicker tool",
			toolFunc: NewGetTickerTool,
			toolName: GetTickerToolID,
			params:   []string{"pair"},
		},
		{
			name:     "GetOrderBook tool",
			toolFunc: NewGetOrderBookTool,
			toolName: GetOrderBookToolID,
			params:   []string{"pair"},
		},
		{
			name:     "CreateOrder tool",
			toolFunc: NewCreateOrderTool,
			toolName: CreateOrderToolID,
			params:   []string{"pair", "type", "volume", "price"},
		},
		{
			name:     "CancelOrder tool",
			toolFunc: NewCancelOrderTool,
			toolName: CancelOrderToolID,
			params:   []string{"order_id"},
		},
		{
			name:     "ListOrders tool",
			toolFunc: NewListOrdersTool,
			toolName: ListOrdersToolID,
			params:   []string{"pair", "limit"},
		},
		{
			name:     "ListTransactions tool",
			toolFunc: NewListTransactionsTool,
			toolName: ListTransactionsToolID,
			params:   []string{"account_id", "min_row", "max_row"},
		},
		{
			name:     "GetTransaction tool",
			toolFunc: NewGetTransactionTool,
			toolName: GetTransactionToolID,
			params:   []string{"account_id", "transaction_id"},
		},
		{
			name:     "ListTrades tool",
			toolFunc: NewListTradesTool,
			toolName: ListTradesToolID,
			params:   []string{"pair", "since"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := tt.toolFunc()

			if tool.Name != tt.toolName {
				t.Errorf("Expected tool name %q, got %q", tt.toolName, tool.Name)
			}

			if tool.Description == "" {
				t.Error("Tool description should not be empty")
			}

			// Verify tool has proper schema structure
			if tool.InputSchema.Type == "" {
				t.Error("Tool should have an input schema type")
				return
			}

			// Verify expected parameters exist
			if tool.InputSchema.Properties == nil {
				if len(tt.params) > 0 {
					t.Error("Tool should have properties for parameters")
				}
				return
			}

			for _, param := range tt.params {
				if _, exists := tool.InputSchema.Properties[param]; !exists {
					t.Errorf("Expected parameter %q to exist", param)
				}
			}
		})
	}
}
