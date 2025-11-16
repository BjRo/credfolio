# REQ-0002: CI for Pull Requests

- Status: planned
- Last Updated: 2025-11-16

## Context

Credfolio uses GitHub and a monorepo containing frontend and backend code, with a shared `Makefile` exposing common developer workflows (e.g. `make test-backend`, `make test-frontend`, `make lint-backend`, `make lint-frontend`, `make fmt-backend`, `make fmt-frontend`).

To keep code quality high and reduce regressions, we want an automated quality gate on pull requests that enforces that tests, linting, and formatting are all clean before code is merged. This quality gate should be implemented using GitHub Actions and should integrate with GitHub’s pull request checks so that merges are blocked when the gate fails.

## Requirement

- GitHub Actions MUST run automatically for all pull requests targeting the main development branch (e.g. `main`) on every push to a branch that has an open pull request.
- The CI configuration MUST define one or more workflows that, together, perform at least the following checks for each pull request:
  - Backend tests MUST be executed via the canonical make target (currently `make test-backend`).
  - Frontend tests MUST be executed via the canonical make target (currently `make test-frontend`).
  - Backend linting MUST be executed via the canonical make target (currently `make lint-backend`) and the workflow MUST fail if lint errors are present.
  - Frontend linting MUST be executed via the canonical make target (currently `make lint-frontend`) and the workflow MUST fail if lint errors are present.
  - Backend formatting MUST be verified using the canonical formatting command (currently `make fmt-backend` or an equivalent check-only variant) such that the workflow fails if backend files are not properly formatted.
  - Frontend formatting MUST be verified using the canonical formatting command (currently `make fmt-frontend` or an equivalent check-only variant) such that the workflow fails if frontend files are not properly formatted.
- All of the above checks MUST surface their status back to the pull request as GitHub Checks, and the pull request MUST be marked as failing (red) when any check fails.
- The main development branch (e.g. `main`) MUST be configured (via repository branch protection settings) so that the CI checks defined above MUST pass before a pull request can be merged.
- CI workflows SHOULD structure jobs such that backend and frontend checks run in parallel where feasible to keep overall runtime reasonable.
- CI workflows SHOULD use caching (e.g. for Node.js and Go modules) when possible to keep typical pull request runs under a reasonable time budget (e.g. under 15 minutes under normal load).
- CI workflows SHOULD use clear job and step names (e.g. “Backend Tests”, “Frontend Lint”) so that failures are easy for developers to understand.
- CI workflows MUST run using the principle of least privilege for GitHub tokens and secrets (e.g. `GITHUB_TOKEN` permissions limited to what is required for checks and status reporting).

## Acceptance Criteria

- Creating a pull request against the main development branch triggers a GitHub Actions workflow run without any manual intervention.
- When a backend test fails (e.g. `make test-backend` exits with a non-zero status), the corresponding CI job fails and the overall pull request status is shown as failing, preventing merge.
- When a frontend test fails (e.g. `make test-frontend` exits with a non-zero status), the corresponding CI job fails and the overall pull request status is shown as failing, preventing merge.
- When backend linting fails (e.g. `make lint-backend` reports lint errors), the corresponding CI job fails and the overall pull request status is shown as failing, preventing merge.
- When frontend linting fails (e.g. `make lint-frontend` reports lint errors), the corresponding CI job fails and the overall pull request status is shown as failing, preventing merge.
- When backend or frontend files are not properly formatted according to `make fmt-backend` / `make fmt-frontend` (or their check-only equivalents), the CI workflow detects this and fails, preventing merge until formatting is corrected.
- When all tests, linting, and formatting checks pass for both backend and frontend, the GitHub Actions checks on the pull request are green, and the protected branch settings allow the pull request to be merged.

## Notes

- This requirement intentionally focuses on CI for pull requests (tests, lint, formatting) and does not specify any deployment (CD) behavior.
- Implementation details of GitHub Actions workflows (e.g. specific job layout, caching strategy, or check-only formatting implementation) are left to engineering as long as the above requirements and acceptance criteria are met.


