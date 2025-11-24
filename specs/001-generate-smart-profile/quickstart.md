# Quickstart: Generate Smart Profile

## Prerequisites

- **Go 1.23+**
- **Node.js 20+**
- **Docker** (for Postgres)
- **OpenAI API Key**

## 1. Infrastructure Setup

Start PostgreSQL:

```bash
docker run --name credfolio-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15
```

Create the database:

```bash
docker exec -it credfolio-db createdb -U postgres credfolio
```

## 2. Environment Variables

Create `apps/backend/.env`:

```ini
DATABASE_URL="host=localhost user=postgres password=postgres dbname=credfolio port=5432 sslmode=disable"
OPENAI_API_KEY="sk-..."
PORT=8080
```

## 3. Backend Setup

Install dependencies & Generate Code:

```bash
cd apps/backend
go mod download
# If using oapi-codegen
go generate ./...
```

Run Backend:

```bash
make dev-db
```

## 4. Frontend Setup

Install dependencies:

```bash
pnpm install
```

Run Frontend:

```bash
make dev-db
```

## 5. Verification

1. Open `http://localhost:8080/swagger/index.html` to see the API Docs.
2. Use the "Upload Reference Letter" endpoint to test PDF upload.
3. Check `http://localhost:3000` for the frontend.

