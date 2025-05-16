.PHONY: build test clean run-stdio run-sse docker-build docker-run

# Binary name
BINARY_NAME=mcp-luno

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
	docker build -t mcp-luno:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e LUNO_API_KEY_ID \
		-e LUNO_API_SECRET \
		mcp-luno:latest

# Install the binary to your GOBIN path
install:
	go install ./cmd/server

pre-commit:
	pre-commit install

# Default target
default: build
