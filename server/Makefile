.PHONY: build run test test-e2e test-e2e-verbose clean migrate docker-up docker-down

# Build the server
build:
	go build -o bin/server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# Run unit tests
test:
	go test -v ./internal/...

# Run e2e tests
test-e2e:
	@echo "Running e2e tests..."
	@echo "Make sure PostgreSQL is running and quickpic_test database exists"
	go test ./tests/e2e -count=1

# Run e2e tests with verbose output
test-e2e-verbose:
	@echo "Running e2e tests (verbose)..."
	go test -v ./tests/e2e -count=1

# Run all tests
test-all: test test-e2e

# Clean build artifacts
clean:
	rm -rf bin/

# Create test database
create-test-db:
	createdb quickpic_test || true

# Drop test database
drop-test-db:
	dropdb quickpic_test || true

# Reset test database
reset-test-db: drop-test-db create-test-db

# Start dependencies with docker-compose
docker-up:
	docker-compose up -d

# Stop dependencies
docker-down:
	docker-compose down

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the server binary"
	@echo "  run              - Run the server"
	@echo "  test             - Run unit tests"
	@echo "  test-e2e         - Run e2e tests"
	@echo "  test-e2e-verbose - Run e2e tests with verbose output"
	@echo "  test-all         - Run all tests"
	@echo "  create-test-db   - Create the test database"
	@echo "  reset-test-db    - Reset the test database"
	@echo "  docker-up        - Start docker dependencies"
	@echo "  docker-down      - Stop docker dependencies"
	@echo "  dev              - Run with hot reload"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  deps             - Download dependencies"
