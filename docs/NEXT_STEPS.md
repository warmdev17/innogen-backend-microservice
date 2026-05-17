## All 11 steps complete

### Run the system
```bash
make compose-up && make seed-dev
make run-gateway & make run-auth & make run-runner & make run-submission & make run-submission-worker & make run-repo &
make e2e
```

### Admin API usage
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}' | jq -r .accessToken)

# Create a language
curl -X POST http://localhost:8080/admin/languages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Python 3","pistonAlias":"python","pistonVersion":"3.10.0","fileExtension":".py","defaultFileName":"solution.py"}'

# List subjects
curl http://localhost:8080/admin/subjects -H "Authorization: Bearer $TOKEN"

# Non-admin access → 403
```

### Production readiness
See docs/DECISIONS.md for full backlog.
