# Common Commands

## Implement a feature
1. Create or update files in `apps/frontend` or `apps/backend`.
2. Add/extend tests.
3. Run `make typecheck && make lint && make test`.

## Add an endpoint (backend)
- Create handler in `apps/backend`.
- Wire route with chi.
- Add tests: `go test ./...`.


