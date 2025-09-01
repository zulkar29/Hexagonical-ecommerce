# Development Makefile for E-commerce SaaS

.PHONY: help dev-up dev-down dev-logs dev-clean dev-rebuild services-up services-down

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)E-commerce SaaS Development Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# Development Environment
dev-up: ## Start all development services
	@echo "$(GREEN)Starting development environment...$(NC)"
	docker-compose -f docker-compose.dev.yml up -d
	@echo "$(GREEN)Development environment started!$(NC)"
	@echo ""
	@echo "Services available at:"
	@echo "- Backend API: http://localhost:8080"
	@echo "- Storefront: http://localhost:3000" 
	@echo "- Dashboard: http://localhost:3001"
	@echo "- MinIO Console: http://localhost:9001"
	@echo "- MailHog: http://localhost:8025"
	@echo "- PostgreSQL: localhost:5432"
	@echo "- Redis: localhost:6379"

dev-down: ## Stop all development services
	@echo "$(YELLOW)Stopping development environment...$(NC)"
	docker-compose -f docker-compose.dev.yml down

dev-logs: ## Show logs from all services
	docker-compose -f docker-compose.dev.yml logs -f

dev-clean: ## Stop services and remove volumes
	@echo "$(RED)Cleaning up development environment...$(NC)"
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f

dev-rebuild: ## Rebuild and restart all services
	@echo "$(YELLOW)Rebuilding development environment...$(NC)"
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.dev.yml up -d

# Individual Services
services-up: ## Start only database and cache services
	@echo "$(GREEN)Starting core services (PostgreSQL, Redis, MinIO)...$(NC)"
	docker-compose -f docker-compose.dev.yml up -d postgres redis minio mailhog

services-down: ## Stop only database and cache services
	docker-compose -f docker-compose.dev.yml stop postgres redis minio mailhog

# Backend specific
backend-logs: ## Show backend logs
	docker-compose -f docker-compose.dev.yml logs -f backend

backend-shell: ## Access backend container shell
	docker-compose -f docker-compose.dev.yml exec backend sh

# Frontend specific
storefront-logs: ## Show storefront logs
	docker-compose -f docker-compose.dev.yml logs -f storefront

dashboard-logs: ## Show dashboard logs
	docker-compose -f docker-compose.dev.yml logs -f dashboard

# Database operations
db-migrate: ## Run database migrations
	docker-compose -f docker-compose.dev.yml exec backend go run cmd/migrate/main.go

db-shell: ## Access PostgreSQL shell
	docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d ecommerce_saas_dev

db-reset: ## Reset database (WARNING: Deletes all data)
	@echo "$(RED)This will delete all database data. Are you sure? [y/N]$(NC)" && read ans && [ $${ans:-N} = y ]
	docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS ecommerce_saas_dev;"
	docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -c "CREATE DATABASE ecommerce_saas_dev;"
	$(MAKE) db-migrate

# Development utilities
install: ## Install dependencies for all projects
	@echo "$(GREEN)Installing dependencies...$(NC)"
	cd backend && go mod tidy
	cd storefront && npm install
	cd dashboard && npm install

format: ## Format code in all projects
	@echo "$(GREEN)Formatting code...$(NC)"
	cd backend && go fmt ./...
	cd storefront && npm run prettier --write .
	cd dashboard && npm run prettier --write .

test: ## Run tests for all projects
	@echo "$(GREEN)Running tests...$(NC)"
	cd backend && go test ./...
	cd storefront && npm test
	cd dashboard && npm test

# Status check
status: ## Check status of all services
	docker-compose -f docker-compose.dev.yml ps