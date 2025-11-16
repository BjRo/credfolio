SHELL := /bin/zsh
.PHONY: setup dev build start test lint fmt typecheck db-up db-down

## Tools
PNPM := pnpm
TURBO := $(PNPM) turbo
GOLANGCI := golangci-lint
DOCKER_COMPOSE := docker compose -f infra/docker-compose.yml --env-file infra/.env

setup:
	@if command -v corepack >/dev/null 2>&1; then \
		($(PNPM) -v >/dev/null 2>&1 || corepack enable) && corepack prepare pnpm@9.12.0 --activate; \
	else \
		echo "Corepack not found. Attempting to install pnpm via npm..."; \
		if command -v npm >/dev/null 2>&1; then npm i -g pnpm@9.12.0; else echo "npm not found. Please install Node.js (>=18) and retry."; exit 1; fi; \
	fi
	$(PNPM) install
	# Go tools (best-effort, only if Go is available)
	@if command -v go >/dev/null 2>&1; then \
		if ! command -v $(GOLANGCI) >/dev/null 2>&1; then \
			echo "Installing golangci-lint..."; \
			if command -v curl >/dev/null 2>&1; then \
				BIN_DIR="$$(go env GOPATH 2>/dev/null)/bin"; \
				[ -n "$$BIN_DIR" ] || BIN_DIR="$$HOME/go/bin"; \
				mkdir -p "$$BIN_DIR"; \
				curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$$BIN_DIR" v1.61.0; \
			else \
				echo "curl not found; skipping golangci-lint install"; \
			fi; \
		fi; \
		command -v air >/dev/null 2>&1 || (echo "Installing air..." && go install github.com/air-verse/air@latest); \
	else \
		echo "Go not found; skipping Go tool installs (golangci-lint, air)."; \
	fi

dev:
	$(MAKE) db-up
	$(TURBO) run dev --parallel

build:
	$(TURBO) run build

start:
	$(TURBO) run start

typecheck:
	$(TURBO) run typecheck

test:
	$(TURBO) run test && (cd apps/backend && go test ./...)

lint:
	$(TURBO) run lint && (cd apps/backend && $(GOLANGCI) run)

fmt:
	$(TURBO) run fmt && (cd apps/backend && gofmt -s -w . && test -n "$$(command -v goimports)" && goimports -w . || true)

db-up:
	$(DOCKER_COMPOSE) up -d

db-down:
	$(DOCKER_COMPOSE) down -v


