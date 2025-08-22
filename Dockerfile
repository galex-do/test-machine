# Use official Golang image as base
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies for PostgreSQL driver
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the application (enable CGO for PostgreSQL driver)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Use minimal alpine image for final stage
FROM alpine:latest

# Install ca-certificates and netcat for database connection check
RUN apk --no-cache add ca-certificates netcat-openbsd

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files, templates, and migrations
COPY static/ ./static/
COPY templates/ ./templates/
COPY migrations/ ./migrations/

# Expose port 5000
EXPOSE 5000

# Command to run the application
CMD ["./main"]