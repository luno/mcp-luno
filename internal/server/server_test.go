package server

import (
	"context"
	"testing"

	"github.com/luno/luno-go"
	"github.com/luno/luno-mcp/internal/config"
	"github.com/mark3labs/mcp-go/mcp" // Added import
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/require"
)

const (
	testServerName       = "test-server"
	testServerWithHooks  = "test-server-with-hooks"
	testServerMultiHooks = "test-server-multi-hooks"
	testVersion1         = "1.0.0"
	testVersion2         = "1.0.1"
	testVersion3         = "1.0.2"
	testVersion4         = "0.0.1"
)

func TestNewMCPServer(t *testing.T) {
	tests := []struct {
		name    string
		srvName string
		version string
		hooks   []*mcpserver.Hooks
	}{
		{
			name:    "creates server without hooks",
			srvName: testServerName,
			version: testVersion1,
			hooks:   nil,
		},
		{
			name:    "creates server with single hook",
			srvName: testServerWithHooks,
			version: testVersion2,
			hooks: []*mcpserver.Hooks{
				func() *mcpserver.Hooks {
					h := &mcpserver.Hooks{}
					h.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
						// Hook logic would go here if there were specific actions to test
					})
					return h
				}(),
			},
		},
		{
			name:    "creates server with multiple distinct hook objects",
			srvName: testServerMultiHooks,
			version: testVersion3,
			hooks: []*mcpserver.Hooks{
				func() *mcpserver.Hooks { // Corresponds to original OnAnyHookFunc
					h := &mcpserver.Hooks{}
					h.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
						// Hook logic would go here if there were specific actions to test
					})
					return h
				}(),
				func() *mcpserver.Hooks { // Corresponds to original BeforeAnyHookFunc
					h := &mcpserver.Hooks{}
					h.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
						// Hook logic would go here if there were specific actions to test
					})
					return h
				}(),
				func() *mcpserver.Hooks { // Corresponds to original AfterAnyHookFunc, using AddOnSuccess for generality
					h := &mcpserver.Hooks{}
					h.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
						// Hook logic would go here if there were specific actions to test
					})
					return h
				}(),
				func() *mcpserver.Hooks { // Corresponds to original OnErrorHookFunc
					h := &mcpserver.Hooks{}
					h.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
						// Hook logic would go here if there were specific actions to test
					})
					return h
				}(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lunoClient := luno.NewClient()
			cfg := &config.Config{LunoClient: lunoClient}

			server := NewMCPServer(tc.srvName, tc.version, cfg, tc.hooks...)

			if server == nil {
				t.Error("NewMCPServer should return non-nil server")
			}
			// These should not panic
			registerResources(server, cfg)
			registerTools(server, cfg)
		})
	}
}

func TestServeSSEStructure(t *testing.T) {
	server := &mcpserver.MCPServer{}

	// This will fail due to invalid address, but it proves the function signature is correct
	// We just want to verify it compiles and can be called
	err := ServeSSE(context.Background(), server, "invalid:address")
	// We expect an error since we're providing an invalid address
	require.ErrorContains(t, err, "lookup tcp/address: unknown port")
}
