package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/echarrod/mcp-luno/internal/config"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Error messages
const (
	ErrAPICredentialsRequired = "API credentials are required for this operation. Please set LUNO_API_KEY_ID and LUNO_API_SECRET environment variables."
	ErrTradingPairRequired    = "Trading pair is required"
	ErrTradingPairDesc        = "Trading pair (e.g., XBTZAR)"
)

// Tool IDs
const (
	GetBalancesToolID      = "get_balances"
	GetTickerToolID        = "get_ticker"
	GetOrderBookToolID     = "get_order_book"
	CreateOrderToolID      = "create_order"
	CancelOrderToolID      = "cancel_order"
	ListOrdersToolID       = "list_orders"
	ListTransactionsToolID = "list_transactions"
	GetTransactionToolID   = "get_transaction"
	ListTradesToolID       = "list_trades"
)

// ===== Balance Tools =====

// NewGetBalancesTool creates a new tool for getting account balances
func NewGetBalancesTool() mcp.Tool {
	return mcp.NewTool(
		GetBalancesToolID,
		mcp.WithDescription("Get balances for all Luno accounts"),
	)
}

// HandleGetBalances handles the get_balances tool
func HandleGetBalances(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		balances, err := cfg.LunoClient.GetBalances(ctx, &luno.GetBalancesRequest{})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get balances: %v", err)), nil
		}

		// Enhance the response with additional information
		type EnhancedBalance struct {
			AccountID   string `json:"account_id"`
			Asset       string `json:"asset"`
			Balance     string `json:"balance"`
			Reserved    string `json:"reserved"`
			Unconfirmed string `json:"unconfirmed"`
			Name        string `json:"name"`
		}

		enhancedBalances := make([]EnhancedBalance, 0, len(balances.Balance))
		for _, balance := range balances.Balance {
			enhancedBalances = append(enhancedBalances, EnhancedBalance{
				AccountID:   balance.AccountId,
				Asset:       balance.Asset,
				Balance:     balance.Balance.String(),
				Reserved:    balance.Reserved.String(),
				Unconfirmed: balance.Unconfirmed.String(),
				Name:        balance.Name,
			})
		}

		resultJSON, err := json.MarshalIndent(enhancedBalances, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal balances: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Market Tools =====

// NewGetTickerTool creates a new tool for getting ticker information
func NewGetTickerTool() mcp.Tool {
	return mcp.NewTool(
		GetTickerToolID,
		mcp.WithDescription("Get ticker information for a trading pair"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description(ErrTradingPairDesc),
		),
	)
}

// HandleGetTicker handles the get_ticker tool
func HandleGetTicker(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.Params.Arguments
		pair, ok := arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError(ErrTradingPairRequired), nil
		}

		// Normalize currency pair
		pair = normalizeCurrencyPair(pair)

		ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{
			Pair: pair,
		})
		if err != nil {
			suggestedPairs := getSuggestedTradingPairs()
			errorMsg := fmt.Sprintf("Failed to get ticker: %v\n\nCommon trading pairs on Luno: %s",
				err, suggestedPairs)
			return mcp.NewToolResultError(errorMsg), nil
		}

		resultJSON, err := json.MarshalIndent(ticker, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal ticker: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// NewGetOrderBookTool creates a new tool for getting the order book
func NewGetOrderBookTool() mcp.Tool {
	return mcp.NewTool(
		GetOrderBookToolID,
		mcp.WithDescription("Get order book for a trading pair"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description(ErrTradingPairDesc),
		),
	)
}

// HandleGetOrderBook handles the get_order_book tool
func HandleGetOrderBook(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.Params.Arguments

		pair, ok := arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError(ErrTradingPairRequired), nil
		}

		// Normalize currency pair
		pair = normalizeCurrencyPair(pair)

		orderBook, err := cfg.LunoClient.GetOrderBook(ctx, &luno.GetOrderBookRequest{
			Pair: pair,
		})
		if err != nil {
			suggestedPairs := getSuggestedTradingPairs()
			errorMsg := fmt.Sprintf("Failed to get order book: %v\n\nCommon trading pairs on Luno: %s",
				err, suggestedPairs)
			return mcp.NewToolResultError(errorMsg), nil
		}

		resultJSON, err := json.MarshalIndent(orderBook, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal order book: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Trading Tools =====

// NewCreateOrderTool creates a new tool for creating limit orders
func NewCreateOrderTool() mcp.Tool {
	return mcp.NewTool(
		CreateOrderToolID,
		mcp.WithDescription("Create a new limit order"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description("Trading pair (e.g., XBTZAR)"),
		),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Order type (BUY or SELL)"),
			mcp.Enum("BUY", "SELL"),
		),
		mcp.WithString(
			"volume",
			mcp.Required(),
			mcp.Description("Order volume (amount of cryptocurrency to buy or sell)"),
		),
		mcp.WithString(
			"price",
			mcp.Required(),
			mcp.Description("Limit price as a decimal string"),
		),
	)
}

// HandleCreateOrder handles the create_order tool for limit orders
// TODO: Add HandleCreateMarketOrder function for market orders
func HandleCreateOrder(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		arguments := request.Params.Arguments
		// Extract and validate arguments
		pair, ok := arguments["pair"].(string)
		if !ok || pair == "" {
			// If no pair is provided, return a list of working pairs
			workingPairs := GetWorkingPairs()
			errorMsg := fmt.Sprintf("Trading pair is required. Please use one of these known working pairs: %s",
				strings.Join(workingPairs, ", "))
			return mcp.NewToolResultError(errorMsg), nil
		}

		// Log original pair for debugging
		fmt.Printf("Original pair before normalization: %s\n", pair)

		// Normalize the pair - this should handle BTC->XBT conversion automatically
		pair = normalizeCurrencyPair(pair)
		fmt.Printf("Final normalized pair: %s\n", pair)

		// Validate the trading pair with our improved validation function
		isValid, errorMsg, suggestions := ValidatePair(ctx, cfg, pair)
		if !isValid {
			// If invalid, show a helpful error message with suggestions
			suggestionsStr := strings.Join(suggestions, ", ")
			pairErrorMsg := fmt.Sprintf("Invalid trading pair: %s\n\n%s\n\nPlease try one of these working pairs: %s",
				pair, errorMsg, suggestionsStr)
			return mcp.NewToolResultError(pairErrorMsg), nil
		}

		orderType, ok := arguments["type"].(string)
		if !ok || (orderType != "BUY" && orderType != "SELL") {
			return mcp.NewToolResultError("Order type must be 'BUY' or 'SELL'"), nil
		}

		volumeStr, ok := arguments["volume"].(string)
		if !ok || volumeStr == "" {
			return mcp.NewToolResultError("Order volume is required"), nil
		}

		priceStr, ok := arguments["price"].(string)
		if !ok || priceStr == "" {
			return mcp.NewToolResultError("Limit price is required"), nil
		}

		// Validate numeric values
		volumeDec, err := decimal.NewFromString(volumeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid volume format: %v", err)), nil
		}

		priceDec, err := decimal.NewFromString(priceStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid price format: %v", err)), nil
		}

		// Map BUY/SELL to BID/ASK for limit orders
		var lunoOrderType luno.OrderType
		if orderType == "BUY" {
			lunoOrderType = luno.OrderTypeBid
		} else { // SELL
			lunoOrderType = luno.OrderTypeAsk
		}

		// Get market info - we already validated the pair, but this provides additional info
		marketInfo := GetMarketInfo(ctx, cfg, pair)
		fmt.Println(marketInfo)

		// Log the request parameters for debugging
		fmt.Printf("Creating order with parameters: Pair=%s, Type=%s, Volume=%s, Price=%s\n",
			pair, lunoOrderType, volumeDec.String(), priceDec.String())

		// Create the limit order
		createReq := &luno.PostLimitOrderRequest{
			Pair:   pair,
			Type:   lunoOrderType,
			Volume: volumeDec,
			Price:  priceDec,
		}

		order, err := cfg.LunoClient.PostLimitOrder(ctx, createReq)
		if err != nil {
			// If the order fails despite our validation, provide detailed error information
			errorMsg := fmt.Sprintf("Failed to create limit order: %v\n\n"+
				"Here's what we know about this market:\n%s\n\n"+
				"This may be due to insufficient balance, market conditions, or API limits.",
				err, marketInfo)

			return mcp.NewToolResultError(errorMsg), nil
		}

		// Order succeeded
		resultJSON, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal order result: %v", err)), nil
		}

		successMsg := fmt.Sprintf("Order created successfully!\n\n%s\n\n%s",
			string(resultJSON), marketInfo)
		return mcp.NewToolResultText(successMsg), nil
	}
}

// NewCancelOrderTool creates a new tool for canceling orders
func NewCancelOrderTool() mcp.Tool {
	return mcp.NewTool(
		CancelOrderToolID,
		mcp.WithDescription("Cancel an order"),
		mcp.WithString(
			"order_id",
			mcp.Required(),
			mcp.Description("Order ID to cancel"),
		),
	)
}

// HandleCancelOrder handles the cancel_order tool
func HandleCancelOrder(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		arguments := request.Params.Arguments
		orderID, ok := arguments["order_id"].(string)
		if !ok || orderID == "" {
			return mcp.NewToolResultError("Order ID is required"), nil
		}

		result, err := cfg.LunoClient.StopOrder(ctx, &luno.StopOrderRequest{
			OrderId: orderID,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to cancel order: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// NewListOrdersTool creates a new tool for listing orders
func NewListOrdersTool() mcp.Tool {
	return mcp.NewTool(
		ListOrdersToolID,
		mcp.WithDescription("List open orders"),
		mcp.WithString(
			"pair",
			mcp.Description("Trading pair (e.g., XBTZAR)"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Maximum number of orders to return (default: 100)"),
		),
	)
}

// HandleListOrders handles the list_orders tool
func HandleListOrders(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		arguments := request.Params.Arguments
		var pair string
		if pairArg, ok := arguments["pair"]; ok {
			pair, _ = pairArg.(string)
		}

		limit := 100 // Default limit
		if limitArg, ok := arguments["limit"]; ok {
			if limitFloat, ok := limitArg.(float64); ok {
				limit = int(limitFloat)
				if limit <= 0 {
					limit = 100
				}
			}
		}

		listReq := &luno.ListOrdersRequest{
			Pair:  pair,
			Limit: int64(limit),
		}

		orders, err := cfg.LunoClient.ListOrders(ctx, listReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list orders: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(orders, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal orders: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Transaction Tools =====

// NewListTransactionsTool creates a new tool for listing transactions
func NewListTransactionsTool() mcp.Tool {
	return mcp.NewTool(
		ListTransactionsToolID,
		mcp.WithDescription("List transactions for an account"),
		mcp.WithString(
			"account_id",
			mcp.Required(),
			mcp.Description("Account ID"),
		),
		mcp.WithNumber(
			"min_row",
			mcp.Description("Minimum row ID to return (for pagination, inclusive)"),
		),
		mcp.WithNumber(
			"max_row",
			mcp.Description("Maximum row ID to return (for pagination, exclusive)"),
		),
	)
}

// HandleListTransactions handles the list_transactions tool
func HandleListTransactions(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		arguments := request.Params.Arguments
		accountIDStr, ok := arguments["account_id"].(string)
		if !ok || accountIDStr == "" {
			return mcp.NewToolResultError("Account ID is required"), nil
		}

		// Convert account ID from string to int64
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid account ID format: %v. Please provide a valid numeric account ID.", err)), nil
		}

		listReq := &luno.ListTransactionsRequest{
			Id: accountID,
		}

		// Set default rows
		listReq.MaxRow = 100 // Default max rows
		listReq.MinRow = 1   // Default min rows

		// Set min_row if provided
		if minRowArg, ok := arguments["min_row"]; ok {
			if minRowFloat, ok := minRowArg.(float64); ok {
				minRow := int64(minRowFloat)
				if minRow >= 0 {
					listReq.MinRow = minRow
				}
			}
		}

		// Set max_row if provided
		if maxRowArg, ok := arguments["max_row"]; ok {
			if maxRowFloat, ok := maxRowArg.(float64); ok {
				maxRow := int64(maxRowFloat)
				if maxRow > 0 {
					listReq.MaxRow = maxRow
				}
			}
		}

		transactions, err := cfg.LunoClient.ListTransactions(ctx, listReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list transactions: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(transactions, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal transactions: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// NewGetTransactionTool creates a new tool for getting a specific transaction
func NewGetTransactionTool() mcp.Tool {
	return mcp.NewTool(
		GetTransactionToolID,
		mcp.WithDescription("Get details of a specific transaction"),
		mcp.WithString(
			"account_id",
			mcp.Required(),
			mcp.Description("Account ID"),
		),
		mcp.WithString(
			"transaction_id",
			mcp.Required(),
			mcp.Description("Transaction ID"),
		),
	)
}

// HandleGetTransaction handles the get_transaction tool
func HandleGetTransaction(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Since we're using a private API endpoint, authentication errors will be handled by the API call

		arguments := request.Params.Arguments
		accountIDStr, ok := arguments["account_id"].(string)
		if !ok || accountIDStr == "" {
			return mcp.NewToolResultError("Account ID is required"), nil
		}

		// Convert account ID from string to int64
		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid account ID format: %v. Please provide a valid numeric account ID.", err)), nil
		}

		transactionIDStr, ok := arguments["transaction_id"].(string)
		if !ok || transactionIDStr == "" {
			return mcp.NewToolResultError("Transaction ID is required"), nil
		}

		// Attempt to convert transaction ID to int64 for comparison
		transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid transaction ID format: %v. Please provide a valid numeric transaction ID.", err)), nil
		}

		// Get the list of transactions with MinRow and MaxRow
		listReq := &luno.ListTransactionsRequest{
			Id:     accountID,
			MinRow: 0,    // Start from the beginning
			MaxRow: 1000, // Use a reasonable max to find the transaction
		}

		transactions, err := cfg.LunoClient.ListTransactions(ctx, listReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get transactions: %v", err)), nil
		}

		// Find the specific transaction
		var transaction *luno.Transaction
		for _, txn := range transactions.Transactions {
			if txn.RowIndex == transactionID {
				transaction = &txn
				break
			}
		}

		if transaction == nil {
			return mcp.NewToolResultError(fmt.Sprintf("Transaction not found: %s", transactionIDStr)), nil
		}

		resultJSON, err := json.MarshalIndent(transaction, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal transaction: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Trades Tools =====

// NewListTradesTool creates a new tool for listing trades
func NewListTradesTool() mcp.Tool {
	return mcp.NewTool(
		ListTradesToolID,
		mcp.WithDescription("List recent trades for a currency pair"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description(ErrTradingPairDesc),
		),
		mcp.WithString(
			"since",
			mcp.Description("Fetch trades executed after this timestamp (Unix milliseconds)"),
		),
	)
}

// HandleListTrades handles the list_trades tool
func HandleListTrades(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments := request.Params.Arguments
		pair, ok := arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError(ErrTradingPairRequired), nil
		}

		// Normalize currency pair
		pair = normalizeCurrencyPair(pair)

		req := &luno.ListTradesRequest{
			Pair: pair,
		}

		// Set since timestamp if provided
		if sinceStr, ok := arguments["since"].(string); ok && sinceStr != "" {
			// Parse the timestamp string to int64
			sinceTimestamp, err := strconv.ParseInt(sinceStr, 10, 64)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid since timestamp format: %v", err)), nil
			}
			// Convert to time.Time then to luno.Time
			t := time.Unix(0, sinceTimestamp*1e6)
			req.Since = luno.Time(t)
		}

		trades, err := cfg.LunoClient.ListTrades(ctx, req)
		if err != nil {
			suggestedPairs := getSuggestedTradingPairs()
			errorMsg := fmt.Sprintf("Failed to list trades: %v\n\nCommon trading pairs on Luno: %s",
				err, suggestedPairs)
			return mcp.NewToolResultError(errorMsg), nil
		}

		resultJSON, err := json.MarshalIndent(trades, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal trades: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Helper Functions =====

// normalizeCurrencyPair converts common currency pair formats to Luno's expected format
func normalizeCurrencyPair(pair string) string {
	// Log input for debugging
	originalPair := pair

	// Remove any separators that might be in the pair
	pair = strings.Replace(pair, "-", "", -1)
	pair = strings.Replace(pair, "_", "", -1)
	pair = strings.Replace(pair, "/", "", -1)
	pair = strings.ToUpper(pair)

	// Apply currency code standardization
	// Known mappings between common symbols and Luno's expected format
	currencyMappings := map[string]string{
		"BTC":     "XBT", // Bitcoin is XBT on Luno
		"BITCOIN": "XBT",
		// Add other mappings if needed in the future
	}

	// Apply all mappings
	for common, luno := range currencyMappings {
		pair = strings.Replace(pair, common, luno, -1)
	}

	// Log the normalization for debugging
	fmt.Printf("Normalized pair: %s -> %s\n", originalPair, pair)

	return pair
}

// getSuggestedTradingPairs returns a string with common Luno trading pairs
func getSuggestedTradingPairs() string {
	// Dynamic generation of trading pairs
	// Define base currencies and fiat currencies
	baseCurrencies := []string{"XBT", "ETH", "XRP", "LTC", "BCH"}
	fiatCurrencies := []string{"ZAR", "NGN", "GBP", "EUR", "USD", "MYR", "IDR", "UGX"}

	// Generate all possible combinations
	var pairs []string

	// Base/fiat pairs (most common)
	for _, base := range baseCurrencies {
		for _, fiat := range fiatCurrencies {
			// Not all combinations exist, but including most common ones
			if base == "XBT" || // XBT pairs with all fiats
				(base == "ETH" && (fiat == "ZAR" || fiat == "NGN" || fiat == "GBP" ||
					fiat == "EUR" || fiat == "USD" || fiat == "MYR" || fiat == "IDR")) {
				pairs = append(pairs, base+fiat)
			}
		}
	}

	// Add crypto-to-crypto pairs
	cryptoCryptoPairs := []string{"ETHXBT", "XRPXBT", "LTCXBT", "BCHXBT"}
	pairs = append(pairs, cryptoCryptoPairs...)

	return strings.Join(pairs, ", ")
}
