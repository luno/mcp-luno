.PHONY: build test clean run-stdio run-sse

# Binary name
BINARY_NAME=mcp-luno

# Build the application
build:
	go build -o $(BINARY_NAME) ./cmd/server

# Run all tests
test:
	go test ./...

# Run unit tests only
test-unit:
	go test -v ./... -short

# Run integration tests (needs API credentials)
test-integration:
	go test -v ./internal/tests -run "Integration" -skip=""

# Clean build files
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run in stdio mode
run-stdio:
	go run ./cmd/server

# Run in SSE mode
run-sse:
	go run ./cmd/server --transport sse --sse-address localhost:8080

# Install the binary to your GOBIN path
install:
	go install ./cmd/server

pre-commit:
	pre-commit install

# Default target
default: build
