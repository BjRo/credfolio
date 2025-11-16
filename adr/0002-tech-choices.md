# 0002: Tech choices (Next.js, Tailwind, Go, tsgo)

- Status: accepted
- Date: 2025-11-16

## Context
We prefer fast, strong tooling and agent-friendly stacks.

## Decision
- Frontend: Next.js (App Router) + Tailwind
- Backend: Go with chi, `golangci-lint`, `air`
- TS compiler: `tsgo` (`@typescript/native-preview`) for typecheck; Next emits JS via SWC
- Lint/format: Biome for JS/TS

## Consequences
- Very fast typechecks and lints; modern DX
- Keep `tsc` fallback available if tsgo limitations are hit


