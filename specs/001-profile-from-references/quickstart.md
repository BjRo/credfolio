# Quickstart: Profile Generation from References

## Prerequisites

1. **OpenAI API Key**: You need a valid `OPENAI_API_KEY` set in `apps/backend/.env`.
2. **Database**: Ensure Postgres is running via `make db-up`.

## Running the Application

1. Start the development environment:
   ```bash
   make dev
   ```
2. Access the frontend at `http://localhost:3000`.
3. Access the backend API at `http://localhost:8080`.

## Testing the Feature

### Manual Testing (UI)
1. Go to `http://localhost:3000/upload`.
2. Select a PDF reference letter.
3. Click "Upload".
4. Wait for processing (simulated or real).
5. Verify the profile is populated at `http://localhost:3000/profile`.

### API Testing

**Upload a file:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@/path/to/reference.pdf"
```

**Get Profile:**
```bash
curl http://localhost:8080/api/profile
```

**Tailor Profile:**
```bash
curl -X POST http://localhost:8080/api/profile/tailor \
  -H "Content-Type: application/json" \
  -d '{"job_description": "Looking for a Senior Go Engineer..."}'
```

