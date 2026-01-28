.PHONY: run build test swagger docker-up docker-down migrate clean

# Variables
APP_NAME=go-echo-starter
MAIN_PATH=./cmd/api

# Run the application
run:
	@echo "Starting application..."
	@go run $(MAIN_PATH)/main.go

# Build binary
build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) $(MAIN_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate swagger docs
swagger:
	@echo "Generating swagger docs..."
	@swag init -g $(MAIN_PATH)/main.go -o ./docs

# Docker commands
docker-up:
	@echo "Starting PostgreSQL..."
	@docker-compose up -d

docker-down:
	@echo "Stopping PostgreSQL..."
	@docker-compose down

# Run migrations
migrate:
	@echo "Running migrations..."
	@docker exec -i go-echo-postgres psql -U postgres -d go_echo_db < migrations/001_create_users_table.sql

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf tmp/

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Run checks (tidy, lint, test)
check: deps test
	@echo "All checks passed!"

# Lint (placeholder for now, can be replaced with golangci-lint)
lint:
	@echo "Running lint..."
	@go vet ./...

# Install swag CLI
install-swag:
	@go install github.com/swaggo/swag/cmd/swag@latest
