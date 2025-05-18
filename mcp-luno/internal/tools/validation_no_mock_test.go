package tools

import (
	"testing"
)

// TestValidatePairNoAPI tests just the basic functionality of ValidatePair
// without making actual API calls
func TestValidatePairNoAPI(t *testing.T) {
	// Initialize the cache for testing
	validPairsCache = make(map[string]bool)
	discoveredPairsCache = []string{"XBTZAR", "ETHZAR", "XBTGBP"}

	// Pre-populate the cache
	for _, pair := range discoveredPairsCache {
		validPairsCache[pair] = true
	}

	tests := []struct {
		name          string
		input         string
		expectedValid bool
	}{
		{
			name:          "Valid pair from cache without API call",
			input:         "XBTZAR",
			expectedValid: true,
		},
		{
			name:          "Valid pair after normalization without API call",
			input:         "btc-zar", // should normalize to XBTZAR which is in cache
			expectedValid: true,
		},
		{
			name:          "Bitcoin conversion without API call",
			input:         "BITCOINGBP", // should normalize to XBTGBP which is in cache
			expectedValid: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Since we don't want to rely on the API, we're just testing the cache lookup aspect
			normalizedPair := normalizeCurrencyPair(tc.input)
			isValid := validPairsCache[normalizedPair]

			if isValid != tc.expectedValid {
				t.Errorf("Cache validity for %q (normalized to %q) = %v, want %v",
					tc.input, normalizedPair, isValid, tc.expectedValid)
			}
		})
	}
}

// TestFindSimilarPairsForValidation tests the suggestions component
// of the validation system
func TestFindSimilarPairsForValidation(t *testing.T) {
	// Test base cases before API setup
	discoveredPairsCache = []string{"XBTZAR", "XBTGBP", "ETHZAR"}

	tests := []struct {
		name                 string
		input                string
		shouldHaveSuggestion string
	}{
		{
			name:                 "BTC pair should suggest XBT equivalent",
			input:                "BTCGBP",
			shouldHaveSuggestion: "XBTGBP",
		},
		{
			name:                 "Invalid pair should include known good pairs",
			input:                "INVALIDPAIR",
			shouldHaveSuggestion: "XBTZAR", // One of the known good pairs
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			suggestions := findSimilarPairs(tc.input)

			found := false
			for _, s := range suggestions {
				if s == tc.shouldHaveSuggestion {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("findSimilarPairs(%q) did not include expected suggestion %q in %v",
					tc.input, tc.shouldHaveSuggestion, suggestions)
			}
		})
	}
}
