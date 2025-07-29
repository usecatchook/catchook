# Variables
APP_NAME=webhook-api
DOCKER_COMPOSE_DEV=docker-compose -f docker-compose.dev.yml
DOCKER_COMPOSE=docker-compose
GO_VERSION=1.24
MAIN_PATH=./cmd/api

# Colors for output
GREEN=\033[32m
YELLOW=\033[33m
RED=\033[31m
NC=\033[0m # No Color

.PHONY: help setup dev build run test lint clean docker-up docker-down deps

# Default target
help: ## Show this help message
	@echo "$(GREEN)Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development Setup
setup: ## Setup development environment
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@cp .env.example .env
	@echo "$(YELLOW)✓ Created .env file from template$(NC)"
	@go mod download
	@echo "$(YELLOW)✓ Downloaded Go dependencies$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d postgres redis
	@echo "$(YELLOW)✓ Started PostgreSQL and Redis containers$(NC)"
	@echo "$(GREEN)Setup complete! Edit .env file with your configuration.$(NC)"
	@echo "$(YELLOW)💡 Use 'make dev-services-gui' to start with GUI tools$(NC)"

dev-services: ## Start development services (PostgreSQL + Redis)
	@echo "$(GREEN)Starting development services...$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d postgres redis
	@echo "$(YELLOW)✓ PostgreSQL and Redis started$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432 (postgres/postgres)$(NC)"
	@echo "$(YELLOW)Redis: localhost:6379$(NC)"

dev-services-gui: ## Start development services with GUI tools
	@echo "$(GREEN)Starting development services with GUI...$(NC)"
	@$(DOCKER_COMPOSE_DEV) --profile gui up -d
	@echo "$(YELLOW)✓ All services started with GUI tools$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432 (postgres/postgres)$(NC)"
	@echo "$(YELLOW)Redis: localhost:6379$(NC)"
	@echo "$(YELLOW)pgAdmin: http://localhost:8082 (admin@webhook-api.com/admin)$(NC)"
	@echo "$(YELLOW)Redis Commander: http://localhost:8081 (admin/admin)$(NC)"

dev-services-stop: ## Stop development services
	@echo "$(GREEN)Stopping development services...$(NC)"
	@$(DOCKER_COMPOSE_DEV) down
	@echo "$(YELLOW)✓ Development services stopped$(NC)"

dev-services-logs: ## Show logs from development services
	@$(DOCKER_COMPOSE_DEV) logs -f

dev-services-clean: ## Clean development services (remove volumes)
	@echo "$(GREEN)Cleaning development services...$(NC)"
	@$(DOCKER_COMPOSE_DEV) down --volumes --remove-orphans
	@echo "$(YELLOW)✓ Development services cleaned$(NC)"

# Development
dev: ## Run application in development mode with hot reload
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	@air

run: ## Run application
	@echo "$(GREEN)Running application...$(NC)"
	@go run $(MAIN_PATH)

build: ## Build application binary
	@echo "$(GREEN)Building application...$(NC)"
	@go build -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN_PATH)
	@echo "$(YELLOW)✓ Binary created: bin/$(APP_NAME)$(NC)"

# Dependencies
deps: ## Download and tidy dependencies
	@echo "$(GREEN)Managing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(YELLOW)✓ Dependencies updated$(NC)"

deps-upgrade: ## Upgrade all dependencies
	@echo "$(GREEN)Upgrading dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(YELLOW)✓ Dependencies upgraded$(NC)"

# Testing
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(YELLOW)✓ Coverage report: coverage.html$(NC)"

benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# Code Quality
lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	@golangci-lint run

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...

# Database
db-up: ## Start database container
	@echo "$(GREEN)Starting database...$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d postgres
	@echo "$(YELLOW)✓ PostgreSQL started$(NC)"

db-down: ## Stop database container
	@echo "$(GREEN)Stopping database...$(NC)"
	@$(DOCKER_COMPOSE_DEV) stop postgres
	@echo "$(YELLOW)✓ PostgreSQL stopped$(NC)"

db-reset: ## Reset database (stop, remove, start)
	@echo "$(GREEN)Resetting database...$(NC)"
	@$(DOCKER_COMPOSE_DEV) stop postgres
	@$(DOCKER_COMPOSE_DEV) rm -f postgres
	@$(DOCKER_COMPOSE_DEV) up -d postgres
	@echo "$(YELLOW)✓ Database reset$(NC)"

db-connect: ## Connect to PostgreSQL
	@echo "$(GREEN)Connecting to PostgreSQL...$(NC)"
	@$(DOCKER_COMPOSE_DEV) exec postgres psql -U postgres -d webhook_api

# Redis
redis-up: ## Start Redis container
	@echo "$(GREEN)Starting Redis...$(NC)"
	@$(DOCKER_COMPOSE_DEV) up -d redis
	@echo "$(YELLOW)✓ Redis started$(NC)"

redis-down: ## Stop Redis container
	@echo "$(GREEN)Stopping Redis...$(NC)"
	@$(DOCKER_COMPOSE_DEV) stop redis
	@echo "$(YELLOW)✓ Redis stopped$(NC)"

redis-cli: ## Connect to Redis CLI
	@echo "$(GREEN)Connecting to Redis CLI...$(NC)"
	@$(DOCKER_COMPOSE_DEV) exec redis redis-cli

# Docker
docker-up: ## Start all services with Docker Compose
	@echo "$(GREEN)Starting all services...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(YELLOW)✓ All services started$(NC)"

docker-down: ## Stop all services
	@echo "$(GREEN)Stopping all services...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(YELLOW)✓ All services stopped$(NC)"

docker-logs: ## Show logs from all services
	@$(DOCKER_COMPOSE) logs -f

docker-build: ## Build application Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(YELLOW)✓ Docker image built: $(APP_NAME):latest$(NC)"

# Cleanup
clean: ## Clean build artifacts and containers
	@echo "$(GREEN)Cleaning up...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@$(DOCKER_COMPOSE) down --volumes --remove-orphans
	@echo "$(YELLOW)✓ Cleanup complete$(NC)"

# Production
build-prod: ## Build production binary with optimizations
	@echo "$(GREEN)Building production binary...$(NC)"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN_PATH)
	@echo "$(YELLOW)✓ Production binary created: bin/$(APP_NAME)$(NC)"

# Installation helpers
install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(YELLOW)✓ Development tools installed$(NC)"

# Health checks
health: ## Check application health
	@echo "$(GREEN)Checking application health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Application not responding$(NC)"

# SQLC
sqlc-generate: ## Generate SQLC code
	@echo "$(GREEN)Generating SQLC code...$(NC)"
	@sqlc generate
	@echo "$(YELLOW)✓ SQLC code generated$(NC)"

sqlc-verify: ## Verify SQLC queries
	@echo "$(GREEN)Verifying SQLC queries...$(NC)"
	@sqlc verify
	@echo "$(YELLOW)✓ SQLC queries verified$(NC)"

# Migration helpers (will be implemented later)
migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	@go run cmd/migrate/main.go up

migrate-down: ## Run database migrations down
	@echo "$(GREEN)Running migrations down...$(NC)"
	@go run cmd/migrate/main.go down

migrate-create: ## Create new migration file (usage: make migrate-create NAME=migration_name)
	@echo "$(GREEN)Creating migration: $(NAME)...$(NC)"
	@go run cmd/migrate/main.go create $(NAME)

# Database with schema
db-migrate: dev-services ## Start DB and apply schema
	@echo "$(GREEN)Applying database schema...$(NC)"
	@sleep 2  # Wait for DB to be ready
	@docker exec webhook-api-postgres-dev psql -U postgres -d webhook_api -f /docker-entrypoint-initdb.d/001_users.sql
	@echo "$(YELLOW)✓ Database schema applied$(NC)"

db-schema: ## Apply schema to running database
	@echo "$(GREEN)Applying database schema...$(NC)"
	@docker exec webhook-api-postgres-dev psql -U postgres -d webhook_api -f /docker-entrypoint-initdb.d/001_users.sql
	@echo "$(YELLOW)✓ Database schema applied$(NC)"

sqlc-dev: db-migrate sqlc-generate ## Setup DB, apply schema and generate SQLC
	@echo "$(GREEN)✓ Database ready with SQLC generated code$(NC)"