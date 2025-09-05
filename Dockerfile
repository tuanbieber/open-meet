# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o open-meet main.go

# Run stage
FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/open-meet .

# Expose app port (make sure this matches fly.toml internal_port)
EXPOSE 8080

# Run the app
CMD ["./open-meet"]
