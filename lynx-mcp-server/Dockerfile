# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install make
RUN apk add --no-cache make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux make build-server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/lynxmcpserver .

# Expose the port the server runs on
EXPOSE 9600

# Run the binary
CMD ["./lynxmcpserver"]
