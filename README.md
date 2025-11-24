# credfolio

## Makefile Quick Reference

Available commands from repository root:
- `make setup`: Install JS/TS deps (pnpm) and Go tools (golangci-lint, air)
- `make dev`: Start database (Docker Compose) and run all apps in dev mode in parallel
- `make build`: Build all apps
- `make start`: Start all apps in production mode (where applicable)
- `make typecheck`: Run type checking across the monorepo
- `make test-backend`: Run Go tests in `apps/backend`
- `make test-frontend`: Run Vitest tests in `apps/frontend`
- `make test`: Run backend tests first, then frontend tests
- `make lint-backend`: Run `golangci-lint` in `apps/backend` (tidies modules first)
- `make lint-frontend`: Run Biome lint in `apps/frontend`
- `make lint`: Run backend lint first, then frontend lint
- `make fmt-backend`: Run `gofmt -s -w` and, if available, `goimports -w` in `apps/backend`
- `make fmt-frontend`: Run Biome format in `apps/frontend`
- `make fmt`: Format backend first, then frontend
- `make db-up`: Start Postgres and pgAdmin via Docker Compose (host ports: Postgres 55432, pgAdmin 8081)
- `make db-down`: Stop and remove DB containers and volumes