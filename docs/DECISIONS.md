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
