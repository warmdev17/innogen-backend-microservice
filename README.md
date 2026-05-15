# RinnoGen Backend

A LeetCode-like curriculum-based programming practice platform.

## Architecture

Go monorepo with microservice-style separation:

| Service | Port | Description |
|---------|------|-------------|
| `api_gateway` | 8080 | Single public entrypoint, JWT auth, proxy |
| `auth_service` | 8081 | Login, user info |
| `run_service` | 8082 | Code execution via Piston |
| `submission_service` | 8083 | Submission CRUD, async judging |
| `repo_service` | 8084 | GitHub repository commit |

Infrastructure: PostgreSQL 16, Redis 7, Piston (code execution engine)

## Quick Start

```bash
# 1. Start infrastructure
make compose-up

# 2. Apply schema and seed data
make seed-dev

# 3. Start all services (each in a separate terminal)
make run-gateway &
make run-auth &
make run-runner &
make run-submission &
make run-submission-worker &
make run-repo &

# 4. Run E2E validation
make e2e
```

## Development

```bash
make fmt          # Format all Go files
make tidy         # Clean up dependencies
make test         # Run all tests
```

## API Endpoints (via API Gateway :8080)

| Method | Path | Auth |
|--------|------|------|
| POST | /auth/login | No |
| GET | /auth/me | Bearer |
| GET | /subjects | No |
| POST | /run | Bearer |
| POST | /submit | Bearer |
| GET | /submissions/{id} | Bearer |
| GET | /me/submissions | Bearer |

## Environment

Copy `.env.example` to `.env` and adjust as needed:

```bash
cp .env.example .env
```
