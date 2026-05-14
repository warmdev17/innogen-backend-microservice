# AGENTS.md — innogen-backend

## Project overview

RinnoGen backend — a LeetCode-like platform for curriculum-based programming practice. Go monorepo with microservice-style separation, one root `go.mod`.

## Architecture

- **One root `go.mod`** — module `innogen-backend`, Go 1.26.2
- **`shared/`** — packages imported by all services: `config`, `logger`, `database`, `response`, `errors`, `models`, `middleware`
- **Services** (each has `cmd/main.go` as entrypoint):
  - `api_gateway/` — API Gateway (port 8080)
  - `auth_service/` — Authentication (port 8081)
  - `run_service/` — Code execution via Piston (port 8082)
  - `submission_service/` — Submission + judging (API port 8083, worker has no HTTP)
  - `repo_service/` — GitHub repository management (port 8084)
- **`schema.sql`** — Full PostgreSQL schema (PostgreSQL, pgcrypto). This is the source of truth for the data model.

## Key commands

```bash
# Run individual services
make run-gateway       # API Gateway on :8080
make run-auth          # Auth on :8081
make run-runner        # Run service on :8082
make run-submission    # Submission API on :8083
make run-submission-worker  # Submission worker (no HTTP)
make run-repo          # Repo service on :8084

# Development
make tidy              # go mod tidy
make fmt               # gofmt -s -w .
make test              # go test ./...
```

## Dependencies

- **Database**: `github.com/jackc/pgx/v5/pgxpool` (PostgreSQL connection pool)
- **No web framework** — use `net/http` with Go 1.22+ `ServeMux` pattern routing (`"GET /path"`)
- **Logging**: stdlib `log/slog` with JSON handler
- **JSON**: `encoding/json` only (no third-party JSON libraries)

## Conventions

- All log messages and error messages **must be in English**
- Go variables use **camelCase**
- Run `gofmt -s -w .` before committing
- Config is loaded from environment variables via `shared/config.Load()`
- Default config values are in `shared/config/config.go`
- Infrastructure dependencies (Postgres, Redis) are in `docker-compose.yml`

## Database

- PostgreSQL 16
- Database: `innogen`, User: `innogen`, Password: `innogen`
- Schema file: `schema.sql` — includes tables, triggers, indexes, seed data
- Key triggers to be aware of:
  - **anti-spam**: 10-second cooldown between submissions per user/problem (DB-level trigger)
  - **acceptance rate**: auto-updated on submission insert/status change

## External services

- **Piston** — code execution engine, default URL `http://localhost:2000`
- **GitHub App** — for repository management (integration not yet implemented)
