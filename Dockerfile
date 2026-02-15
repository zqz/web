# Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /build

# Copy module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# Install goose for migrations (version from go.mod)
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

# Run stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy server binary and goose from builder
COPY --from=builder /build/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy templates, static assets, and migrations
COPY --from=builder /build/templates ./templates
COPY --from=builder /build/static ./static
COPY --from=builder /build/db/migrations ./db/migrations

# Optional: run as non-root (create app user)
RUN adduser -D -g "" appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

# Default: run the server. Override with goose for migrations, e.g.:
#   docker run ... goose -dir db/migrations postgres "$DATABASE_URL" up
CMD ["./server"]
