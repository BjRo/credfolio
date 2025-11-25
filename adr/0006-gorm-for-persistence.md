# 0006: GORM for Persistence and Migrations

- Status: accepted
- Date: 2025-11-24

## Context

We need a robust persistence layer for the backend. The requirements specify using PostgreSQL as the database. To speed up development and handle database interactions and schema migrations effectively, we need an ORM solution. The project stakeholders have explicitly requested the use of GORM.

## Decision

We will use **GORM** (`gorm.io/gorm`) as the ORM library for the Go backend, coupled with the `pgx` driver (`gorm.io/driver/postgres`).

We will utilize GORM's `AutoMigrate` feature for schema migrations during the initial development phase to support rapid iteration.

## Consequences

- **Pros**:
  - Rapid development with struct-based schema definitions.
  - Built-in support for associations, hooks, and transactions.
  - `AutoMigrate` simplifies schema evolution in early stages.
  - Large community and ecosystem.

- **Cons**:
  - `AutoMigrate` may not handle complex schema changes (e.g., column renaming or data transformations) safely in production without manual intervention or a shift to versioned migrations later.
  - Runtime reflection overhead (though negligible for this use case).
  - Risk of tight coupling between domain entities and database tags if not carefully managed (e.g., via DTOs or Repository pattern).

