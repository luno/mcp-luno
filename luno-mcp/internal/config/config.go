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
	
	// Default Luno API domain
	DefaultLunoDomain = "api.luno.com"
)

// Config holds the configuration for the application
type Config struct {
	// Luno API credentials
	LunoAPIKeyID     string
	LunoAPIKeySecret string
	
	// Luno client
	LunoClient *luno.Client
	
	// Luno API domain
	LunoDomain string
}

// Load loads the configuration from environment variables
func Load(domainOverride string) (*Config, error) {
	apiKeyID := os.Getenv(EnvLunoAPIKeyID)
	apiKeySecret := os.Getenv(EnvLunoAPIKeySecret)
	
	if apiKeyID == "" || apiKeySecret == "" {
		return nil, errors.New("Luno API credentials not found. Please set LUNO_API_KEY_ID and LUNO_API_SECRET environment variables")
	}
	
	// Set domain - use override if provided, otherwise use default
	domain := DefaultLunoDomain
	if domainOverride != "" {
		domain = domainOverride
	}
	
	// Create Luno client
	client := luno.NewClient(luno.WithApiKeyAuth(apiKeyID, apiKeySecret))
	if domain != DefaultLunoDomain {
		client.SetBaseURL(fmt.Sprintf("https://%s", domain))
	}
	
	return &Config{
		LunoAPIKeyID:     apiKeyID,
		LunoAPIKeySecret: apiKeySecret,
		LunoClient:       client,
		LunoDomain:       domain,
	}, nil
}

// FormatCurrency formats a decimal amount with the currency code
func FormatCurrency(amount decimal.Decimal, currency string) string {
	return fmt.Sprintf("%s %s", amount.String(), strings.ToUpper(currency))
}
