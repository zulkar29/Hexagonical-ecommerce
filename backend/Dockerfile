# Multi-stage build for Go backend
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Set proper port
EXPOSE 8080

# Add proper command
CMD ["./main"]