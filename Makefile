# GO build commands
GO_CMD                    = go
BUILD_CMD                 = $(GO_CMD) build
LINUX_GO_BUILD_CMD        = CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD_CMD)

# Paths and filenames
BIN_DIR                   = ./bin
LOCAL_BINARY_PATH         = $(BIN_DIR)/local_server
DOCKER_SERVER_BIN_PATH    = $(BIN_DIR)/dev_server

# Entry points
SERVER_ENTRYPOINT         = main.go
MIGRATE_ENTRYPOINT        = cmd/migrate/main.go

# Containerization
CONTAINER_CMD 			  = docker
COMPOSE_CMD 			  = $(CONTAINER_CMD) compose
CONTAINER_FILE            = Dockerfile
COMPOSE_FILE              = docker-compose.dev.yml
CONTAINER_IMAGE_NAME      = vaxen-api
CONTAINER_IMAGE_TAG       = latest

.PHONY: help build run test clean migrate swagger install dev

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

build: ## Build the application
	$(BUILD_CMD) -o $(LOCAL_BINARY_PATH) $(SERVER_ENTRYPOINT)

run: build ## Build and run the application
	$(LOCAL_BINARY_PATH)

dev: ## Run with auto-reload (requires air)
	air

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf $(BIN_DIR)
	rm -f coverage.out

swagger: ## Generate Swagger documentation
	swag init -g $(SERVER_ENTRYPOINT) --output docs

migrate: ## Run database migrations
	@echo "Running database migrations..."
	go run $(MIGRATE_ENTRYPOINT)

seed-admin: ## Create a platform admin user (usage: make seed-admin EMAIL=admin@vaxen.io PASSWORD=securepass)
	@if [ -z "$(EMAIL)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Usage: make seed-admin EMAIL=admin@vaxen.io PASSWORD=yourpassword"; \
		exit 1; \
	fi
	go run cmd/seed/main.go admin $(EMAIL) $(PASSWORD)

dockerbinary: $(DOCKER_SERVER_BIN_PATH) ## Build the Go binary for Docker

$(DOCKER_SERVER_BIN_PATH): $(wildcard *.go) go.mod go.sum
	@echo "🔨 Building Go binary for Alpine Linux..."
	@mkdir -p $(BIN_DIR)
	$(LINUX_GO_BUILD_CMD) -o $(DOCKER_SERVER_BIN_PATH) $(SERVER_ENTRYPOINT)

docker-build: ## Build Docker image
	$(CONTAINER_CMD) build -t $(CONTAINER_IMAGE_NAME):$(CONTAINER_IMAGE_TAG) -f $(CONTAINER_FILE) .

docker-run: ## Run Docker container
	$(CONTAINER_CMD) run -p 8080:8080 --env-file .env $(CONTAINER_IMAGE_NAME):$(CONTAINER_IMAGE_TAG)

docker-compose-up: dockerbinary ## Start services with docker-compose
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) up -d

docker-compose-down: ## Stop services with docker-compose
	$(COMPOSE_CMD) -f $(COMPOSE_FILE) down

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...
