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
