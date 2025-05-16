package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/echarrod/luno-mcp/internal/config"
	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool IDs
const (
	GetBalancesToolID       = "get_balances"
	GetTickerToolID         = "get_ticker"
	GetOrderBookToolID      = "get_order_book"
	CreateOrderToolID       = "create_order"
	CancelOrderToolID       = "cancel_order"
	ListOrdersToolID        = "list_orders"
	ListTransactionsToolID  = "list_transactions"
	GetTransactionToolID    = "get_transaction"
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
func HandleGetBalances(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
				Balance:     balance.Balance,
				Reserved:    balance.Reserved,
				Unconfirmed: balance.Unconfirmed,
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
			mcp.Description("Trading pair (e.g., XBTZAR)"),
		),
	)
}

// HandleGetTicker handles the get_ticker tool
func HandleGetTicker(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pair, ok := request.Params.Arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError("Trading pair is required"), nil
		}

		ticker, err := cfg.LunoClient.GetTicker(ctx, &luno.GetTickerRequest{
			Pair: pair,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get ticker: %v", err)), nil
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
			mcp.Description("Trading pair (e.g., XBTZAR)"),
		),
	)
}

// HandleGetOrderBook handles the get_order_book tool
func HandleGetOrderBook(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pair, ok := request.Params.Arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError("Trading pair is required"), nil
		}

		orderBook, err := cfg.LunoClient.GetOrderBook(ctx, &luno.GetOrderBookRequest{
			Pair: pair,
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get order book: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(orderBook, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal order book: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// ===== Trading Tools =====

// NewCreateOrderTool creates a new tool for creating orders
func NewCreateOrderTool() mcp.Tool {
	return mcp.NewTool(
		CreateOrderToolID,
		mcp.WithDescription("Create a new order"),
		mcp.WithString(
			"pair",
			mcp.Required(),
			mcp.Description("Trading pair (e.g., XBTZAR)"),
		),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Order type (BID or ASK)"),
			mcp.Enum("BID", "ASK"),
		),
		mcp.WithString(
			"price",
			mcp.Required(),
			mcp.Description("Order price"),
		),
		mcp.WithString(
			"volume",
			mcp.Required(),
			mcp.Description("Order volume"),
		),
	)
}

// HandleCreateOrder handles the create_order tool
func HandleCreateOrder(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract and validate arguments
		pair, ok := request.Params.Arguments["pair"].(string)
		if !ok || pair == "" {
			return mcp.NewToolResultError("Trading pair is required"), nil
		}

		orderType, ok := request.Params.Arguments["type"].(string)
		if !ok || (orderType != "BID" && orderType != "ASK") {
			return mcp.NewToolResultError("Order type must be 'BID' or 'ASK'"), nil
		}

		priceStr, ok := request.Params.Arguments["price"].(string)
		if !ok || priceStr == "" {
			return mcp.NewToolResultError("Order price is required"), nil
		}

		volumeStr, ok := request.Params.Arguments["volume"].(string)
		if !ok || volumeStr == "" {
			return mcp.NewToolResultError("Order volume is required"), nil
		}

		// Validate numeric values
		_, err := decimal.NewFromString(priceStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid price format: %v", err)), nil
		}

		_, err = decimal.NewFromString(volumeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid volume format: %v", err)), nil
		}

		// Create the order
		createReq := &luno.PostOrderRequest{
			Pair:   pair,
			Type:   orderType,
			Price:  priceStr,
			Volume: volumeStr,
		}

		order, err := cfg.LunoClient.PostOrder(ctx, createReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create order: %v", err)), nil
		}

		resultJSON, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal order: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
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
func HandleCancelOrder(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orderID, ok := request.Params.Arguments["order_id"].(string)
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
		mcp.WithInteger(
			"limit",
			mcp.Description("Maximum number of orders to return (default: 100)"),
		),
	)
}

// HandleListOrders handles the list_orders tool
func HandleListOrders(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var pair string
		if pairArg, ok := request.Params.Arguments["pair"]; ok {
			pair, _ = pairArg.(string)
		}

		limit := 100 // Default limit
		if limitArg, ok := request.Params.Arguments["limit"]; ok {
			if limitFloat, ok := limitArg.(float64); ok {
				limit = int(limitFloat)
				if limit <= 0 {
					limit = 100
				}
			}
		}

		listReq := &luno.ListOrdersRequest{
			Pair:  pair,
			Limit: strconv.Itoa(limit),
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
		mcp.WithInteger(
			"limit",
			mcp.Description("Maximum number of transactions to return (default: 100)"),
		),
		mcp.WithInteger(
			"min_row",
			mcp.Description("Minimum row ID to return (for pagination)"),
		),
		mcp.WithInteger(
			"max_row",
			mcp.Description("Maximum row ID to return (for pagination)"),
		),
	)
}

// HandleListTransactions handles the list_transactions tool
func HandleListTransactions(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountID, ok := request.Params.Arguments["account_id"].(string)
		if !ok || accountID == "" {
			return mcp.NewToolResultError("Account ID is required"), nil
		}

		listReq := &luno.ListTransactionsRequest{
			Id: accountID,
		}

		// Set limit if provided
		if limitArg, ok := request.Params.Arguments["limit"]; ok {
			if limitFloat, ok := limitArg.(float64); ok {
				limit := int(limitFloat)
				if limit > 0 {
					listReq.Limit = int64(limit)
				}
			}
		}

		// Set min_row if provided
		if minRowArg, ok := request.Params.Arguments["min_row"]; ok {
			if minRowFloat, ok := minRowArg.(float64); ok {
				minRow := int64(minRowFloat)
				if minRow > 0 {
					listReq.MinRow = minRow
				}
			}
		}

		// Set max_row if provided
		if maxRowArg, ok := request.Params.Arguments["max_row"]; ok {
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
func HandleGetTransaction(cfg *config.Config) mcp.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountID, ok := request.Params.Arguments["account_id"].(string)
		if !ok || accountID == "" {
			return mcp.NewToolResultError("Account ID is required"), nil
		}

		transactionID, ok := request.Params.Arguments["transaction_id"].(string)
		if !ok || transactionID == "" {
			return mcp.NewToolResultError("Transaction ID is required"), nil
		}

		// First, get the list of transactions
		listReq := &luno.ListTransactionsRequest{
			Id:    accountID,
			Limit: 100, // Use a reasonable limit to find the transaction
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
			return mcp.NewToolResultError(fmt.Sprintf("Transaction not found: %s", transactionID)), nil
		}

		resultJSON, err := json.MarshalIndent(transaction, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal transaction: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
