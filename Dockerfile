# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o open-meet cmd/*.go

# Run stage
FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/open-meet .

# Copy and rename .env.example to .env
COPY .env.example .env

# Expose port 8080
EXPOSE 8080

# Run the app
CMD ["./open-meet"]
