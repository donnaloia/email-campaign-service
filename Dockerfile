# Start from the official Go image
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Start a new stage with a minimal image
FROM alpine:latest

# Add necessary certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port (adjust as needed)
EXPOSE 8080

# Run the binary
CMD ["./main"] 