# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go build -o open-meet main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/open-meet .
EXPOSE 8080
CMD ["./open-meet"]
