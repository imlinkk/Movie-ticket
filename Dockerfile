# Stage 1: Build the Go binary
FROM golang:alpine AS builder


WORKDIR /app

# Install git and build essentials
RUN apk add --no-cache git ca-certificates tzdata

# Copy source code and dependencies
COPY . .

# Download dependencies
RUN go mod download

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o movie-ticket-api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /app

# Install security certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/movie-ticket-api .

# Expose server port
EXPOSE 8080

# Run binary
ENTRYPOINT ["./movie-ticket-api"]
