.PHONY: dev build run-backend run-frontend deps down clean logs ps help

# Default target
all: dev

dev: deps ## Start the full development environment (DB + properties + Frontend)
	@echo "Starting Backend and Frontend..."
	make -j2 run-backend run-frontend

run-backend: ## Run the Go backend
	go run cmd/council/main.go

run-frontend: #
# Frontend commands
frontend-dev: ## Run the React frontend in development mode
	cd frontend && npm run dev

frontend-build: ## Build the React frontend for production
	cd frontend && npm run build

# Database commands
migrate-up: ## Apply database migrations
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down: ## Revert database migrations
	migrate -path migrations -database "${DATABASE_URL}" down

build: ## Build both backend and frontend
	make frontend-build
	go build -o bin/council cmd/council/main.go

# --- Testing ---
test: ## Run unit tests (excluding infrastructure wrappers)
	@echo "Running tests..."
	go test -v -coverprofile=coverage.out ./...
	@# Filter out infrastructure and cmd/main.go from coverage report
	@grep -v -E "internal/infrastructure|cmd/" coverage.out > coverage.filtered.out
	@mv coverage.filtered.out coverage.out
	@go tool cover -func=coverage.out

# --- Dependencies (Docker) ---

deps: ## Start Postgres and other dependencies
	@echo "Starting dependencies..."
	docker compose up -d

down: ## Stop dependencies
	docker compose down

clean: ## Stop dependencies and delete data volumes
	docker compose down -v

logs: ## Follow docker logs
	docker compose logs -f

ps: ## Show running containers
	docker compose ps

help: ## Show this help message
	@echo "Usage: make [target]"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
