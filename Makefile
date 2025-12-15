.PHONY: dev build run-backend run-frontend deps down clean logs ps help

# Default target
all: dev

dev: deps ## Start the full development environment (DB + properties + Frontend)
	@echo "Starting Backend and Frontend..."
	make -j2 run-backend run-frontend

run-backend: ## Run the Go backend
	go run cmd/server/main.go

run-frontend: ## Run the React frontend
	cd frontend && npm run dev

build: ## Build both backend and frontend
	cd frontend && npm run build
	go build -o bin/server cmd/server/main.go

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
