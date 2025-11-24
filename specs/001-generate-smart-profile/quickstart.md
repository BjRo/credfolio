# Quickstart: Generate Smart Profile

## Prerequisites

- **Go 1.23+**
- **Node.js 20+**
- **Docker** (for Postgres)
- **OpenAI API Key**

## 1. Infrastructure Setup

Create `infra/.env`:

```ini
POSTGRES_USER=credfolio
POSTGRES_PASSWORD=credfolio
POSTGRES_DB=credfolio
PGADMIN_DEFAULT_EMAIL=admin@example.com
PGADMIN_DEFAULT_PASSWORD=adminadmin
```

Start PostgreSQL and pgAdmin:

```bash
make db-up
```

This starts:
- PostgreSQL on port `55432` (host) → `5432` (container)
- pgAdmin on port `8081` (host) → `80` (container)

## 2. Environment Variables

Create `apps/backend/.env`:

```ini
DATABASE_URL=postgres://credfolio:credfolio@localhost:55432/credfolio?sslmode=disable
OPENAI_API_KEY=sk-...
PORT=8080
```

Create `apps/frontend/.env.local`:

```ini
NEXT_PUBLIC_API_BASE=http://localhost:8080
```

## 3. Project Setup

Install dependencies:

```bash
make setup
```

This installs:
- JavaScript/TypeScript dependencies (pnpm)
- Go tools (golangci-lint, air)

## 4. Backend Setup

Install Go dependencies:

```bash
cd apps/backend
go mod download
```

Note: Code generation (OpenAPI server stubs) is already done. If you need to regenerate after modifying `apps/backend/api/openapi.yaml`, run:

```bash
go generate ./...
```

Run Backend (from repo root):

```bash
make dev
```

Or run backend with database startup:

```bash
make dev-db
```

The backend will be available at `http://localhost:8080`

## 5. Frontend Setup

The frontend dependencies are installed by `make setup`. To run the frontend:

```bash
make dev
```

The frontend will be available at `http://localhost:3000`

## 6. Verification

1. Check backend health: `curl http://localhost:8080/healthz`
2. Check frontend: Open `http://localhost:3000` in your browser
3. Test API endpoints using the frontend UI or curl:
   - Upload reference letter: `POST http://localhost:8080/reference-letters`
   - Generate profile: `POST http://localhost:8080/profile/generate`
   - Get profile: `GET http://localhost:8080/profile`

