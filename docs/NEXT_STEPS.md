## Next Steps

### Running the system
```bash
make compose-up
psql postgres://innogen:innogen@localhost:5432/innogen < schema.sql

# If upgrading from STEP 7:
psql postgres://innogen:innogen@localhost:5432/innogen < migrations/001_add_github_owner.sql

# Seed GitHub account (example):
psql postgres://innogen:innogen@localhost:5432/innogen -c "
INSERT INTO github_accounts (user_id, installation_id, github_owner, github_owner_type)
VALUES (1, '12345678', 'your-username', 'User');
"

# Start services:
make run-gateway          # :8080
make run-auth             # :8081
make run-runner           # :8082
make run-submission       # :8083
make run-submission-worker
make run-repo             # :8084
```

### Upcoming (STEP 9)
- Implement api_gateway reverse proxy to all services
- Add JWT auth passthrough
- Route: /api/auth/* → auth_service
- Route: /api/subjects/* → api_gateway (curriculum)
- Route: /api/run/* → run_service
- Route: /api/submissions/* → submission_service
- Route: /api/repos/* → repo_service
- Single entry point for frontend at :8080
