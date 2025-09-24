.PHONY: help build run test clean docker-build docker-run migrate seed setup worker lint fmt

help: ## Show help message
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  worker       - Run the background worker"
	@echo "  test         - Run all tests"
	@echo "  test-unit    - Run unit tests only"
	@echo "  test-user    - Run user module tests"
	@echo "  test-product - Run product module tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  test-bench   - Run benchmark tests"
	@echo "  test-race    - Run tests with race detection"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  migrate      - Run database migrations"
	@echo "  seed         - Populate database with demo data"
	@echo "  setup        - Run migrations and seed data"
	@echo "  dev          - Run in development mode with hot reload"

build: ## Build the application
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/worker ./cmd/worker

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-unit: ## Run unit tests only
	go test -v ./internal/...

test-user: ## Run user module tests
	go test -v ./internal/user

test-product: ## Run product module tests
	go test -v ./internal/product

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-bench: ## Run benchmark tests
	go test -bench=. -benchmem ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-build: ## Build Docker image
	docker build -t ecommerce-saas-backend:latest -f Dockerfile.dev .

docker-run: ## Run Docker container
	docker run -p 8080:8080 ecommerce-saas-backend:latest

migrate: ## Run database migrations
	go run ./cmd/migrate/main.go -action=up

seed: ## Populate database with demo data
	go run ./cmd/migrate/main.go -action=seed

setup: ## Run migrations and seed data
	go run ./cmd/migrate/main.go -action=up
	go run ./cmd/migrate/main.go -action=seed

dev: ## Run in development mode with hot reload
	air -c .air.toml

worker: ## Run the background worker
	go run ./cmd/worker/main.go

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .