FROM golang:1.23.2-alpine

# Set the working directory inside the container
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY ./cmd/admin_manager ./cmd/admin_manager
COPY ./pkg ./pkg

# Build runnable binary from source
RUN go build -o /bin/admin-manager ./cmd/admin_manager/manager.go

# Run the built binary
CMD ["admin-manager"]
