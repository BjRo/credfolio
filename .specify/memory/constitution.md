<!--
Sync Impact Report:
Version change: 1.1.0 → 1.2.0 (MINOR: new principle added)
Modified principles: N/A
Added sections: Context7 API Verification principle
Removed sections: N/A
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section aligns with all principles
  ✅ spec-template.md - No changes needed (user stories already align with DoD)
  ✅ tasks-template.md - No changes needed (task structure aligns with DoD)
  ⚠ pending: Manual review recommended for any custom command templates
Follow-up TODOs:
  - TODO(RATIFICATION_DATE): Original adoption date unknown - needs historical research or project start date
-->

# credfolio Constitution

## Core Principles

### I. Definition of Done (NON-NEGOTIABLE)

For every new feature or code change, ensure all of the following succeed locally before calling a feature done.

THIS IS NOT NEGOTIABLE.

For changes that affect only the backend or only the frontend, you may run the commands only for the affected side, unless you suspect cross-cutting impact.

ALWAYS RUN THE VERIFICATION TOOLS FROM THE PROJECT ROOT.

**Mandatory Verification Checklist:**

1. **Lint is clean** (no errors):
   - Backend: run `make lint-backend`
   - Frontend: run `make lint-frontend`

2. **Tests pass**:
   - Backend: run `make test-backend`
   - Frontend: run `make test-frontend`

3. **Code is properly formatted**:
   - Backend: run `make fmt-backend`
   - Frontend: run `make fmt-frontend`

**Rationale**: Ensures code quality, consistency, and prevents regressions before code review. All PRs MUST pass these gates before submission. If any command fails, fix the reported issues before proceeding.

### II. Architecture Decision Records (ADR)

All technical decisions that have a long-term, cross-cutting impact on the system MUST be documented as Architecture Decision Records (ADRs). The goal is to maintain a clear historical record of why decisions were made, enabling future understanding and informed evolution of the system.

**When to Create an ADR:**

- New libraries or tools are introduced
- Previous libraries or tools are phased out and removed
- New patterns are introduced
- Old patterns are phased out and removed

**ADR Creation Requirements:**

1. Use the template found in `adr/templates/adr.md`
2. Create a new ADR file with an incremented number (e.g., `0004-<title>.md`)
3. Fill out all placeholders in the template (ID, Title, Status, Date, Context, Decision, Consequences)
4. Add the new ADR entry to `adr/index.md` with all required metadata

**Rationale**: ADRs provide institutional memory for architectural decisions, preventing repeated debates and enabling informed future changes. They document the context, decision rationale, and consequences, making it easier to understand system evolution and evaluate when decisions should be revisited.

### III. Context7 API Verification

When working with external APIs (frontend or backend libraries), you MUST use context7 to verify correct API usage when uncertainty exists. This ensures code is generated and fixed according to official library documentation rather than assumptions or outdated patterns.

**When to Use Context7:**

- Fixing linting errors related to a library API
- Fixing tests that fail due to incorrect API usage
- Generating new code for a library previously not used in the codebase
- Any situation where API usage is unclear or ambiguous

**Usage Requirement:**

Append `use context7` to your request so the related MCP (Model Context Protocol) is included and the actual library documentation can be consulted. This provides access to up-to-date, authoritative API documentation.

**Rationale**: Incorrect API usage leads to bugs, test failures, and technical debt. Consulting official documentation via context7 ensures code correctness from the start, reduces debugging time, and maintains alignment with library best practices and current versions.

## Governance

This constitution supersedes all other development practices and guidelines. All PRs and code reviews MUST verify compliance with all principles, including Definition of Done, Architecture Decision Records, and Context7 API Verification requirements.

**Amendment Procedure**: Amendments require documentation of rationale, version bump according to semantic versioning (MAJOR for backward-incompatible changes, MINOR for new principles/sections, PATCH for clarifications), and update of dependent templates.

**Compliance Review**: The Constitution Check gate in implementation plans must verify adherence to all principles before proceeding with feature work.

**Version**: 1.2.0 | **Ratified**: TODO(RATIFICATION_DATE) | **Last Amended**: 2025-01-27
