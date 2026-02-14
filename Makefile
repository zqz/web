.PHONY: help build test test-coverage run clean sqlc-generate migrate-up migrate-down migrate-create docker-up docker-down

# Default target
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  test             - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  run              - Run the application locally"
	@echo "  clean            - Clean build artifacts"
	@echo "  sqlc-generate    - Generate sqlc code from SQL queries"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback last database migration"
	@echo "  migrate-create   - Create a new migration (use NAME=migration_name)"
	@echo "  docker-up        - Start Docker containers for local development"
	@echo "  docker-down      - Stop Docker containers"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linters"
	@echo "  mod-tidy         - Tidy go modules"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server ./cmd/server

# Build for Linux ARM64 (e.g. for deploy). CGO_ENABLED=0 required when cross-compiling from macOS.
build-linux:
	@echo "Building for linux/arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/server-linux-arm64 ./cmd/server

# Run all tests
test:
	@echo "Running tests..."
	go test -v -race ./internal/... -timeout 300s

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./internal/... -timeout 300s
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/server

# Run the application with hot reload (requires air)
dev:
	@echo "Running application in dev mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not found. Install with: go install github.com/air-verse/air@latest"; \
		go run ./cmd/server; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate sqlc code
sqlc-generate:
	@echo "Generating sqlc code..."
	sqlc generate

# Run migrations
migrate-up:
	@echo "Running migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "ERROR: DATABASE_URL environment variable not set"; \
		exit 1; \
	fi
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "ERROR: DATABASE_URL environment variable not set"; \
		exit 1; \
	fi
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

# Create a new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "ERROR: NAME variable not set. Usage: make migrate-create NAME=create_users"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	goose -dir db/migrations create $(NAME) sql

# Start Docker containers for development
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofumpt -l -w .

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Run all pre-commit checks
pre-commit: fmt lint test
	@echo "All pre-commit checks passed!"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully!"
