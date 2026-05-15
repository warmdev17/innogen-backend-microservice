## Next Steps

### Running the system

```bash
# Start infrastructure
make compose-up

# Apply schema
psql postgres://innogen:innogen@localhost:5432/innogen < schema.sql

# Start services (each in separate terminal)
make run-gateway          # port 8080
make run-auth             # port 8081
make run-runner           # port 8082
make run-submission       # port 8083
make run-submission-worker  # background process
```

### Test commands

```bash
# Health checks
curl http://localhost:8083/health

# Submit and judge
curl -X POST http://localhost:8083/submit \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"problemId":1,"languageId":1,"code":"console.log(1+2)"}'

# Check result (use the returned UUID)
curl -H "X-User-ID: 1" http://localhost:8083/submissions/<uuid>

# List submissions
curl -H "X-User-ID: 1" http://localhost:8083/me/submissions
```

### Upcoming (STEP 7)

- Implement repo_service mock commit endpoint
- Create GitHub repository path builder
- Store commit SHA on submission
- Implement GET /repos and POST /repos endpoints
- Wire submission_service to call repo_service after Accepted submission
