## Task Log

### 2026-05-15 — Piston Infrastructure

- Added `piston` service to `docker-compose.yml` (ghcr.io/engineer-man/piston:latest, port 2000)
- Added `compose-up`, `compose-down`, `logs-piston` targets to `Makefile`
- Documented Piston access URLs in `docs/DECISIONS.md` and `docs/NEXT_STEPS.md`
- Validation: `docker compose up -d postgres redis piston` starts all three infrastructure services

### 2026-05-15 — Submission Service API (STEP 5)

- Implemented POST /submit (create Pending submission)
- Implemented GET /submissions/{id} (get submission by UUID with auth)
- Implemented GET /me/submissions (list user's submissions, no code in response)
- Implemented GET /me/submissions/{problemId}/latest (latest for user+problem)
- Added `Submission` model to shared/models
- Added `XUserID()` middleware for temporary header-based auth
- Added anti-spam detection (429 response on 10s cooldown trigger)
- Added UUID format validation on path parameters
- Added pagination limit (50) on list endpoint
- Files: dto.go, repository.go, service.go, handler.go, route.go, cmd/api/main.go
- Validated: go build, go vet, go test all pass
- Commit: pending

### 2026-05-15 — Redis Queue and Judge Worker (STEP 6)

- Extracted shared packages: `shared/piston/`, `shared/judge/`, `shared/constants/`, `shared/languageutil/`
- Migrated run_service to use shared packages
- Created Redis-backed job queue (`submission_service/internal/queue/`) using `go-redis/v9`
- Implemented async worker (`submission_service/internal/worker/`) with full judging flow:
  - Dequeue from Redis → load submission → mark Running → execute all test cases via Piston → judge results → update submission
  - Status priority: CompilationError > InternalError > TLE > RuntimeError > WrongAnswer > Accepted
  - Per-test TLE detection via problem's time_limit_ms + Piston run_timeout
- Updated POST /submit to enqueue submission for async judging
- Added `UpdateSubmissionStatus` and `UpdateSubmissionResult` repository methods
- Updated handler to use proper sentinel error for spam cooldown
- Deleted old run_service/internal/piston/ and run_service/internal/judge/ (moved to shared/)
- Validated: go build, go vet, go test (18/18 pass) all clean
- Commit: pending
