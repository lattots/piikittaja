FROM golang:1.23.2-alpine

# Set the working directory inside the container
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY ./cmd/messenger ./cmd/messenger
COPY ./pkg ./pkg

# Build runnable binary from source
RUN go build -o /bin/messenger ./cmd/messenger/messenger.go

# Run the built binary
CMD ["messenger"]
