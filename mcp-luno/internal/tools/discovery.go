package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/echarrod/mcp-luno/internal/config"
	"github.com/luno/luno-go"
)

// Cache of validated pairs to avoid repeated API calls
var (
	validPairsCache      map[string]bool
	discoveredPairsCache []string
)

// InitializePairDiscovery should be called at server startup to populate the valid pairs cache
func InitializePairDiscovery(ctx context.Context, cfg *config.Config) {
	// Initialize the cache maps
	validPairsCache = make(map[string]bool)
	discoveredPairsCache = []string{}

	// Populate with known working pairs first for immediate use
	for _, pair := range GetWorkingPairs() {
		validPairsCache[pair] = true
		discoveredPairsCache = append(discoveredPairsCache, pair)
	}

	// Launch discovery in background to avoid slowing down startup
	go func() {
		pairs := DiscoverAvailablePairs(context.Background(), cfg, false)
		fmt.Printf("Background pair discovery complete. Found %d valid pairs.\n", len(pairs))

		// Update cache with discovered pairs
		for _, pair := range pairs {
			if !validPairsCache[pair] {
				validPairsCache[pair] = true
				discoveredPairsCache = append(discoveredPairsCache, pair)
			}
		}
	}()
}

// DiscoverAvailablePairs attempts to find valid trading pairs on Luno by querying ticker API
// This is a discovery method we can use to automatically find valid pairs
func DiscoverAvailablePairs(ctx context.Context, cfg *config.Config, includeErrors bool) []string {
	baseCurrencies := []string{"XBT", "ETH", "XRP", "LTC", "BCH"}
	fiatCurrencies := []string{"ZAR", "NGN", "GBP", "EUR", "USD", "MYR", "IDR", "UGX"}
	cryptoBase := []string{"ETH", "XRP", "LTC", "BCH"}

	var validPairs []string
	var tryPairs []string

	// Generate base/fiat pairs
	for _, base := range baseCurrencies {
		for _, fiat := range fiatCurrencies {
			tryPairs = append(tryPairs, base+fiat)
		}
	}

	// Generate crypto/crypto pairs
	for _, coin := range cryptoBase {
		tryPairs = append(tryPairs, coin+"XBT")
	}

	fmt.Println("Attempting to discover valid trading pairs...")

	// Try each pair against the ticker API
	for _, pair := range tryPairs {
		ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{Pair: pair})
		if err == nil && ticker != nil {
			validPairs = append(validPairs, pair)
			fmt.Printf("Found valid pair: %s (Last trade price: %s)\n", pair, ticker.LastTrade.String())

			// Update the cache
			validPairsCache[pair] = true

			// Add to discovered pairs if not already there
			if !containsPair(discoveredPairsCache, pair) {
				discoveredPairsCache = append(discoveredPairsCache, pair)
			}
		} else if includeErrors {
			fmt.Printf("Invalid pair: %s (%v)\n", pair, err)
		}
	}

	fmt.Printf("Discovery complete. Found %d valid trading pairs.\n", len(validPairs))
	return validPairs
}

// Helper function to check if a slice contains a specific pair
func containsPair(pairs []string, pair string) bool {
	for _, p := range pairs {
		if p == pair {
			return true
		}
	}
	return false
}

// GetMarketInfo returns a detailed description of the market situation
func GetMarketInfo(ctx context.Context, cfg *config.Config, pair string) string {
	// First check if the pair is valid by trying to get ticker info
	ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{Pair: pair})
	if err != nil {
		return fmt.Sprintf("Could not get market info for %s: %v", pair, err)
	}

	orderBook, err := cfg.LunoClient.GetOrderBook(ctx, &luno.GetOrderBookRequest{Pair: pair})
	if err != nil {
		return fmt.Sprintf("Got ticker but could not get order book for %s: %v", pair, err)
	}

	var marketInfo strings.Builder

	marketInfo.WriteString(fmt.Sprintf("Market info for %s:\n", pair))
	marketInfo.WriteString(fmt.Sprintf("Last trade price: %s\n", ticker.LastTrade.String()))
	marketInfo.WriteString(fmt.Sprintf("Ask (Sell) price: %s\n", ticker.Ask.String()))
	marketInfo.WriteString(fmt.Sprintf("Bid (Buy) price: %s\n", ticker.Bid.String()))
	marketInfo.WriteString(fmt.Sprintf("24-hour volume: %s\n\n", ticker.Rolling24HourVolume.String()))

	// Add some order book info
	marketInfo.WriteString("Current Order Book:\n")
	if len(orderBook.Asks) > 0 {
		marketInfo.WriteString(fmt.Sprintf("Top 3 asks (Sell orders): \n"))
		for i := 0; i < 3 && i < len(orderBook.Asks); i++ {
			marketInfo.WriteString(fmt.Sprintf("  %s @ %s\n",
				orderBook.Asks[i].Volume.String(),
				orderBook.Asks[i].Price.String()))
		}
	}

	if len(orderBook.Bids) > 0 {
		marketInfo.WriteString(fmt.Sprintf("Top 3 bids (Buy orders): \n"))
		for i := 0; i < 3 && i < len(orderBook.Bids); i++ {
			marketInfo.WriteString(fmt.Sprintf("  %s @ %s\n",
				orderBook.Bids[i].Volume.String(),
				orderBook.Bids[i].Price.String()))
		}
	}

	return marketInfo.String()
}

// GetWorkingPairs returns a list of known working pairs based on testing
func GetWorkingPairs() []string {
	// Return cached pairs if available, or fallback to known working pairs
	if len(discoveredPairsCache) > 0 {
		return discoveredPairsCache
	}
	return []string{"XBTZAR", "ETHZAR", "XBTNGN", "XBTGBP", "XBTUSD", "ETHXBT"}
}

// ValidatePair checks if a trading pair is valid or returns suggestions if not
func ValidatePair(ctx context.Context, cfg *config.Config, pair string) (bool, string, []string) {
	// Normalize the pair first
	normalizedPair := normalizeCurrencyPair(pair)

	// Check if it's in our cache of valid pairs
	if validPairsCache[normalizedPair] {
		return true, "", nil
	}

	// Try to validate against the API
	ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{Pair: normalizedPair})
	if err == nil && ticker != nil {
		// It's valid, add to cache
		validPairsCache[normalizedPair] = true
		if !containsPair(discoveredPairsCache, normalizedPair) {
			discoveredPairsCache = append(discoveredPairsCache, normalizedPair)
		}
		return true, "", nil
	}

	// Not valid, find similar pairs for suggestions
	suggestions := findSimilarPairs(normalizedPair)
	errorMsg := fmt.Sprintf("Invalid trading pair: %s (%v)", normalizedPair, err)
	return false, errorMsg, suggestions
}

// findSimilarPairs finds trading pairs that are close to the given invalid pair
func findSimilarPairs(pair string) []string {
	var suggestions []string

	// Extract potential base currency
	// Assuming crypto codes are typically 3 chars
	var baseCurrency string
	if len(pair) >= 6 {
		baseCurrency = pair[:3]
		// We could extract the quote currency here if needed in the future
	} else if len(pair) >= 3 {
		baseCurrency = pair[:3]
	}

	// Check for common substitutions - especially BTC/XBT confusion
	if baseCurrency == "BTC" {
		for _, p := range GetWorkingPairs() {
			if strings.HasPrefix(p, "XBT") {
				suggestions = append(suggestions, p)
			}
		}
	} else if strings.HasPrefix(pair, "XBT") || strings.HasPrefix(pair, "ETH") {
		// Suggest pairs that start with the same base currency
		for _, p := range GetWorkingPairs() {
			if strings.HasPrefix(p, baseCurrency) {
				suggestions = append(suggestions, p)
			}
		}
	}

	// If no specific suggestions, return all working pairs
	if len(suggestions) == 0 {
		suggestions = GetWorkingPairs()
	}

	return suggestions
}
