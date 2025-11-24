# Implementation Plan: Generate Smart Profile & Credibility

**Branch**: `001-generate-smart-profile` | **Date**: 2025-11-24 | **Spec**: [specs/001-generate-smart-profile/spec.md](../spec.md)
**Input**: Feature specification from `/specs/001-generate-smart-profile/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

The "Generate Smart Profile & Credibility" feature automates the creation of professional profiles by extracting data from reference letters using AI (OpenAI). It includes a "Credibility" section derived from employer sentiment and allows users to tailor their profiles and CVs (PDF) for specific job descriptions.

## Technical Context

**Language/Version**: Go 1.23+ (Backend), TypeScript/Next.js 15+ (Frontend)
**Primary Dependencies**:
- **Backend**: `chi` (Router), `gorm` (ORM), `openai-go` (AI), `maroto` (PDF Gen), `ledongthuc/pdf` (PDF Parse), `oapi-codegen` (OpenAPI)
- **Frontend**: `react`, `next`, `tailwindcss`
**Storage**: PostgreSQL (via GORM)
**Testing**: Go standard `testing`, `testify` (Backend); `vitest`, `react-testing-library` (Frontend)
**Target Platform**: Web (Linux container)
**Project Type**: Monorepo (`apps/backend`, `apps/frontend`)
**Performance Goals**: Profile generation < 60s.
**Constraints**: No auth/security implemented yet (Mock User).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Definition of Done**: Will run lint/test/fmt make commands mentioned in constitution.
- **ADR**:
  - Need new ADR for GORM (as requested).
  - Need new ADR for Maroto/PDF? (Maybe, "New libraries or tools").
  - Using OpenAPI (ADR-0005 complient).
- **Context7**: Will use context7 for GORM/OpenAI generation.

## Project Structure

### Documentation (this feature)

```text
specs/001-generate-smart-profile/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
apps/backend/
├── api/
│   └── openapi.yaml             # API Contract
├── internal/
│   ├── domain/                  # Core entities (Profile, WorkExperience)
│   ├── service/                 # Business Logic (ProfileService, AIService)
│   ├── repository/              # Data Access (GormProfileRepo)
│   └── handler/                 # HTTP Handlers (generated or Chi)
├── pkg/
│   ├── ai/                      # OpenAI Abstraction
│   └── pdf/                     # PDF Generation/Extraction
└── cmd/server/

apps/frontend/
├── src/
│   ├── app/
│   │   └── profile/             # Pages
│   ├── components/
│   │   └── profile/             # UI Components
│   └── lib/
│       └── api/                 # Generated API Client
```

**Structure Decision**: Standard Clean Architecture within the Monorepo backend. Feature-based grouping in Frontend.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| GORM      | User Request | Using raw SQL or `sqlc` rejected per user constraint. |
| OpenAI Abstraction | Future Proofing | Direct API calls would make switching models harder later (User Request). |
