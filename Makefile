SHELL := /bin/zsh
.PHONY: setup dev build start test lint fmt typecheck db-up db-down

## Tools
PNPM := pnpm
TURBO := $(PNPM) turbo
GOLANGCI := golangci-lint
DOCKER_COMPOSE := docker compose -f infra/docker-compose.yml --env-file infra/.env

setup:
	$(PNPM) -v || corepack enable
	corepack prepare pnpm@9.12.0 --activate
	$(PNPM) install
	# Go tools (best-effort)
	@which golangci-lint >/dev/null 2>&1 || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0)
	@which air >/dev/null 2>&1 || (echo "Installing air..." && go install github.com/air-verse/air@latest)

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


