# Go Microservice Template Makefile
# Usage: make init MODULE=github.com/your-org/your-service SERVICE=your-service

MODULE ?= github.com/your-org/your-service
SERVICE ?= your-service
CURRENT_MODULE = github.com/minisource/template_go

.PHONY: help init build run test lint clean docker-build docker-run

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## Initialize project with new module name (make init MODULE=github.com/your-org/service SERVICE=service-name)
	@echo "Initializing project..."
	@echo "  Module: $(MODULE)"
	@echo "  Service: $(SERVICE)"
	@if [ "$(MODULE)" = "github.com/your-org/your-service" ]; then \
		echo "Error: Please provide MODULE parameter"; \
		echo "Usage: make init MODULE=github.com/your-org/your-service SERVICE=your-service"; \
		exit 1; \
	fi
	@echo "Replacing module name in go.mod..."
	@cd src && sed -i 's|$(CURRENT_MODULE)|$(MODULE)|g' go.mod
	@echo "Replacing imports in all Go files..."
	@find src -name "*.go" -type f -exec sed -i 's|$(CURRENT_MODULE)|$(MODULE)|g' {} \;
	@echo "Updating config files..."
	@find src/config -name "*.yml" -type f -exec sed -i 's|DiviPay|$(SERVICE)|g' {} \;
	@echo "Updating Dockerfile..."
	@sed -i 's|backend|$(SERVICE)|g' src/Dockerfile
	@echo "Updating docker-compose..."
	@sed -i 's|backend|$(SERVICE)|g' docker/docker-compose.yml
	@echo "Running go mod tidy..."
	@cd src && go mod tidy
	@echo ""
	@echo "✅ Project initialized successfully!"
	@echo "   Module: $(MODULE)"
	@echo "   Service: $(SERVICE)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Update src/config/config-development.yml with your database settings"
	@echo "  2. Update src/docs/docs.go with your API info"
	@echo "  3. Run 'make run' to start the server"

build: ## Build the application
	@cd src && go build -o bin/server ./cmd/

run: ## Run the application
	@cd src && go run ./cmd/main.go

test: ## Run all tests
	@cd src && go test -v -cover ./...

test-unit: ## Run unit tests only
	@cd src && go test -v -cover ./tests/unit/...

test-integration: ## Run integration tests only
	@cd src && go test -v -cover ./tests/integration/...

lint: ## Run linter
	@cd src && golangci-lint run ./...

fmt: ## Format code
	@cd src && gofmt -w .

tidy: ## Tidy and download dependencies
	@cd src && go mod tidy && go mod download

clean: ## Clean build artifacts
	@rm -rf src/bin
	@rm -rf src/logs/*.log

swagger: ## Generate swagger docs
	@cd src && swag init -g cmd/main.go -o docs

docker-build: ## Build Docker image
	@docker build -t $(SERVICE):latest -f src/Dockerfile src/

docker-run: ## Run with Docker Compose
	@cd docker && docker-compose up -d

docker-stop: ## Stop Docker Compose
	@cd docker && docker-compose down

migrate: ## Run database migrations
	@cd src && go run ./cmd/main.go migrate

dev: ## Run in development mode with hot reload (requires air)
	@cd src && air

install-tools: ## Install development tools
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Development tools installed"
