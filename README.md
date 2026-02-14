# upl

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/a75e339ad8a045949fd2d032103f71cd)](https://app.codacy.com/gh/zqz/upl/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codebeat badge](https://codebeat.co/badges/d0afb6f1-2490-4ec1-b4a4-213f88800fe4)](https://codebeat.co/projects/github-com-zqz-upl-master)
[![Maintainability](https://api.codeclimate.com/v1/badges/74bcc076dbf4d07c141d/maintainability)](https://codeclimate.com/github/zqz/upl/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/zqz/upl)](https://goreportcard.com/report/github.com/zqz/upl)
![Deploy to Production](https://github.com/zqz/upl/workflows/Deploy%20to%20Production/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=zqz_web&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=zqz_web)


File hosting in Go. Chunked uploads, optional thumbnails, Postgres + sqlc, HTMX UI.

**Requirements:** Go 1.26+, Postgres. Optional: sqlc, goose (or use the Makefile targets that assume them).

```bash
cp .env.example .env   # set DATABASE_URL, SESSION_SECRET at least
make migrate-up
make run
```

Runs on port 3000 by default. See `.env.example` for config.

### Commands

| Command | Description |
|--------|-------------|
| `make run` | Run the server |
| `make dev` | Run with hot reload (air) |
| `make test` | Tests |
| `make migrate-up` / `make migrate-down` | Migrations |
| `make sqlc-generate` | Regenerate sqlc after editing `internal/repository/queries/*.sql` |

API lives under `/api/v1` (e.g. `POST /api/v1/files`, `GET /api/v1/files/{slug}`). Migrations are in `db/migrations/`.
