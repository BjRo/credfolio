# Quickstart: Generate Smart Profile

## Prerequisites

- **Go 1.23+**
- **Node.js 20+** with pnpm
- **Docker** (for PostgreSQL)
- **OpenAI API Key**

## 1. Initial Setup

Clone and install dependencies:

```bash
# Install all dependencies (Go tools, pnpm packages)
make setup
```

## 2. Environment Variables

Create `infra/.env` for the database:

```ini
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=credfolio
```

Create `apps/backend/.env` for the backend:

```ini
DATABASE_URL="host=localhost user=postgres password=postgres dbname=credfolio port=5432 sslmode=disable"
OPENAI_API_KEY="sk-your-api-key-here"
PORT=8080
```

## 3. Start Development Environment

Start everything (database + backend + frontend):

```bash
make dev-db
```

Or start components individually:

```bash
# Start database only
make db-up

# Start dev servers (frontend + backend)
make dev
```

## 4. Run Tests

```bash
# Run all tests
make test

# Run backend tests only
make test-backend

# Run frontend tests only
make test-frontend
```

## 5. Verification

1. **Frontend**: Open `http://localhost:3000` - you should see the Credfolio home page
2. **Backend API**: Open `http://localhost:8080/api/v1/profile` - should return 401 (unauthorized)
3. **Upload a reference letter**: Navigate to `/profile/generate` and upload a `.txt` or `.md` file
4. **Generate profile**: Click "Generate Profile with AI" to process the uploaded letters
5. **View profile**: See your generated profile at `/profile`
6. **Tailor profile**: Go to `/profile/tailor` and paste a job description
7. **Download CV**: Click "Download CV" to generate a PDF

## Troubleshooting

### Database connection issues

```bash
# Restart database
make db-down
make db-up
```

### Port already in use

Check if ports 3000 (frontend), 5432 (database), or 8080 (backend) are in use:

```bash
lsof -i :3000
lsof -i :5432
lsof -i :8080
```

### OpenAI API errors

Ensure your `OPENAI_API_KEY` is valid and has sufficient credits.

