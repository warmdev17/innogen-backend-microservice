## Next Steps

### Running the system
```bash
make compose-up
psql postgres://innogen:innogen@localhost:5432/innogen < schema.sql

# Terminals:
make run-gateway          # :8080
make run-auth             # :8081
make run-runner           # :8082
make run-submission       # :8083
make run-submission-worker
make run-repo             # :8084
```

### Upcoming (STEP 8)
- Implement real GitHub App authentication (installation token, JWT)
- Create real GitHub repositories via GitHub API
- Push real commits with accepted solution code
- Replace mock commit SHA with real GitHub commit SHA
- Implement OAuth/GitHub login flow
