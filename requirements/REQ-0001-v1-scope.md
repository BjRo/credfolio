# REQ-0001: Credfolio v1 Scope

- Status: planned
- Last Updated: 2025-11-16

## Context
First iteration focuses on infrastructure and minimal running app.

## Requirement
- Monorepo with frontend and backend
- Local Postgres via Docker
- Scripts for dev, test, lint, typecheck

## Acceptance Criteria
- `make dev` runs both apps
- `make db-up` brings up Postgres
- `make test`/`make lint` pass cleanly


