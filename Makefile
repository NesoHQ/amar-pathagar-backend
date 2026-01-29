############################
# Makefile for Amar Pathagar Backend
############################

# --------------------------------------------------
# Load environment variables from .env
# --------------------------------------------------
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# --------------------------------------------------
# Configuration
# --------------------------------------------------
COMPOSE_FILE     = docker-compose.yml
COMPOSE_DEV_FILE = docker-compose.dev.yml
BINARY_NAME      = amar-pathagar-api
MAIN_PATH        = ./cmd/api

.DEFAULT_GOAL := help

# --------------------------------------------------
# Help
# --------------------------------------------------
.PHONY: help
help: ## Show this help message
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║         Amar Pathagar Backend - Makefile Commands         ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --------------------------------------------------
# Development
# --------------------------------------------------
.PHONY: dev
dev: ## Start development environment (with hot reload)
	docker compose -f $(COMPOSE_DEV_FILE) up -d
	@echo "✅ Development environment started"
	@echo "📝 API: http://localhost:8080"
	@echo "🔍 Health: http://localhost:8080/health"
	@echo "📋 Logs: make logs"

.PHONY: logs
logs: ## Follow application logs
	docker compose -f $(COMPOSE_DEV_FILE) logs -f backend

.PHONY: restart
restart: ## Restart development environment
	docker compose -f $(COMPOSE_DEV_FILE) restart backend
	@echo "✅ Backend restarted"

.PHONY: stop
stop: ## Stop development environment
	docker compose -f $(COMPOSE_DEV_FILE) stop

# --------------------------------------------------
# Production
# --------------------------------------------------
.PHONY: up
up: ## Start production environment
	docker compose -f $(COMPOSE_FILE) up -d --build
	@echo "✅ Production environment started"

.PHONY: down
down: ## Stop and remove all containers
	docker compose -f $(COMPOSE_FILE) down
	docker compose -f $(COMPOSE_DEV_FILE) down
	@echo "✅ All containers stopped and removed"

.PHONY: build
build: ## Build Docker image
	docker compose -f $(COMPOSE_FILE) build --no-cache

# --------------------------------------------------
# Database
# --------------------------------------------------
.PHONY: db-shell
db-shell: ## Open PostgreSQL shell
	docker compose -f $(COMPOSE_DEV_FILE) exec postgres psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: db-reset
db-reset: ## Reset database (drop and recreate)
	@echo "⚠️  This will delete all data. Press Ctrl+C to cancel..."
	@sleep 3
	docker compose -f $(COMPOSE_DEV_FILE) exec postgres psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	docker compose -f $(COMPOSE_DEV_FILE) exec postgres psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	docker compose -f $(COMPOSE_DEV_FILE) exec postgres psql -U $(DB_USER) -d $(DB_NAME) -f /docker-entrypoint-initdb.d/init.sql
	@echo "✅ Database reset complete"

.PHONY: db-backup
db-backup: ## Backup database
	@mkdir -p backups
	docker compose -f $(COMPOSE_DEV_FILE) exec -T postgres pg_dump -U $(DB_USER) $(DB_NAME) > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✅ Database backed up to backups/"

.PHONY: db-restore
db-restore: ## Restore database from backup (usage: make db-restore FILE=backups/backup.sql)
	@if [ -z "$(FILE)" ]; then echo "❌ Usage: make db-restore FILE=backups/backup.sql"; exit 1; fi
	docker compose -f $(COMPOSE_DEV_FILE) exec -T postgres psql -U $(DB_USER) $(DB_NAME) < $(FILE)
	@echo "✅ Database restored from $(FILE)"

# --------------------------------------------------
# Local Development (without Docker)
# --------------------------------------------------
.PHONY: run
run: ## Run locally (without Docker)
	@echo "🚀 Starting server..."
	go run $(MAIN_PATH)/main.go

.PHONY: run-watch
run-watch: ## Run locally with hot reload (air)
	@echo "🚀 Starting server with hot reload..."
	air -c .air.toml

.PHONY: install
install: ## Install dependencies
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

.PHONY: build-binary
build-binary: ## Build standalone binary
	@echo "🔨 Building binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) $(MAIN_PATH)/main.go
	@echo "✅ Binary built: $(BINARY_NAME)"

# --------------------------------------------------
# Testing
# --------------------------------------------------
.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detector
	go test -race -v ./...

# --------------------------------------------------
# Code Quality
# --------------------------------------------------
.PHONY: lint
lint: ## Run linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Install: https://golangci-lint.run/usage/install/"; \
	fi

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	@echo "✅ Code formatted"

.PHONY: vet
vet: ## Run go vet
	go vet ./...
	@echo "✅ Code vetted"

.PHONY: tidy
tidy: ## Tidy go.mod
	go mod tidy
	@echo "✅ go.mod tidied"

# --------------------------------------------------
# Docker Utilities
# --------------------------------------------------
.PHONY: ps
ps: ## Show running containers
	docker compose -f $(COMPOSE_DEV_FILE) ps

.PHONY: shell
shell: ## Open shell in backend container
	docker compose -f $(COMPOSE_DEV_FILE) exec backend sh

.PHONY: clean
clean: ## Clean up containers, volumes, and build artifacts
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_DEV_FILE) down -v
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "✅ Cleanup complete"

# --------------------------------------------------
# Monitoring
# --------------------------------------------------
.PHONY: stats
stats: ## Show container stats
	docker stats --no-stream

.PHONY: health
health: ## Check API health
	@curl -s http://localhost:8080/health | jq . || echo "❌ API not responding"

# --------------------------------------------------
# Database Seeding
# --------------------------------------------------
.PHONY: seed
seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@# Add your seed script here
	@echo "✅ Database seeded"

# --------------------------------------------------
# Quick Commands
# --------------------------------------------------
.PHONY: quick-start
quick-start: dev logs ## Quick start (dev + logs)

.PHONY: quick-restart
quick-restart: restart logs ## Quick restart (restart + logs)

.PHONY: quick-clean
quick-clean: down clean ## Quick clean (down + clean)

# --------------------------------------------------
# Info
# --------------------------------------------------
.PHONY: info
info: ## Show project information
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║              Amar Pathagar Backend Info                   ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📦 Project: Amar Pathagar Backend API"
	@echo "🔧 Language: Go $(shell go version | awk '{print $$3}')"
	@echo "🐳 Docker: $(shell docker --version | awk '{print $$3}' | tr -d ',')"
	@echo "📂 Main: $(MAIN_PATH)"
	@echo "🔨 Binary: $(BINARY_NAME)"
	@echo ""
	@echo "🌐 Endpoints:"
	@echo "   - API: http://localhost:8080"
	@echo "   - Health: http://localhost:8080/health"
	@echo "   - Database: localhost:5432"
	@echo ""
	@echo "📚 Documentation: README.md"
	@echo ""
