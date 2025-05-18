package main

import (
	"fmt"
)

// normalizeCurrencyPair converts common currency pair formats to Luno's expected format
func normalizeCurrencyPair(pair string) string {
	// Remove any separators that might be in the pair
	pair = strings.Replace(pair, "-", "", -1)
	pair = strings.Replace(pair, "_", "", -1)
	pair = strings.Replace(pair, "/", "", -1)
	pair = strings.ToUpper(pair)

	// Map common symbols to Luno-specific symbols
	mappings := map[string]string{
		"BTC":    "XBT", // Bitcoin is XBT on Luno
		"BTCZAR": "XBTZAR",
		"BTCNGN": "XBTNGN",
		"BTCMYR": "XBTMYR",
		"BTCIDR": "XBTIDR",
		"BTCUGX": "XBTUGX",
		"BTCZMT": "XBTZMT",
		"BTCGBP": "XBTGBP",
		"BTCEUR": "XBTEUR",
		"BTCUSD": "XBTUSD",
	}

	if mapped, exists := mappings[pair]; exists {
		return mapped
	}

	return pair
}

func main() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"BTC", "XBT"},
		{"BTCGBP", "XBTGBP"},
		{"BTC-GBP", "XBTGBP"},
		{"BTC/GBP", "XBTGBP"},
		{"BTC_GBP", "XBTGBP"},
		{"btcgbp", "XBTGBP"},
		{"xbTGbP", "XBTGBP"},
		{"ETHZAR", "ETHZAR"},
		{"ETH-ZAR", "ETHZAR"},
	}

	for _, tc := range testCases {
		result := normalizeCurrencyPair(tc.input)
		fmt.Printf("Input: %s, Expected: %s, Got: %s, Match: %v\n",
			tc.input, tc.expected, result, result == tc.expected)
	}
}

var strings = struct {
	Replace func(string, string, string, int) string
	ToUpper func(string) string
}{
	Replace: func(s, old, new string, n int) string {
		return s
	},
	ToUpper: func(s string) string {
		return s
	},
}
