## Makefile quick reference

Run these from the repository root.

- setup: Install JS/TS deps (pnpm) and Go tools (golangci-lint, air).
- dev: Start database (Docker Compose) and run all apps in dev mode in parallel.
- build: Build all apps.
- start: Start all apps in production mode (where applicable).
- typecheck: Run type checking across the monorepo.
- test-backend: Run Go tests in `apps/backend`.
- test-frontend: Run Vitest tests in `apps/frontend`.
- test: Run backend tests first, then frontend tests.
- lint-backend: Run `golangci-lint` in `apps/backend` (tidies modules first).
- lint-frontend: Run Biome lint in `apps/frontend`.
- lint: Run backend lint first, then frontend lint.
- fmt-backend: Run `gofmt -s -w` and, if available, `goimports -w` in `apps/backend`.
- fmt-frontend: Run Biome format in `apps/frontend`.
- fmt: Format backend first, then frontend.
- db-up: Start Postgres and pgAdmin via Docker Compose (host ports: Postgres 55432, pgAdmin 8081).
- db-down: Stop and remove DB containers and volumes.

## Definition of done for new features

For every new feature or change, ensure all of the following succeed locally before opening a PR:

1) Lint is clean for both backend and frontend (no errors):
   - make lint-backend
   - make lint-frontend
   - or: make lint

2) Tests pass for both backend and frontend:
   - make test-backend
   - make test-frontend
   - or: make test

3) Code is properly formatted:
   - make fmt-backend
   - make fmt-frontend
   - or: make fmt

If any command fails, fix the reported issues before proceeding.

