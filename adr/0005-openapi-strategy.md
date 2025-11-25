# 0005: OpenAPI Strategy

- Status: accepted
- Date: 2025-11-24

## Context

As the application grows, maintaining a clear contract between the backend and frontend is crucial. We need a reliable way to document the API, ensure type safety across the stack without manual duplication, and provide a mechanism for testing API endpoints interactively.

## Decision

We have decided to adopt an OpenAPI-first approach with the following constraints:

1.  **OpenAPI for Documentation**: The backend must use OpenAPI to document all exposed API endpoints.
2.  **Interactive Documentation Hosting**: The backend will host an OpenAPI frontend (e.g., Swagger UI) that allows for interactive exploration and testing of the API.
3.  **Security**: Access to the interactive documentation and the ability to execute requests through it is restricted to authenticated users only.
4.  **Code Generation**:
    *   **No Hardcoded Types**: Frontend types for the API must be generated directly from the OpenAPI specifications. Manual definition of API response types in the frontend is prohibited.
    *   **Deprecation Handling**: The generated frontend types and client code must support handling of deprecated backend fields, ensuring backward compatibility and signaling usage of deprecated fields.

## Consequences

*   **Consistency**: The frontend and backend will stay in sync automatically regarding data structures.
*   **Efficiency**: Reduces manual overhead of writing and maintaining TypeScript interfaces for API responses.
*   **Testability**: Developers can test backend endpoints in isolation via the hosted Swagger UI.
*   **Maintenance**: Requires setting up and maintaining the code generation pipeline as part of the build/dev process.

