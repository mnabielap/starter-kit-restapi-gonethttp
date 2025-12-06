# --- Builder Stage ---
FROM golang:1.25.4-alpine AS builder

# Install git required for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is crucial for static linking (since we use the pure-go sqlite driver now)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# --- Final Stage ---
FROM alpine:latest

WORKDIR /usr/src/app

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create a directory for media uploads (persistence)
RUN mkdir -p /usr/src/app/media

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy entrypoint
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /usr/src/app
USER appuser

# Expose the port defined in .env.docker
EXPOSE 5005

ENTRYPOINT ["./entrypoint.sh"]

CMD ["./main"]