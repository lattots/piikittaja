# Use Go version 1.20 on Alpine OS for building
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY ./cmd/telegram_bot ./cmd/telegram_bot
COPY ./pkg ./pkg

# Build runnable binary from source
RUN go build -o /bin/telegram-bot ./cmd/telegram_bot

# Final image for Telegram bot using Debian
FROM debian:12

RUN apt-get update && apt-get install -y ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /bin/telegram-bot /usr/local/bin/telegram-bot
COPY ./assets/telegram_bot ./assets/telegram_bot
