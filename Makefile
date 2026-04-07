# ============================================================
# Vaxen API — Makefile
# ============================================================

GO          = go
BIN_DIR     = ./bin
SERVER_BIN  = $(BIN_DIR)/server
IMAGE_NAME  = vaxen-api
IMAGE_TAG   = latest

COMPOSE_DEV  = docker compose -f compose.dev.yml
COMPOSE_PROD = docker compose -f compose.prod.yml

.PHONY: help install build run dev test test-coverage clean swagger migrate \
        seed-admin seed-rates \
        dev-up dev-down dev-logs dev-rebuild dev-seed-admin dev-seed-rates dev-migrate \
        prod-up prod-down prod-logs prod-migrate prod-seed-admin \
        docker-build lint fmt vet

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==================== Local Development ====================

install: ## Install dependencies and tools
	$(GO) mod download
	$(GO) mod tidy
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

build: ## Build the server binary
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(SERVER_BIN) main.go

run: build ## Build and run locally
	$(SERVER_BIN)

dev: ## Run with auto-reload (requires air)
	air

test: ## Run all tests
	CGO_ENABLED=1 $(GO) test ./... -count=1

test-coverage: ## Run tests with coverage report
	CGO_ENABLED=1 $(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR) coverage.out

swagger: ## Generate Swagger docs
	swag init -g main.go --output docs

migrate: ## Run database migrations (local)
	$(GO) run cmd/migrate/main.go

seed-admin: ## Create admin user (local). Usage: make seed-admin EMAIL=x PASSWORD=y
	@if [ -z "$(EMAIL)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Usage: make seed-admin EMAIL=admin@vaxen.io PASSWORD=yourpassword"; exit 1; fi
	$(GO) run cmd/seed/main.go admin $(EMAIL) $(PASSWORD)

seed-rates: ## Seed exchange rates (local). Usage: make seed-rates [BASE=USD TARGETS=EUR,GBP]
	@if [ -z "$(BASE)" ]; then $(GO) run cmd/seed/main.go exchange-rates --all; \
	else $(GO) run cmd/seed/main.go exchange-rates $(BASE) $(TARGETS); fi

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

# ==================== Docker Dev Environment ====================

dev-up: ## Start dev stack (API + Postgres + Redis + pgAdmin)
	$(COMPOSE_DEV) up -d --build
	@echo ""
	@echo "API:     http://localhost:8080"
	@echo "Swagger: http://localhost:8080/docs/index.html"
	@echo "pgAdmin: http://localhost:5050"
	@echo "DB:      localhost:5432 (postgres/password)"

dev-down: ## Stop dev stack
	$(COMPOSE_DEV) down

dev-rebuild: ## Rebuild and restart the API container
	$(COMPOSE_DEV) up -d --build api

dev-logs: ## Tail API logs
	$(COMPOSE_DEV) logs -f api

dev-migrate: ## Run migrations inside the dev stack
	$(COMPOSE_DEV) exec api /app/migrate

dev-seed-admin: ## Create admin in dev stack. Usage: make dev-seed-admin EMAIL=x PASSWORD=y
	@if [ -z "$(EMAIL)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Usage: make dev-seed-admin EMAIL=admin@vaxen.io PASSWORD=yourpassword"; exit 1; fi
	$(COMPOSE_DEV) exec api /app/seed admin $(EMAIL) $(PASSWORD)

dev-seed-rates: ## Seed exchange rates in dev stack
	$(COMPOSE_DEV) exec api /app/seed exchange-rates --all

# ==================== Docker Production ====================

prod-up: ## Start production stack
	$(COMPOSE_PROD) up -d --build

prod-down: ## Stop production stack
	$(COMPOSE_PROD) down

prod-logs: ## Tail production API logs
	$(COMPOSE_PROD) logs -f api

prod-migrate: ## Run migrations in production
	$(COMPOSE_PROD) exec api /app/migrate

prod-seed-admin: ## Create admin in production. Usage: make prod-seed-admin EMAIL=x PASSWORD=y
	@if [ -z "$(EMAIL)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Usage: make prod-seed-admin EMAIL=admin@vaxen.io PASSWORD=yourpassword"; exit 1; fi
	$(COMPOSE_PROD) exec api /app/seed admin $(EMAIL) $(PASSWORD)

# ==================== Docker Image ====================

docker-build: ## Build Docker image
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
