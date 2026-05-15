## Next Steps

### Immediate

1. Run infrastructure: `make compose-up` (starts postgres, redis, piston)
2. Apply schema: `psql postgres://innogen:innogen@localhost:5432/innogen < schema.sql`
3. Seed test data into the database
4. Start run_service: `make run-runner`
5. Test: `curl -X POST http://localhost:8082/run -H "Content-Type: application/json" -d '{"problemId":1,"languageId":1,"code":"..."}'`

### Piston Notes

- Local Piston API: `http://localhost:2000`
- Docker-network URL for containerized Go services: `http://piston:2000`
- Verify Piston is running: `curl http://localhost:2000/api/v2/runtimes`
- Piston may take 1-5 minutes to download language packages on first start

### Upcoming

- Implement submission_service POST /submissions
- Redis queue for judge worker
- repo_service mock commit
