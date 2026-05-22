# Deployment Guide

## Docker Architecture

| Service | Port | Depends On |
|---------|------|-----------|
| api_gateway | 8080 | auth, run, submission, repo |
| auth_service | 8081 | postgres |
| run_service | 8082 | postgres, piston |
| submission_service | 8083 | postgres, redis |
| submission_worker | — | postgres, redis, piston, repo |
| repo_service | 8084 | postgres |

## Quick Start (Docker)

```bash
# Build and start all services
make compose-full-up

# Apply schema + seed
make seed-dev

# Check health
curl http://localhost:8080/health

# View logs
make compose-full-logs

# Stop everything
make compose-full-down

# Reset volumes
make docker-reset
```

## Service URLs (Docker)

Internal: `http://<service_name>:<port>`
Host: `http://localhost:<port>`

## GitHub Secrets

Mount private key:
```yaml
volumes:
  - ./secrets:/app/secrets:ro
```
Set: `GITHUB_PRIVATE_KEY_PATH=/app/secrets/github-key.pem`
