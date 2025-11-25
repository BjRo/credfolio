# Phase 0: Research & Technical Decisions

**Feature**: Generate Smart Profile & Credibility
**Date**: 2025-11-24

## 1. ORM & Migrations (GORM)

**Decision**: Use **GORM** (`gorm.io/gorm`) with **PostgreSQL** driver (`gorm.io/driver/postgres`).
**Rationale**: Explicit user request. GORM provides a full-featured ORM with associations, hooks, and migration support, fitting the requirement.
**Implementation Details**:
- **Driver**: `pgx` is recommended for performance, GORM uses it under the hood.
- **Migrations**: Use GORM's `AutoMigrate` for the initial phase as it fits the "generate" rapid dev cycle. For production versioning later, we might wrap it or use a separate tool, but for now `AutoMigrate` satisfies "migrations I want to use gorm.io".
- **Pattern**: Repository pattern in `internal/repository` (or `services/persistence`) to decouple business logic from GORM specific syntax (Clean Architecture).

## 2. AI Integration (OpenAI & Abstraction)

**Decision**: Use **`github.com/openai/openai-go`** (Official SDK) with a custom **Provider Interface**.
**Rationale**: The official SDK supports the latest features like **Structured Outputs** (native JSON schema parsing) which is a requirement.
**Abstraction Pattern**:
- Define an `LLMProvider` interface in the domain layer:
  ```go
  type LLMProvider interface {
      GenerateProfile(ctx context.Context, prompt string) (*ProfileData, error)
      // ... other methods
  }
  ```
- Implement `OpenAIProvider` struct that satisfies this interface.
- This allows swapping the implementation (e.g. to Anthropic or local LLM) without changing business logic.

## 3. PDF Generation (CV)

**Decision**: Use **`github.com/go-pdf/fpdf`** (community fork of gofpdf) or **`github.com/johnfercher/maroto`** (High-level wrapper).
**Rationale**: `maroto` is excellent for declarative layouts like CVs/Resumes. It simplifies grid-based designs compared to raw `fpdf`.
**Selection**: **Maroto v2**.

## 4. Text Extraction

**Decision**: Extract text from Markdown or txt files
**Rationale**: Pure Go library, no CGO dependencies. We need to extract text to send it to the LLM for processing.
**Flow**:
1. User uploads txt or markdown files.
2. Backend reads document.
3. Raw text is sent to OpenAI with a prompt to "Extract structured resume data".

## 5. OpenAPI Setup

**Decision**: Use **`github.com/swaggo/swag`** (declarative comments) or **Standard OpenAPI 3.0 YAML/JSON**.
**Rationale**: Go community often uses `swaggo` to generate docs from code comments. However, "OpenAPI-first" (ADR-0005) implies defining the spec first or at least having a strict contract.
**Adjustment**: Given ADR-0005 requires *generating frontend types*, having a `openapi.yaml` is crucial.
**Approach**: We will define the `openapi.yaml` in the `contracts/` directory as the source of truth (Design First), and use `oapi-codegen` (deepmap) or similar to generate Go server stubs and TypeScript clients. This enforces the contract better than comment-based generation.
**Tool**: `github.com/oapi-codegen/oapi-codegen/v2` for Go generation.

## 6. Architecture & Structure

- **Monorepo**: `apps/backend` (Go), `apps/frontend` (Next.js).
- **Auth**: "Current User" mock middleware. `r.WithContext(ctx)` to inject a hardcoded `userID`.
- **Security**: None for now (skipped per instructions).

## 7. Dependencies Checklist

- [ ] `gorm.io/gorm`
- [ ] `gorm.io/driver/postgres`
- [ ] `github.com/openai/openai-go`
- [ ] `github.com/johnfercher/maroto/v2`
- [ ] `github.com/ledongthuc/pdf`
- [ ] `github.com/oapi-codegen/oapi-codegen/v2` (CLI tool)

## 8. Take things step by step
- Commit every finished task individually to the git branch
- make sure to run `make test`, `make lint` and `make fmt` before and fix any errors


