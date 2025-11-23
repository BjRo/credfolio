# Implementation Plan: Profile Generation from References

**Branch**: `001-profile-from-references` | **Date**: 2025-11-23 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-profile-from-references/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

The system will allow users to upload reference letters (PDF/DOCX), extracting structured data (Work Experience, Skills, Employer Feedback) using an LLM-based approach. This data will populate a professional profile ("Credfolio") that users can edit, view, and download as a CV.

## Technical Context

**Language/Version**: Go 1.22.5 (Backend), TypeScript 5.6 (Frontend)
**Primary Dependencies**:
- Backend: `chi` (Router), `pgx` (Database), `go-openai` (LLM), `ledongthuc/pdf` & `unidoc` (PDF Extraction - TBD)
- Frontend: Next.js 14, TailwindCSS, React Query (presumed)
**Storage**: Postgres 16 (Data), Local Filesystem (Uploads - MVP)
**Testing**: `go test` (Backend), `vitest` (Frontend)
**Target Platform**: Linux server (Dockerized)
**Project Type**: Web Application (Monorepo)
**Performance Goals**: Profile generation < 30s
**Constraints**: MVP English-only, Text-based PDFs
**Scale/Scope**: User-focused, MVP scale

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Lint is clean**: Verified via Makefile
- [x] **Tests pass**: Verified via Makefile
- [x] **ADR Required**: Yes, for PDF Library and LLM Provider choices. (Will include decision in research.md)
- [x] **Context7**: Required for new libraries (OpenAI, PDF lib).

## Project Structure

### Documentation (this feature)

```text
specs/001-profile-from-references/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── openapi.yaml
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
apps/backend/
├── src/
│   ├── api/
│   │   ├── handlers/      # HTTP Handlers
│   │   └── router.go      # Route definitions
│   ├── services/
│   │   ├── extractor/     # PDF/LLM extraction logic
│   │   └── profile/       # Profile management logic
│   ├── db/
│   │   ├── migrations/    # SQL migrations
│   │   └── queries/       # DB access
│   └── models/            # Go structs
└── tests/

apps/frontend/
├── app/
│   ├── dashboard/         # Profile view
│   └── upload/            # Upload flow
├── components/
│   ├── profile/           # Profile UI components
│   └── upload/            # File upload components
└── services/
    └── api.ts             # API client
```

**Structure Decision**: Standard Clean Architecture for Go backend; App Router for Next.js frontend.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | | |
