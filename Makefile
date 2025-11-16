SHELL := /bin/bash
.PHONY: setup dev build start test test-backend test-frontend lint lint-backend lint-frontend fmt fmt-backend fmt-frontend typecheck db-up db-down

## Tools
PNPM := pnpm
TURBO := $(PNPM) turbo
GO := go
GOPATH := $(shell $(GO) env GOPATH 2>/dev/null)
GOBIN := $(shell $(GO) env GOBIN 2>/dev/null)
BIN_DIR := $(if $(GOBIN),$(GOBIN),$(GOPATH)/bin)
GOLANGCI := $(BIN_DIR)/golangci-lint
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
		echo "Ensuring golangci-lint is installed with local Go toolchain..."; \
		GOTOOLCHAIN=local $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
		command -v air >/dev/null 2>&1 || (echo "Installing air..." && go install github.com/air-verse/air@latest); \
	else \
		echo "Go not found; skipping Go tool installs (golangci-lint, air)."; \
	fi

dev:
	$(TURBO) run dev --parallel

dev-db:
	$(MAKE) db-up
	$(MAKE) dev

build:
	$(TURBO) run build

start:
	$(TURBO) run start

typecheck:
	$(TURBO) run typecheck

test-backend:
	(cd apps/backend && GOTOOLCHAIN=local go test ./...)

test-frontend:
	$(TURBO) run test --filter=@credfolio/frontend

test:
	$(MAKE) test-backend && $(MAKE) test-frontend

lint-backend:
	(cd apps/backend && GOTOOLCHAIN=local go mod tidy && GOTOOLCHAIN=local go mod download && GOTOOLCHAIN=local $(GOLANGCI) run)

lint-frontend:
	$(TURBO) run lint --filter=@credfolio/frontend

lint:
	$(MAKE) lint-backend && $(MAKE) lint-frontend

fmt-backend:
	(cd apps/backend && gofmt -s -w . && test -n "$$(command -v goimports)" && goimports -w . || true)

fmt-frontend:
	$(TURBO) run fmt --filter=@credfolio/frontend

fmt:
	$(MAKE) fmt-backend && $(MAKE) fmt-frontend

db-up:
	$(DOCKER_COMPOSE) up -d

db-down:
	$(DOCKER_COMPOSE) down -v


