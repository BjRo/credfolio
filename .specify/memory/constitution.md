<!--
Sync Impact Report:
Version change: 1.3.0 → 1.4.0 (MINOR: new principle added)
Modified principles: N/A
Added sections: Structured Outputs and Personas in LLM Prompts principle
Removed sections: N/A
Templates requiring updates:
  ✅ plan-template.md - Constitution Check section aligns with all principles
  ✅ spec-template.md - No changes needed (user stories already align with DoD)
  ✅ tasks-template.md - No changes needed (principle applies to LLM interaction, not task structure)
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
   - run `make lint`

2. **Tests pass**:
   - run `make test`

3. **Code is properly formatted**:
   - run `make fmt`

ALWAYS RUN THESE COMMAND FROM THE ROOT PROJECT DIRECTORY. NEVER FALL DOWN TO OTHER TOOLS. NEVER SKIP THESE. ALWAYS RUN THESE BEFORE FINISHING A TASK

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

### IV. Unit-Testing-First

We adopt a unit-testing-first approach where unit tests guide our development process. This ensures fast feedback loops, high code quality, and maintainable test suites that don't leak implementation details.

**Core Requirements:**

- **Unit Tests Guide Development**: Unit tests drive the development process, not follow it.
- **Mandatory Unit Tests**: All developments MUST be accompanied by unit tests (both in frontend and backend). No code is considered complete without corresponding unit tests.
- **Preference for Unit Tests**: We favor unit tests over integration tests because they enable fast feedback loops and isolate behavior.
- **AAA Style**: All tests MUST be written in AAA (Arrange-Act-Assert) style for clarity and consistency.
- **Implementation Hiding**: Tests MUST be written to avoid leaking implementation details whenever possible. Tests should verify behavior, not internal structure.
- **Naming Convention**: Test names MUST include context, trigger, and expectation (e.g., `test_UserService_whenEmailIsInvalid_returnsError`).
- **No External Calls**: Tests MUST NOT call out to external systems via HTTP/HTTPS. External calls MUST be stubbed and simulated with test data.

**Rationale**: Unit-testing-first development ensures code correctness from the start, enables rapid iteration with fast feedback, prevents regressions, and creates living documentation of system behavior. By avoiding external dependencies in tests, we maintain test speed and reliability while keeping tests focused on the unit under test.

### V. Structured Outputs and Personas in LLM Prompts

When interacting with Large Language Models (LLMs) for code generation, analysis, or development assistance, we MUST use structured outputs and personas to improve the quality and consistency of results.

**Core Requirements:**

- **Structured Outputs**: LLM prompts MUST request structured outputs (e.g., JSON schemas, specific formats, or well-defined response structures) when the response needs to be parsed, validated, or used programmatically. This ensures predictable, machine-readable results.
- **Personas**: LLM prompts MUST include appropriate personas or role definitions that guide the model's behavior and expertise level. Personas help the model adopt the right perspective, knowledge domain, and communication style for the task.
- **Context-Specific Personas**: Personas MUST be tailored to the specific task (e.g., "senior backend engineer" for API design, "security expert" for vulnerability analysis, "code reviewer" for review tasks).
- **Structured Prompt Design**: Prompts MUST be organized with clear sections: context, role/persona, task description, expected output format, and constraints.

**When to Apply:**

- Any LLM interaction where structured, consistent output is needed

**Rationale**: Structured outputs reduce parsing errors, enable automated validation, and ensure consistent response formats across interactions. Personas improve the relevance and quality of LLM responses by providing appropriate context, expertise level, and behavioral guidance. Together, these techniques significantly improve the reliability and usefulness of LLM-assisted development.

## Governance

This constitution supersedes all other development practices and guidelines. All PRs and code reviews MUST verify compliance with all principles, including Definition of Done, Architecture Decision Records, Context7 API Verification, Unit-Testing-First, and Structured Outputs and Personas in LLM Prompts requirements.

**Amendment Procedure**: Amendments require documentation of rationale, version bump according to semantic versioning (MAJOR for backward-incompatible changes, MINOR for new principles/sections, PATCH for clarifications), and update of dependent templates.

**Compliance Review**: The Constitution Check gate in implementation plans must verify adherence to all principles before proceeding with feature work.

**Version**: 1.4.0 | **Ratified**: TODO(RATIFICATION_DATE) | **Last Amended**: 2025-01-27
