.PHONY: build test clean run-stdio run-sse docker-build docker-run

# Binary name
BINARY_NAME=luno-mcp

# Build the application
build:
	go build -o $(BINARY_NAME) ./cmd/server

# Run tests
test:
	go test ./...

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

# Build Docker image
docker-build:
	docker build -t luno-mcp:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e LUNO_API_KEY_ID \
		-e LUNO_API_SECRET \
		luno-mcp:latest

# Install the binary to your GOBIN path
install:
	go install ./cmd/server

# Default target
default: build
