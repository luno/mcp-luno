# Luno MCP Server

A [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server that provides access to the Luno cryptocurrency exchange API through the official Luno Go SDK.

This server enables integration with VS Code's Copilot and other MCP-compatible clients, providing contextual information and functionality related to the Luno cryptocurrency exchange.

## Features

- **Resources**: Access to account balances, transaction history, and more
- **Tools**: Functionality for creating and managing orders, checking prices, and viewing transaction details
- **Security**: Secure authentication using Luno API keys
- **VS Code Integration**: Easy integration with VS Code's Copilot features

## Installation

### Prerequisites

- Go 1.20 or later
- Luno account with API key and secret

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/echarrod/mcp-luno
cd mcp-luno
```

2. Build the binary:
```bash
go build -o luno-mcp ./cmd/server
```

3. Make it available system-wide (optional):
```bash
sudo mv luno-mcp /usr/local/bin/
```

## Usage

### Setting up credentials

The server requires your Luno API key and secret. These can be obtained from your Luno account settings.

Set the following environment variables:

```bash
export LUNO_API_KEY_ID=your_api_key_id
export LUNO_API_SECRET=your_api_secret
```

### Running the server

#### Standard I/O mode (default)

```bash
luno-mcp
```

#### Server-Sent Events (SSE) mode

```bash
luno-mcp --transport sse --sse-address localhost:8080
```

### Command-line options

- `--transport`: Transport type (`stdio` or `sse`, default: `stdio`)
- `--sse-address`: Address for SSE transport (default: `localhost:8080`)
- `--domain`: Luno API domain (default: `api.luno.com`)
- `--log-level`: Log level (`debug`, `info`, `warn`, `error`, default: `info`)

## VS Code Integration

To integrate with VS Code, add the following to your settings.json file:

### For stdio transport:

```json
"mcp": {
  "servers": {
    "luno": {
      "command": "mcp-luno",
      "args": [],
      "env": {
        "LUNO_API_KEY_ID": "your_api_key_id",
        "LUNO_API_SECRET": "your_api_secret"
      }
    }
  }
}
```

### For SSE transport:

```json
"mcp": {
  "servers": {
    "luno": {
      "type": "sse",
      "url": "http://localhost:8080/sse"
    }
  }
}
```

## Available Resources

- `luno://wallets`: List all wallets/balances in your Luno account
- `luno://transactions`: List recent transactions
- `luno://accounts/{id}`: Get details for a specific account by ID

## Available Tools

### Market Data
- `get_ticker`: Get current ticker information for a trading pair
- `get_order_book`: Get the order book for a trading pair

### Account Information
- `get_balances`: Get balances for all accounts

### Trading
- `create_order`: Create a new buy or sell order
- `cancel_order`: Cancel an existing order
- `list_orders`: List open orders

### Transactions
- `list_transactions`: List transactions for an account
- `get_transaction`: Get details of a specific transaction

## Examples

### Working with wallets

You can ask Copilot to show your wallet balances:
```
What are my current wallet balances on Luno?
```

### Trading

You can ask Copilot to help you trade:
```
Create a limit order to buy 0.001 BTC at 50000 ZAR
```

### Transaction history

You can ask Copilot to show your transaction history:
```
Show me my recent Bitcoin transactions
```

## Security Considerations

This tool requires API credentials that have access to your Luno account. Be cautious when using API keys, especially ones with withdrawal permissions. It's recommended to create API keys with only the permissions needed for your specific use case.

## License

[MIT License](LICENSE)

## Disclaimer

This software is not officially affiliated with or endorsed by Luno. Use at your own risk.
