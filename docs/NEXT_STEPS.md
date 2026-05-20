## All 12 steps complete

### Quick validation
```bash
make compose-up && make seed-dev
make run-gateway & make run-auth & make run-runner & make run-submission & make run-submission-worker & make run-repo &
make e2e
```

### Manual webhook test
```bash
SECRET="change-me"
BODY='{"action":"created","installation":{"id":123,"account":{"login":"test","type":"User"}}}'
SIG=$(echo -n "$BODY" | openssl dgst -sha256 -hmac "$SECRET" | sed 's/^.* //')
curl -X POST http://localhost:8084/webhooks/github \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: installation" \
  -H "X-Hub-Signature-256: sha256=$SIG" \
  -d "$BODY"
```

### STEP 17 — Commit attribution fix using OAuth identity
- Use OAuth-linked GitHub identity (username + noreply email) as commit author
- Fall back to installation owner if OAuth not connected
- Ensure commits show the correct author avatar/name on GitHub

### Production readiness
- CORS middleware at gateway
- Request ID middleware
- Panic recovery middleware
- Rate limiting middleware
- CI/CD pipeline
- Dockerfiles for services
- Frontend

### Test UI
```bash
make test-ui-install && make test-ui   # http://localhost:5173
```

### GitHub E2E Test
```bash
make e2e-github   # Requires services + seed data running
```
