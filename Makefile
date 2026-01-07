.PHONY: build run test migrate-up migrate-down migrate-status lint fmt

# Build the application
build:
	go build -o bin/movie-booking cmd/main.go

# Run the API server
run:
	go run cmd/main.go --api

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run migrations up
migrate-up:
	go run cmd/main.go --migrate --migration-command=up

# Run migrations down
migrate-down:
	go run cmd/main.go --migrate --migration-command=down

# Check migration status
migrate-status:
	go run cmd/main.go --migrate --migration-command=status

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy
