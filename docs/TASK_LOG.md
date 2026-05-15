## Task Log

### 2026-05-15 — Piston Infrastructure

- Added `piston` service to `docker-compose.yml` (ghcr.io/engineer-man/piston:latest, port 2000)
- Added `compose-up`, `compose-down`, `logs-piston` targets to `Makefile`
- Documented Piston access URLs in `docs/DECISIONS.md` and `docs/NEXT_STEPS.md`
- Validation: `docker compose up -d postgres redis piston` starts all three infrastructure services
