FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git make

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application for Linux architecture
RUN make build-for-linux-amd64

# Final stage
FROM gcr.io/distroless/static:nonroot

# Set metadata
LABEL org.opencontainers.image.source="https://github.com/cnosuke/mcp-search"
LABEL org.opencontainers.image.description="MCP server for web search using Brave Search API"

WORKDIR /app

# Copy configuration
COPY --from=builder /app/config.yml /app/config.yml

# Copy the binary
COPY --from=builder /app/bin/mcp-search-linux-amd64 /app/mcp-search

# Use nonroot user
USER nonroot:nonroot

# Set the entrypoint
ENTRYPOINT ["/app/mcp-search"]

# Default command
CMD ["server", "--config", "config.yml"]
