## Infrastructure Decisions

### Piston (Code Execution Engine)

- **Local development**: Piston runs as a Docker container (`innogen-piston`) alongside Postgres and Redis.
- **Access URL**: `http://localhost:2000` from host machine, `http://piston:2000` from other Docker containers on the same network.
- **Image**: `ghcr.io/engineer-man/piston:latest`
- **Privileged mode**: Required for Piston's internal container-based isolation.
- **Volume**: `pistondata:/piston` persists installed language packages across restarts.
- **Decision**: Piston is NOT built into the `innogen-backend` monorepo. It is a separate external service, consumed via REST API.

### X-User-ID Temporary Authentication

- **Decision**: For STEP 5 (submission_service) and development convenience, a temporary `X-User-ID` HTTP header is used for authentication.
- **Implementation**: `shared/middleware/auth.go` provides `XUserID()` middleware that reads the integer user ID from the header and injects it into the request context using the same `UserIDKey` used by JWT `Auth()`.
- **Production**: Replace with JWT Bearer token authentication via `middleware.Auth(jwtSecret)` before deploying.
- **Routes using X-User-ID**: POST /submit, GET /submissions/{id}, GET /me/submissions, GET /me/submissions/{problemId}/latest

### Redis Queue for Async Judging

- **Decision**: Use Redis list-based queue (LPUSH/BRPOP) for submission job distribution.
- **Queue name**: `submission_jobs`
- **Payload**: `{"submissionId":"<uuid>"}` (JSON)
- **Delivery semantics**: At-most-once. Jobs removed from Redis on dequeue before processing.
- **Crash recovery**: MVP uses a placeholder reconciliation sweep. Production should implement periodic re-enqueue of stale Pending/Running submissions.
- **Library**: `github.com/redis/go-redis/v9`

### Shared Packages Extraction

- **Decision**: Extracted `piston` client, `judge` logic, status `constants`, and `languageutil` into `shared/` to avoid code duplication between run_service and submission_service worker.
- **Affected packages**:
  - `shared/piston/` — Piston API v2 HTTP client (exported types: Request, Response, Stage, Client)
  - `shared/judge/` — Output evaluation with TLE detection via timeLimitMs
  - `shared/constants/` — Status string constants (Accepted, WrongAnswer, etc.)
  - `shared/languageutil/` — File name resolution from language config
- **Run service compatibility**: run_service re-exports status constants from shared/constants for backward compatibility.

### Mock Commit Flow (STEP 7)

- **Decision**: Use mock commit SHA (40-char random hex) generated via `crypto/rand` for MVP. No real GitHub API calls are made.
- **Integration**: `submission_service` worker imports `repo_service` directly (same Go module, monorepo pattern — no HTTP overhead).
- **Repository model**: One repository per `(user_id, subject_id)` pair, created/updated via `INSERT ... ON CONFLICT`.
- **Path format**: `<subjectSlug>/Session-<NN>/Lesson-<NN>/Problem-<NN>-<problemSlug>/<fileName>` with zero-padded order indexes.
- **Repo name format**: `<subjectSlug>-RinnoGen`
- **Transaction safety**: Commit DML (upsert repo, update submission, insert commit) wrapped in a single `pgx` transaction.
- **Idempotency**: `CommitSubmission` checks `commit_sha` before proceeding — safe to call multiple times.
- **Packages**: `repo_service/repository/` and `repo_service/service/` are outside `internal/` to allow cross-service imports by `submission_service`. Pathbuilder remains in `internal/` as it's only used by repo_service.

### Time Limit Detection

- **Decision**: Piston `run_timeout` field is set from `problems.time_limit_ms`. If a process is killed by signal AND a time limit was configured, it's classified as Time Limit Exceeded instead of Runtime Error. Without a configured limit, signals are treated as Runtime Error (backward compatible with run_service).
