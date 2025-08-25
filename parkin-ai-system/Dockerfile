# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates for SSL/TLS verification
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for SSL/TLS verification
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Singapore

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy configuration files and resources
COPY --from=builder /app/manifest ./manifest
COPY --from=builder /app/resource ./resource

# Create logs directory
RUN mkdir -p /app/log

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/ || exit 1

# Run the application
CMD ["./main"]
