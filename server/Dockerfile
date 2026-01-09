# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO/SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o /quickpic-server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Create data directory
RUN mkdir -p /data

# Copy binary from builder
COPY --from=builder /quickpic-server .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./quickpic-server"]
