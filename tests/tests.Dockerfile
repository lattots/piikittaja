FROM golang:1.23.2-alpine AS test_base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
