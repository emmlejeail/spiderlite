# Start from a lightweight Go image
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build binary
COPY . .
RUN go build -o crawler ./main.go

# Use a minimal runtime image
FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/crawler .

# Use non-root user (optional security)
RUN adduser -D appuser
USER appuser

# Default command (can be overridden)
ENTRYPOINT ["./crawler"]
