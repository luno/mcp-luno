FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o luno-mcp ./cmd/server

# Create a minimal final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/luno-mcp .

# Expose the default SSE port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/luno-mcp", "--transport", "sse", "--sse-address", "0.0.0.0:8080"]
