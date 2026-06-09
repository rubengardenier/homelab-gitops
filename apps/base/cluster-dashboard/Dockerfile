# Build stage for Go backend
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod and source files
COPY backend/go.mod ./
COPY backend/*.go ./

# Download dependencies
RUN go mod tidy && go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final stage
FROM alpine:3.24

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy static files
COPY src/ ./static/

# Create non-root user
RUN adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

# Environment variables
ENV PORT=8080
ENV STATIC_DIR=/app/static

EXPOSE 8080

CMD ["./server"]
