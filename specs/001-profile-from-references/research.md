# Research: Profile Generation from References

**Status**: Complete
**Date**: 2025-11-23

## Decisions

### 1. PDF & Document Extraction

**Context**: We need to extract text from User uploaded PDF and DOCX files to feed into the LLM for parsing.
**Decision**: Use `ledongthuc/pdf` (Go) for PDF text extraction and `baliance/gooxml` (or similar) for DOCX.
**Rationale**:
- `ledongthuc/pdf` is a pure Go library, avoiding CGO dependencies (easier deployment).
- It supports text extraction sufficient for reference letters.
- OCR is out of scope for MVP (as per spec).
**Alternatives Considered**:
- `unidoc/unipdf`: Powerful but has strict licensing (AGPL/Commercial).
- `dslipak/pdf`: Requires CGO (pdfium).

### 2. LLM Provider & Integration

**Context**: We need to extract structured data (Role, Company, Skills, Sentiment) from unstructured text.
**Decision**: OpenAI API (GPT-4o) via `sashabaranov/go-openai` library.
**Rationale**:
- Best-in-class performance for unstructured text extraction and JSON formatting.
- `sashabaranov/go-openai` is the standard, well-maintained community client.
- "Employer Feedback" extraction requires high semantic understanding (sentiment/nuance), which GPT-4 excels at.
**Alternatives Considered**:
- Local LLaMA: Too resource-intensive for MVP hosting.
- Anthropic Claude: Good alternative, but OpenAI has better structured output (JSON mode) support currently.

### 3. Database Driver

**Context**: Persisting User Profiles and extracted data in Postgres.
**Decision**: `jackc/pgx/v5`.
**Rationale**:
- High performance, standard for modern Go.
- Better feature set than `lib/pq` (which is in maintenance mode).
**Alternatives Considered**:
- `lib/pq`: Older, slower, less features.
- GORM: Adds overhead/complexity; raw SQL/sqlc is preferred for performance and clarity in this project.

### 4. File Storage

**Context**: Storing original reference letters.
**Decision**: Local filesystem (MVP) behind an interface `FileStorage`.
**Rationale**:
- Keeps infrastructure simple for MVP.
- Interface allows easy swap to S3 for production.

