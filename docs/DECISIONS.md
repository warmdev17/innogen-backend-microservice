## Infrastructure Decisions

### Piston (Code Execution Engine)

- **Local development**: Piston runs as a Docker container (`innogen-piston`) alongside Postgres and Redis.
- **Access URL**: `http://localhost:2000` from host machine, `http://piston:2000` from other Docker containers on the same network.
- **Image**: `ghcr.io/engineer-man/piston:latest`
- **Privileged mode**: Required for Piston's internal container-based isolation.
- **Volume**: `pistondata:/piston` persists installed language packages across restarts.
- **Decision**: Piston is NOT built into the `innogen-backend` monorepo. It is a separate external service, consumed via REST API.
