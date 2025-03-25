# Start from a lightweight Go image
FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV DD_AGENT_HOST=localhost
ENV DD_TRACE_AGENT_PORT=8126
ENV DD_ENV=prod
ENV DD_SERVICE=spiderlite

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

# Create data directory
RUN mkdir /data && chmod 777 /data

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=1 go build -o spiderlite ./cmd/crawler
RUN CGO_ENABLED=1 go build -o spiderlite-server ./cmd/server

EXPOSE 8080

# Default command (can be overridden)
CMD ["./spiderlite-server"]
