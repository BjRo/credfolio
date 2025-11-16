# 0003 – Environment Variable Management in Monorepo

Date: 2025-11-16

## Status
Accepted

## Context
Both the frontend (Next.js) and backend (Go) require configuration (e.g., API URLs, ports). We need a consistent approach across development and production environments that:

- Works locally without container orchestration
- Avoids committing secrets
- Plays well with the monorepo structure
- Keeps public vs private variables clear for the frontend

## Decision
- Frontend (Next.js):
  - Use `.env.local` for development overrides (git-ignored). Use `.env` for shared non-secret defaults as needed.
  - Use the `NEXT_PUBLIC_` prefix for variables that may be exposed to the browser.
  - Example: `NEXT_PUBLIC_BACKEND_URL=http://localhost:8080`
  - Next.js automatically loads `.env*, .env.local`.

- Backend (Go):
  - Read environment variables from the OS by default.
  - In development, automatically load `.env` and `.env.local` if present using `github.com/joho/godotenv`. This makes `pnpm run dev`/`go run .` pick up local overrides without extra tooling.
  - Example: `PORT=8080`

- Infra (Docker Compose):
  - Keep a separate `infra/.env` for Compose-only variables (DB credentials), not loaded by app processes.

- Versioning:
  - Provide `*.example` files to document required variables without committing secrets.
  - Example files:
    - `apps/frontend/.env.example` (or `.env.local.example`)
    - `apps/backend/.env.example`
    - `infra/.env.example`

## Consequences
- Developers can start both apps with `pnpm run dev` using local `.env.local` files.
- Production deployments set real environment variables via the platform (e.g., systemd, Docker, cloud provider), with no reliance on `.env.local`.
- Public/private boundary is explicit for frontend via the `NEXT_PUBLIC_` prefix.

## Example Variables
Frontend (`apps/frontend/.env.local`):
```\nNEXT_PUBLIC_BACKEND_URL=http://localhost:8080\n```\n
Backend (`apps/backend/.env.local`):
```\nPORT=8080\n```\n


