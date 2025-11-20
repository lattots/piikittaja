# Use Go version 1.20 on Alpine OS for building
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY ./cmd/web_app ./cmd/web_app
COPY ./pkg ./pkg

# Build runnable binary from source
RUN go build -o /bin/web-app ./cmd/web_app/server.go

# Final image for web app using Debian
FROM debian:13

RUN apt-get update && apt-get install -y ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /bin/web-app /usr/local/bin/web-app
COPY ./assets/web_app ./assets/web_app
