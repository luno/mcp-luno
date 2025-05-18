package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/echarrod/mcp-luno/internal/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Tool IDs
const (
	ValidatePairToolID = "validate_pair"
)

// NewValidatePairTool creates a new tool for validating trading pairs without making trades
func NewValidatePairTool() mcp.Tool {
	return mcp.NewTool(
		ValidatePairToolID,
		mcp.WithDescription("Validate a trading pair without creating an order"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description("Trading pair to validate (e.g., BTCGBP, XBT-ZAR)"),
		),
	)
}

// HandleValidatePair handles the validate_pair tool
func HandleValidatePair(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.Params.Arguments
		pair, ok := arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError(ErrTradingPairRequired), nil
		}

		// Log original pair for debugging
		originalPair := pair

		// Normalize the pair
		pair = normalizeCurrencyPair(pair)

		// Validate the pair
		isValid, errorMsg, suggestions := ValidatePair(ctx, cfg, pair)

		type ValidationResult struct {
			OriginalPair   string   `json:"original_pair"`
			NormalizedPair string   `json:"normalized_pair"`
			IsValid        bool     `json:"is_valid"`
			Message        string   `json:"message"`
			Suggestions    []string `json:"suggestions,omitempty"`
		}

		var result ValidationResult

		if isValid {
			// We gather market info for display in the response message
			// No need to store it in a variable here since we're using it directly in the response

			result = ValidationResult{
				OriginalPair:   originalPair,
				NormalizedPair: pair,
				IsValid:        true,
				Message:        fmt.Sprintf("Trading pair '%s' is valid. Original input: '%s'", pair, originalPair),
				Suggestions:    nil,
			}
		} else {
			result = ValidationResult{
				OriginalPair:   originalPair,
				NormalizedPair: pair,
				IsValid:        false,
				Message:        errorMsg,
				Suggestions:    suggestions,
			}
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
		}

		// Format the message nicely
		var respMessage strings.Builder
		if isValid {
			respMessage.WriteString(fmt.Sprintf("✅ Valid trading pair: %s\n\n", pair))

			// Add market info
			marketInfo := GetMarketInfo(ctx, cfg, pair)
			respMessage.WriteString(marketInfo)

		} else {
			respMessage.WriteString(fmt.Sprintf("❌ Invalid trading pair: %s\n\n", pair))
			respMessage.WriteString(fmt.Sprintf("Error: %s\n\n", errorMsg))
			respMessage.WriteString("Suggestions:\n")
			for _, s := range suggestions {
				respMessage.WriteString(fmt.Sprintf("- %s\n", s))
			}
			respMessage.WriteString("\nNote: Luno uses XBT instead of BTC for Bitcoin.")
		}

		// Also add the raw JSON for programmatic access
		respMessage.WriteString("\n\nRaw data:\n")
		respMessage.WriteString(string(resultJSON))

		return mcp.NewToolResultText(respMessage.String()), nil
	}
}
