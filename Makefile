.PHONY: help build run test clean swagger migrate docker deps lint fmt vet

# Variables
APP_NAME := subscription-service
BUILD_DIR := ./bin
CONFIG_PATH := ./configs/config.yaml
MIGRATIONS_DIR := ./internal/infrastructure/database/postgres/migrations

# Help target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Dependencies
deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Code quality
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run

# Swagger documentation
swagger: ## Generate swagger documentation
	@echo "Generating swagger documentation..."
	swag init -g cmd/app/main.go -o ./api/swagger --parseDependency --parseInternal --parseVendor
	@echo "Swagger docs generated at ./api/swagger/"

swagger-fix: ## Fix swagger version issues and regenerate
	@echo "Fixing swagger issues..."
	rm -rf ./api/swagger/
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	swag init -g cmd/app/main.go -o ./api/swagger --parseDependency --parseInternal
	@echo "Swagger docs regenerated!"

swagger-serve: swagger ## Generate and serve swagger docs
	@echo "Swagger UI available at: http://localhost:8080/docs/"

# Database migrations
migrate-up: ## Run database migrations up
	@echo "Running database migrations up..."
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=up

migrate-down: ## Run database migrations down  
	@echo "Running database migrations down..."
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=down

migrate-create: ## Create new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

migrate-version: ## Show current migration version
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=version

migrate-force: ## Force migration to specific version (usage: make migrate-force version=0)
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=VERSION_NUMBER"; exit 1; fi
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=force -version=$(version)

migrate-reset: ## Reset migrations and start fresh
	@echo "Resetting migrations..."
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=force -version=1
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=down
	go run cmd/migrator/main.go -config=$(CONFIG_PATH) -migrations-dir="file://$(MIGRATIONS_DIR)" -action=up

# Build targets
build: deps fmt vet swagger ## Build the application
	@echo "Building $(APP_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/app/main.go
	go build -o $(BUILD_DIR)/migrator cmd/migrator/main.go

build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux cmd/app/main.go

# Run targets  
run: ## Run the application (without swagger for now)
	@echo "Starting $(APP_NAME)..."
	go run cmd/app/main.go -config=$(CONFIG_PATH)

run-with-swagger: swagger ## Run with swagger generation
	@echo "Starting $(APP_NAME) with swagger..."
	go run cmd/app/main.go -config=$(CONFIG_PATH)

run-dev: ## Run in development mode
	@echo "Starting $(APP_NAME) in development mode..."
	CONFIG_PATH=$(CONFIG_PATH) air

# Test targets
test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker targets
docker-build: ## Build docker image
	@echo "Building docker image..."
	docker build -t $(APP_NAME):latest .

docker-run: docker-build ## Run docker container
	@echo "Running docker container..."
	docker run --rm -p 8080:8080 --env-file .env $(APP_NAME):latest

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up --build

docker-compose-down: ## Stop services with docker-compose  
	@echo "Stopping services with docker-compose..."
	docker-compose down

docker-compose-logs: ## Show docker-compose logs
	docker-compose logs -f

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -rf ./api/swagger/docs.go ./api/swagger/swagger.json ./api/swagger/swagger.yaml
	rm -f coverage.out coverage.html
	docker-compose down --volumes --remove-orphans 2>/dev/null || true
	docker system prune -f

# Development setup
setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp configs/.env.example .env; echo ".env file created from example"; fi
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Edit .env file with your configuration"  
	@echo "2. Run: make migrate-up"
	@echo "3. Run: make run"

# CI/CD helpers
ci: deps lint test swagger build ## CI pipeline
	@echo "CI pipeline completed successfully!"

# Quick start
dev: setup migrate-up run ## Quick development start

dev-docker: ## Quick start with Docker
	@echo "Starting development environment with Docker..."
	docker-compose up --build -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 10
	make migrate-up
	make run

start-db: ## Start only PostgreSQL in Docker
	@echo "Starting PostgreSQL..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	sleep 10
	@echo "PostgreSQL is ready!"

stop-db: ## Stop PostgreSQL
	docker-compose down

full-stack: ## Start full application stack
	docker-compose up --build

# Default target
all: clean ci docker-build ## Build everything