# 0001: Monorepo structure with Turborepo

- Status: accepted
- Date: 2025-11-16

## Context
We want fast feedback loops and agent-friendly structure across frontend and backend.

## Decision
Use a Turborepo monorepo with `apps/frontend` (Next.js) and `apps/backend` (Go).

## Consequences
- Single command orchestration for dev/test/lint/typecheck.
- Clear separation of concerns; easy CI reuse.


