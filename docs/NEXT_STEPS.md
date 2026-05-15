## Next Steps

### Immediate

1. Run infrastructure: `make compose-up`
2. Apply schema: `psql postgres://innogen:innogen@localhost:5432/innogen < schema.sql`
3. Start services:
   - `make run-gateway` (port 8080)
   - `make run-auth` (port 8081)
   - `make run-runner` (port 8082)
   - `make run-submission` (port 8083)

### Test commands

```bash
# Health checks
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Create submission
curl -X POST http://localhost:8083/submit \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"problemId":1,"languageId":1,"code":"console.log(1+2)"}'

# List submissions
curl -H "X-User-ID: 1" http://localhost:8083/me/submissions
```

### Piston Notes

- Local Piston API: `http://localhost:2000`
- Docker-network URL: `http://piston:2000`
- Verify: `curl http://localhost:2000/api/v2/runtimes`

### Upcoming (STEP 6)

- Implement Redis-backed job queue for submission judging
- Create judge worker that consumes from queue
- Worker calls run_service for code execution
- Worker updates submission status after judging
- Implement GET /submissions/{id}/results endpoint
