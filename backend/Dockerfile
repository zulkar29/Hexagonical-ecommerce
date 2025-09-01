# TODO: Implement actual Dockerfile during development
FROM golang:1.21-alpine AS builder

WORKDIR /app

# TODO: Copy go.mod and go.sum files
# COPY go.mod go.sum ./
# RUN go mod download

# TODO: Copy source code and build
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# TODO: Copy binary from builder stage
# COPY --from=builder /app/main .

# TODO: Set proper port
# EXPOSE 8080

# TODO: Add proper command
# CMD ["./main"]