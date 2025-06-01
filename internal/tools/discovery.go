package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/luno/luno-go"
	"github.com/luno/luno-mcp/internal/config"
)

// GetMarketInfo returns a detailed description of the market situation
func GetMarketInfo(ctx context.Context, cfg *config.Config, pair string) (string, error) {
	// First check if the pair is valid by trying to get ticker info
	ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{Pair: pair})
	if err != nil {
		return "", fmt.Errorf("could not get market info for %s: %w", pair, err)
	}

	orderBook, err := cfg.LunoClient.GetOrderBook(ctx, &luno.GetOrderBookRequest{Pair: pair})
	if err != nil {
		return "", fmt.Errorf("got ticker but could not get order book for %s: %w", pair, err)
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
		marketInfo.WriteString("Top 3 asks (Sell orders): \n")
		for i := 0; i < 3 && i < len(orderBook.Asks); i++ {
			marketInfo.WriteString(fmt.Sprintf("  %s @ %s\n",
				orderBook.Asks[i].Volume.String(),
				orderBook.Asks[i].Price.String()))
		}
	}

	if len(orderBook.Bids) > 0 {
		marketInfo.WriteString("Top 3 bids (Buy orders): \n")
		for i := 0; i < 3 && i < len(orderBook.Bids); i++ {
			marketInfo.WriteString(fmt.Sprintf("  %s @ %s\n",
				orderBook.Bids[i].Volume.String(),
				orderBook.Bids[i].Price.String()))
		}
	}

	return marketInfo.String(), nil
}
