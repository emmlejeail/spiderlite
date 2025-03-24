# Start from a lightweight Go image
FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV DD_AGENT_HOST=localhost
ENV DD_TRACE_AGENT_PORT=8126
ENV DD_ENV=prod
ENV DD_SERVICE=spiderlite

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build binary
COPY . .
RUN CGO_ENABLED=1 apk add --no-cache gcc musl-dev sqlite-dev && \
    go build -o spiderlite ./cmd/spiderlite

# Use a minimal runtime image
FROM alpine:latest
WORKDIR /app

# Install SQLite runtime dependencies
RUN apk add --no-cache sqlite-libs

# Copy the binary from builder
COPY --from=builder /app/spiderlite .

# Create data directory for SQLite database
RUN mkdir data && \
    chown -R nobody:nobody /app

# Use non-root user
USER nobody

# Default command (can be overridden)
ENTRYPOINT ["./spiderlite"]
