# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required packages
RUN apk add --no-cache \
    curl \
    postgresql-client

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy config directory
COPY --from=builder /app/config/. ./config/
# Copy wait-for script
COPY scripts/wait-for.sh /wait-for.sh
# Make script executable
RUN chmod +x /wait-for.sh

# Expose port
EXPOSE 8080

# Use wait-for script as entrypoint
ENTRYPOINT ["/wait-for.sh", "db", "./main"]