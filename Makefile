# Catchook Makefile

.PHONY: help dev-api dev-app build-api build-app test lint clean

# Default target
help: ## Show this help message
	@echo "Catchook Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development
dev-api: ## Start API in development mode
	cd cmd/api && go run main.go

dev-app: ## Start frontend in development mode
	cd app && npm run dev

dev: ## Start both API and frontend in development mode
	make -j2 dev-api dev-app

# Build
build-api: ## Build API binary
	go build -o bin/api cmd/api/main.go

build-app: ## Build frontend for production
	cd app && npm run build

build: build-api build-app ## Build both API and frontend

# Testing
test: ## Run all tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Linting
lint: ## Run linters
	golangci-lint run
	cd app && npm run lint

# Database
include .env
export

migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir storage/postgres/schema -seq $$name

migrate-up: ## Run all pending migrations
	migrate -path internal/platform/storage/postgres/schema -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

migrate-down: ## Rollback the last migration
	migrate -path internal/platform/storage/postgres/schema -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" down 1

migrate-force: ## Force set migration version
	@read -p "Enter version: " version; \
	migrate -path internal/platform/storage/postgres/schema -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" force $$version

migrate-status: ## Show migration status
	migrate -path internal/platform/storage/postgres/schema -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" version

# Dependencies
deps: ## Install all dependencies
	go mod tidy
	cd app && npm install

# Docker
docker-dev: ## Start development environment with Docker
	docker-compose -f docker-compose.dev.yml up -d

docker-stop: ## Stop Docker development environment
	docker-compose -f docker-compose.dev.yml down

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf app/.next/
	rm -rf app/out/
	go clean -cache

# Setup
setup: deps docker-dev ## Complete development setup
	@echo "✅ Setup complete! Run 'make dev' to start development."