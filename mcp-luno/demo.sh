#!/bin/bash
# Demo script to showcase Luno MCP Server functionality

# Check if the API credentials are set
if [ -z "$LUNO_API_KEY_ID" ] || [ -z "$LUNO_API_SECRET" ]; then
  echo "Error: Luno API credentials not set"
  echo "Please set LUNO_API_KEY_ID and LUNO_API_SECRET environment variables"
  echo "Example:"
  echo "  export LUNO_API_KEY_ID=your_api_key_id"
  echo "  export LUNO_API_SECRET=your_api_secret"
  exit 1
fi

# Print header
echo "================================================"
echo "Luno MCP Server Demo"
echo "================================================"
echo ""

# Build the server if not already built
echo "Building Luno MCP Server..."
go build -o mcp-luno ./cmd/server
echo "Build complete!"
echo ""

# Start the server in the background (SSE mode)
echo "Starting Luno MCP Server in SSE mode..."
./mcp-luno --transport sse --sse-address localhost:8080 --log-level debug &
SERVER_PID=$!

# Wait for server to start
sleep 2
echo "Server started with PID: $SERVER_PID"
echo ""

# Print instructions
echo "The server is now running in SSE mode on localhost:8080"
echo ""
echo "You can integrate it with VS Code by adding the following to your settings.json:"
echo ""
echo '{
  "mcp": {
    "servers": {
      "luno": {
        "type": "sse",
        "url": "http://localhost:8080/sse"
      }
    }
  }
}'
echo ""
echo "Alternatively, for stdio mode, use:"
echo ""
echo '{
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
}'
echo ""

# Wait for user to press a key
echo "Press any key to stop the server and exit..."
read -n 1 -s

# Kill the server
echo "Stopping server..."
kill $SERVER_PID
echo "Server stopped"
echo ""

echo "Thank you for trying the Luno MCP Server!"
