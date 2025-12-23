FROM golang:1.25.1-alpine3.21 AS builder

# Install git and ca-certificates (needed for go get if any)
RUN apk update && apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Install tools
RUN go install github.com/githubnemo/CompileDaemon@latest

# Copy source code
COPY . .

# BadgerDB data location
VOLUME ["/app/data"]

# Proper CompileDaemon command
CMD ["sh", "-c", "go build -o main . && ./main"]
