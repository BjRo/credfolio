# 0004: Unit Testing Strategy

- Status: accepted
- Date: 2025-11-24

## Context

To ensure high code quality, maintainability, and fast feedback loops during development, we need a clear and consistent strategy for testing our codebase. We want to avoid brittle tests and tests that are slow to run.

## Decision

We will adopt a unit-testing-first approach with the following constraints and guidelines:

*   **Unit Tests Guide Development:** We use unit tests to guide our development.
*   **Mandatory Unit Tests:** All developments need to be accompanied by unit tests (both in frontend and backend).
*   **Preference for Unit Tests:** We favor unit tests over integration tests (because they enable fast feedback loops).
*   **AAA Style:** All of our tests are written in AAA (Arrange Act Assert) style.
*   **Implementation Hiding:** Whenever possible we write tests in such a way that they don't leak implementation details.
*   **Naming Convention:** We use good naming for tests: Test names should include context, trigger and expectation.
*   **No External Calls:** We don't call out to external systems via HTTP/HTTPs. We stub and simulate these calls with test data.

## Consequences

*   **Positive:**
    *   Faster feedback loops for developers.
    *   Higher confidence in code correctness.
    *   Tests serve as documentation for the expected behavior.
    *   Reduced flakiness by avoiding external dependencies in tests.
*   **Negative:**
    *   Requires discipline and initial time investment to write comprehensive tests.
    *   Mocking complex external systems can sometimes be non-trivial.

