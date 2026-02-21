.PHONY: help build run test clean migrate swagger install dev

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

build: ## Build the application
	go build -o bin/api main.go

run: ## Run the application
	go run main.go

dev: ## Run with auto-reload (requires air)
	air

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

swagger: ## Generate Swagger documentation
	swag init -g main.go --output docs

migrate: ## Run database migrations
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

docker-build: ## Build Docker image
	docker build -t vaxen-api .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env vaxen-api

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...
