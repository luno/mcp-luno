package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/luno/luno-go"
)

// Only in place to confirm we can access API through the SDK with env vars
func TestListOrders_SDKIntegration(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Verify API key environment variables are set
	apiKeyID := os.Getenv("LUNO_API_KEY_ID")
	apiKeySecret := os.Getenv("LUNO_API_SECRET")
	apiDomain := os.Getenv("LUNO_API_DOMAIN")

	if apiKeyID == "" || apiKeySecret == "" {
		t.Fatalf("API credentials not set in environment. LUNO_API_KEY_ID and LUNO_API_SECRET must be set.")
	}

	fmt.Printf("API Key ID: %s (length: %d)\n", maskValue(apiKeyID), len(apiKeyID))
	fmt.Printf("API Secret: %s (length: %d)\n", maskValue(apiKeySecret), len(apiKeySecret))

	if apiDomain != "" {
		fmt.Printf("Using custom API domain: %s\n", apiDomain)
	} else {
		fmt.Println("Using default API domain (api.luno.com)")
	}

	// Create Luno client with the API credentials
	client := luno.NewClient()
	client.SetAuth(apiKeyID, apiKeySecret)

	// Set custom API domain if provided
	if apiDomain != "" {
		client.SetBaseURL(fmt.Sprintf("https://%s", apiDomain))
	}

	// Try to list orders
	ctx := context.Background()
	listReq := &luno.ListOrdersRequest{
		// No parameters to keep it simple, just test the API connection
	}

	orders, err := client.ListOrders(ctx, listReq)
	if err != nil {
		t.Fatalf("Failed to list orders: %v", err)
	}

	// Verify the response
	if orders == nil {
		t.Fatalf("Orders response is nil")
	}

	// Print order count
	fmt.Printf("Successfully retrieved %d orders\n", len(orders.Orders))

	// Print first order details if available
	if len(orders.Orders) > 0 {
		order := orders.Orders[0]
		fmt.Printf("Sample order: OrderID=%s, Type=%s, State=%s\n",
			order.OrderId,
			order.Type,
			order.State)
	} else {
		fmt.Println("No orders found in account")
	}

	// Test passes if we get here
	fmt.Println("API connection test successful")
}

// Mask a string to show only the first 4 characters and replace the rest with asterisks
func maskValue(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-4)
}
