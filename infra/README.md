# Infrastructure

Docker Compose services:

- Postgres (16-alpine)
- pgAdmin (optional)

Create a `.env` file in this directory (not committed) with:

```
POSTGRES_USER=credfolio
POSTGRES_PASSWORD=credfolio
POSTGRES_DB=credfolio
PGADMIN_DEFAULT_EMAIL=admin@example.com
PGADMIN_DEFAULT_PASSWORD=adminadmin
```

Then run:

```
make db-up
```

Ports:

- Postgres is exposed on host port `55432` (container `5432`)
- pgAdmin is exposed on host port `8081` (container `80`)


