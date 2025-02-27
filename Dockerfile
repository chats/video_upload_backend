# Dockerfile
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 go build -o /go/bin/api cmd/api/main.go

# Create final image
FROM alpine:3.16

# Install runtime dependencies
RUN apk add --no-cache ffmpeg ca-certificates

# Copy binary from builder
COPY --from=builder /go/bin/api /usr/local/bin/api

# Create necessary directories
RUN mkdir -p /app/tmp

# Set working directory
WORKDIR /app

# Copy migrations
COPY migrations ./migrations

# Set environment variables
ENV APP_ENV=production

# Expose port
EXPOSE 8080

# Run application
CMD ["api"]