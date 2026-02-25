# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o taskmaster-api ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/taskmaster-api .

# Expose port
EXPOSE 8080

# Environment variables (defaults)
ENV PORT=8080

# Command to run
CMD ["./taskmaster-api"]
