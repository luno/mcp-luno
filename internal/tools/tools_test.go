package tools

import (
	"testing"
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
	essentialPairs := map[string]bool{
		"XBTZAR": false,
		"XBTGBP": false,
	}

	for _, pair := range pairs {
		if _, ok := essentialPairs[pair]; ok {
			essentialPairs[pair] = true
		}
	}

	for pair, found := range essentialPairs {
		if !found {
			t.Errorf("GetWorkingPairs() missing essential pair %s", pair)
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
				found := false
				for _, pair := range suggestions {
					if pair == tc.expectedResult {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("findSimilarPairs(%q) did not include expected suggestion %q",
						tc.inputPair, tc.expectedResult)
				}
			}
		})
	}
}

// TestContainsPair tests the containsPair helper function
func TestContainsPair(t *testing.T) {
	testPairs := []string{"XBTZAR", "ETHZAR", "XBTUSD"}

	tests := []struct {
		name     string
		pair     string
		expected bool
	}{
		{"Existing pair", "XBTZAR", true},
		{"Missing pair", "XBTGBP", false},
		{"Case sensitive check", "xbtzar", false}, // Should be case sensitive
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := containsPair(testPairs, tc.pair)
			if result != tc.expected {
				t.Errorf("containsPair(%v, %q) = %v, want %v",
					testPairs, tc.pair, result, tc.expected)
			}
		})
	}
}
