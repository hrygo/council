# ============================================================================
# The Council - Development Makefile
# ============================================================================
# Usage: make [target]
# Run `make help` to see all available commands
# ============================================================================

.PHONY: all help \
        start stop restart status \
        start-all stop-all \
        start-db stop-db start-backend stop-backend start-frontend stop-frontend \
        build test test-backend test-frontend lint fmt check clean install \
        coverage coverage-backend coverage-frontend

# ============================================================================
# 🎨 Colors
# ============================================================================
BOLD   := \033[1m
CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
RESET  := \033[0m

# ============================================================================
# 📦 Variables (loaded from .env if exists)
# ============================================================================
-include .env
export

GO_BIN         := bin/council
DATABASE_URL   ?= postgres://council:council_password@localhost:5432/council_db?sslmode=disable
LLM_PROVIDER   ?= gemini
LLM_MODEL      ?= gemini-2.0-flash

# ============================================================================
# 🚀 Default
# ============================================================================
all: help

# ============================================================================
# 🔄 LIFECYCLE COMMANDS (Primary)
# ============================================================================

start: start-all ## 🚀 Start everything (DB + Backend + Frontend)

stop: stop-all ## 🛑 Stop everything

restart: stop start ## 🔄 Restart everything

status: ## 📊 Show status of all services
	@echo "$(BOLD)$(CYAN)📊 Service Status$(RESET)"
	@echo "════════════════════════════════════════"
	@echo ""
	@echo "$(BOLD)🐳 Docker Services:$(RESET)"
	@docker compose ps 2>/dev/null || echo "   Not running"
	@echo ""
	@echo "$(BOLD)🔧 Backend (port 8080):$(RESET)"
	@lsof -ti:8080 >/dev/null 2>&1 && echo "   $(GREEN)● Running$(RESET) (PID: $$(lsof -ti:8080))" || echo "   $(RED)○ Stopped$(RESET)"
	@echo ""
	@echo "$(BOLD)🎨 Frontend (port 5173/5174):$(RESET)"
	@lsof -ti:5173 >/dev/null 2>&1 && echo "   $(GREEN)● Running$(RESET) on :5173" || \
		(lsof -ti:5174 >/dev/null 2>&1 && echo "   $(GREEN)● Running$(RESET) on :5174" || echo "   $(RED)○ Stopped$(RESET)")
	@echo ""

# ============================================================================
# 🐳 DOCKER SERVICES
# ============================================================================

start-db: ## 🐳 Start database services (Postgres + Redis)
	@echo "$(CYAN)🐳 Starting Docker services...$(RESET)"
	@docker compose up -d
	@echo "$(GREEN)✅ Docker services started$(RESET)"
	@docker compose ps

stop-db: ## 🛑 Stop database services
	@echo "$(YELLOW)🛑 Stopping Docker services...$(RESET)"
	@docker compose down
	@echo "$(GREEN)✅ Docker services stopped$(RESET)"

restart-db: stop-db start-db ## 🔄 Restart database services

logs-db: ## 📜 Follow database logs
	@docker compose logs -f

reset-db: ## ⚠️ Reset database (DELETE ALL DATA)
	@echo "$(RED)$(BOLD)⚠️ WARNING: This will DELETE all data!$(RESET)"
	@read -p "Are you sure? [y/N]: " confirm && [ "$$confirm" = "y" ] || exit 1
	@docker compose down -v
	@docker compose up -d
	@sleep 3
	@echo "$(GREEN)✅ Database reset complete$(RESET)"

# ============================================================================
# 🔧 BACKEND
# ============================================================================

start-backend: ## 🔧 Start Go backend
	@echo "$(CYAN)🔧 Starting Backend on :8080...$(RESET)"
	@lsof -ti:8080 >/dev/null 2>&1 && { echo "$(YELLOW)⚠️ Port 8080 already in use. Stopping...$(RESET)"; make stop-backend; sleep 1; } || true
	@env DATABASE_URL="$(DATABASE_URL)" \
		LLM_PROVIDER="$(LLM_PROVIDER)" \
		LLM_MODEL="$(LLM_MODEL)" \
		GEMINI_API_KEY="$(GEMINI_API_KEY)" \
		go run cmd/council/main.go &
	@sleep 3
	@lsof -ti:8080 >/dev/null 2>&1 && echo "$(GREEN)✅ Backend started$(RESET)" || echo "$(RED)❌ Backend failed to start. Check: make logs-backend$(RESET)"

stop-backend: ## 🛑 Stop Go backend
	@echo "$(YELLOW)🛑 Stopping Backend...$(RESET)"
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@echo "$(GREEN)✅ Backend stopped$(RESET)"

restart-backend: stop-backend start-backend ## 🔄 Restart backend

logs-backend: ## 📜 Tail backend logs (if using file logging)
	@echo "$(YELLOW)Backend logs are in terminal output$(RESET)"

# ============================================================================
# 🎨 FRONTEND
# ============================================================================

start-frontend: ## 🎨 Start React frontend
	@echo "$(CYAN)🎨 Starting Frontend...$(RESET)"
	@cd frontend && npm run dev &
	@sleep 2
	@echo "$(GREEN)✅ Frontend started$(RESET)"

stop-frontend: ## 🛑 Stop React frontend
	@echo "$(YELLOW)🛑 Stopping Frontend...$(RESET)"
	@lsof -ti:5173 | xargs kill -9 2>/dev/null || true
	@lsof -ti:5174 | xargs kill -9 2>/dev/null || true
	@echo "$(GREEN)✅ Frontend stopped$(RESET)"

restart-frontend: stop-frontend start-frontend ## 🔄 Restart frontend

# ============================================================================
# 🚀 COMBINED LIFECYCLE
# ============================================================================

start-all: ## 🚀 Start all services
	@echo "$(GREEN)$(BOLD)🚀 Starting The Council...$(RESET)"
	@echo ""
	@make start-db
	@echo ""
	@make start-backend
	@echo ""
	@make start-frontend
	@echo ""
	@echo "$(GREEN)$(BOLD)════════════════════════════════════════$(RESET)"
	@echo "$(GREEN)$(BOLD)✅ All services started!$(RESET)"
	@echo "   $(CYAN)Backend:  http://localhost:8080$(RESET)"
	@echo "   $(CYAN)Frontend: http://localhost:5173$(RESET)"
	@echo "$(GREEN)$(BOLD)════════════════════════════════════════$(RESET)"

stop-all: ## 🛑 Stop all services
	@echo "$(YELLOW)$(BOLD)🛑 Stopping The Council...$(RESET)"
	@make stop-frontend
	@make stop-backend
	@make stop-db
	@echo "$(GREEN)✅ All services stopped$(RESET)"

# ============================================================================
# 🏗️ BUILD & TEST
# ============================================================================

build: lint ## 🏗️ Build production binaries
	@echo "$(GREEN)$(BOLD)🏗️ Building...$(RESET)"
	@cd frontend && npm run build
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o $(GO_BIN) cmd/council/main.go
	@echo "$(GREEN)✅ Build complete: $(GO_BIN)$(RESET)"

test: test-backend test-frontend ## 🧪 Run all tests (Backend + Frontend)

test-backend: ## 🔧 Run Go backend tests
	@echo "$(CYAN)🧪 Running backend tests...$(RESET)"
	@go test -v -race -coverprofile=coverage.out ./...

test-frontend: ## 🎨 Run React frontend tests
	@echo "$(CYAN)📅 Running frontend tests...$(RESET)"
	@cd frontend && npm run test -- --run

test-short: ## ⚡ Quick backend tests (no race detector)
	@go test -short ./...

coverage: coverage-backend coverage-frontend ## 📊 Run all coverage (Dashboard)
	@echo ""
	@echo "$(BOLD)$(CYAN)📈 FINAL COVERAGE DASHBOARD$(RESET)"
	@echo "════════════════════════════════════════════════════════════════"
	@printf "  $(BOLD)%-20s$(RESET) | $(BOLD)%s$(RESET)\n" "Domain" "Coverage Score"
	@echo "────────────────────────────────────────────────────────────────"
	@BE_RAW=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	FE_RAW=$$(grep -E "Lines" frontend/coverage_summary.txt | grep -oE "[0-9.]+" | head -1 || echo "0"); \
	printf "  %-20s | $(GREEN)%s%%$(RESET)\n" "Backend (Go)" "$$BE_RAW"; \
	printf "  %-20s | $(GREEN)%s%%$(RESET)\n" "Frontend (React)" "$$FE_RAW"; \
	echo "────────────────────────────────────────────────────────────────"; \
	AVG=$$(echo "scale=2; ($$BE_RAW + $$FE_RAW) / 2" | bc 2>/dev/null || echo "N/A"); \
	printf "  $(BOLD)%-20s$(RESET) | $(BOLD)%s%%$(RESET)\n" "Overall Average" "$$AVG"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "$(CYAN)Detailed reports:$(RESET)"
	@echo "  Backend  -> $(BOLD)coverage.html$(RESET)"
	@echo "  Frontend -> $(BOLD)frontend/coverage/index.html$(RESET)"
	@echo ""

coverage-backend: test-backend ## 🔧 Run backend coverage summary (Package List)
	@echo "$(CYAN)📊 Backend Coverage by Package:$(RESET)"
	@echo "-------------------------------------------|---------"
	@printf "  %-40s | %s\n" "Package" "Coverage"
	@echo "-------------------------------------------|---------"
	@go test -cover ./... | sed 's/github.com\/hrygo\/council\///g' | \
		awk '/^ok/ { printf "  %-40s | %s\n", $$2, $$5 } \
		     /^\?/ { printf "  %-40s | %s\n", $$2, "0.0%*" } \
		     /^[[:space:]]+internal/ { printf "  %-40s | %s\n", $$1, $$3 }' | sort
	@echo "-------------------------------------------|---------"
	@go tool cover -func=coverage.out | grep total | awk '{printf "  $(BOLD)%-40s | %s$(RESET)\n", "TOTAL", $$3}'
	@echo "-------------------------------------------|---------"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(CYAN)* [0.0%*] means no test files in package$(RESET)"

coverage-frontend: ## 🎨 Run frontend coverage (Full Table with Color)
	@echo "$(CYAN)📊 Frontend Coverage Detailed Report:$(RESET)"
	@cd frontend && FORCE_COLOR=1 npx vitest run --coverage --coverage.reporter=text --coverage.reporter=text-summary | tee coverage_summary.txt
	@echo ""

lint: ## 🔍 Run linters
	@echo "$(CYAN)🔍 Linting...$(RESET)"
	@go vet ./...
	@cd frontend && npm run lint

fmt: ## 🎯 Format code
	@gofmt -w -s .
	@echo "$(GREEN)✅ Formatted$(RESET)"

check: lint test ## ✅ Run all checks

# ============================================================================
# 📦 SETUP
# ============================================================================

install: ## 📦 Install dependencies
	@echo "$(CYAN)📦 Installing dependencies...$(RESET)"
	@go mod download
	@cd frontend && npm install
	@[ -f .env ] || cp .env.example .env
	@echo "$(GREEN)✅ Dependencies installed$(RESET)"

clean: stop-all ## 🧹 Clean everything
	@echo "$(YELLOW)🧹 Cleaning...$(RESET)"
	@rm -rf bin/ coverage.out coverage.html
	@cd frontend && rm -rf dist/ node_modules/.cache
	@docker compose down -v 2>/dev/null || true
	@echo "$(GREEN)✅ Clean complete$(RESET)"

# ============================================================================
# ❓ HELP
# ============================================================================

help: ## ❓ Show this help
	@echo ""
	@echo "$(BOLD)$(CYAN)🏛️  The Council$(RESET)"
	@echo "$(BOLD)════════════════════════════════════════════════════════════════$(RESET)"
	@echo ""
	@echo "$(BOLD)🔄 Lifecycle:$(RESET)"
	@echo "  $(CYAN)make start$(RESET)          Start everything"
	@echo "  $(CYAN)make stop$(RESET)           Stop everything"
	@echo "  $(CYAN)make restart$(RESET)        Restart everything"
	@echo "  $(CYAN)make status$(RESET)         Show service status"
	@echo ""
	@echo "$(BOLD)🐳 Docker:$(RESET)"
	@echo "  $(CYAN)make start-db$(RESET)       Start Postgres + Redis"
	@echo "  $(CYAN)make stop-db$(RESET)        Stop Docker services"
	@echo "  $(CYAN)make logs-db$(RESET)        Follow Docker logs"
	@echo "  $(CYAN)make reset-db$(RESET)       Reset database (⚠️ deletes data)"
	@echo ""
	@echo "$(BOLD)🔧 Backend:$(RESET)"
	@echo "  $(CYAN)make start-backend$(RESET)  Start Go server"
	@echo "  $(CYAN)make stop-backend$(RESET)   Stop Go server"
	@echo ""
	@echo "$(BOLD)🎨 Frontend:$(RESET)"
	@echo "  $(CYAN)make start-frontend$(RESET) Start React dev server"
	@echo "  $(CYAN)make stop-frontend$(RESET)  Stop React dev server"
	@echo ""
	@echo "$(BOLD)🏗️ Build & Test:$(RESET)"
	@echo "  $(CYAN)make build$(RESET)          Build for production"
	@echo "  $(CYAN)make test$(RESET)           Run all tests"
	@echo "  $(CYAN)make test-backend$(RESET)   Run backend tests"
	@echo "  $(CYAN)make test-frontend$(RESET)  Run frontend tests"
	@echo "  $(CYAN)make coverage$(RESET)        Run all coverage reports"
	@echo "  $(CYAN)make lint$(RESET)           Run linters"
	@echo "  $(CYAN)make check$(RESET)          Run all checks"
	@echo ""
	@echo "$(BOLD)📦 Setup:$(RESET)"
	@echo "  $(CYAN)make install$(RESET)        Install dependencies"
	@echo "  $(CYAN)make clean$(RESET)          Clean everything"
	@echo ""
