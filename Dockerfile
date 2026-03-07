# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o cbs-simulator .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /app/cbs-simulator .
COPY --from=builder /app/database ./database
COPY --from=builder /app/.env.example ./.env

# Create database directory
RUN mkdir -p /app/database

# Expose port
EXPOSE 8080

# Run the application
CMD ["./cbs-simulator"]
