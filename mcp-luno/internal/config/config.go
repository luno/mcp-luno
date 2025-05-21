package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

const (
	// Environment variables
	EnvLunoAPIKeyID     = "LUNO_API_KEY_ID"
	EnvLunoAPIKeySecret = "LUNO_API_SECRET"
	EnvLunoAPIDomain    = "LUNO_API_DOMAIN"

	// Default Luno API domain
	DefaultLunoDomain = "api.luno.com"
)

// Config holds the configuration for the application
type Config struct {
	// Luno client
	LunoClient *luno.Client
}

// Mask a string to show only the first 4 characters and replace the rest with asterisks
func maskValue(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-4)
}

// Load loads the configuration from environment variables
func Load(domainOverride string) (*Config, error) {
	// Debugging: Print all environment variables to see if they're properly set
	// fmt.Println("*** Environment Variables Debug ***")
	// for _, env := range os.Environ() {
	// 	fmt.Println(env)
	// }
	// fmt.Println("*** End Environment Variables Debug ***")

	apiKeyID := os.Getenv(strings.TrimSpace(EnvLunoAPIKeyID))
	apiKeySecret := os.Getenv(strings.TrimSpace(EnvLunoAPIKeySecret))

	fmt.Printf("LUNO_API_KEY_ID value: %s (length: %d)\n", maskValue(apiKeyID), len(apiKeyID))
	fmt.Printf("LUNO_API_SECRET value: %s (length: %d)\n", maskValue(apiKeySecret), len(apiKeySecret))

	if apiKeyID == "" || apiKeySecret == "" {
		return nil, errors.New("luno API credentials not found, please set LUNO_API_KEY_ID and LUNO_API_SECRET environment variables")
	}

	// Set domain - first check command line override, then env var, then default
	domain := DefaultLunoDomain

	// Check for environment variable override
	if envDomain := os.Getenv(EnvLunoAPIDomain); envDomain != "" {
		domain = envDomain
		fmt.Printf("Using domain from environment variable: %s\n", domain)
	}

	// Command line override takes precedence if provided
	if domainOverride != "" {
		domain = domainOverride
		fmt.Printf("Using domain from command line: %s\n", domain)
	}

	// Create Luno client
	client := luno.NewClient()
	if domain != DefaultLunoDomain {
		client.SetBaseURL(fmt.Sprintf("https://%s", domain))
	}
	err := client.SetAuth(apiKeyID, apiKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to set Luno API credentials: %w", err)
	}
	client.SetDebug(true) // TODO: Remove when we don't need anymore
	return &Config{
		LunoClient: client,
	}, nil
}

// FormatCurrency formats a decimal amount with the currency code
func FormatCurrency(amount decimal.Decimal, currency string) string {
	return fmt.Sprintf("%s %s", amount.String(), strings.ToUpper(currency))
}
