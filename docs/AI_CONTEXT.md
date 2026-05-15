# RinnoGen / Innogen Backend AI Context

## Project

This is a Go microservice-style monorepo for a LeetCode-like curriculum learning platform.

The platform organizes programming exercises as:

Subject -> Session -> Lesson -> Problem

## Architecture

One root Go module.

Services:

- api_gateway
- auth_service
- run_service
- submission_service
- repo_service

Shared packages:

- shared/config
- shared/database
- shared/logger
- shared/response
- shared/middleware
- shared/models
- shared/errors

## Rules

- Use net/http only.
- Use PostgreSQL through shared/database.
- Use Redis for queue later.
- Use Piston for code execution.
- Use GitHub App for repo creation and commits later.
- Use schema.sql as the source of truth.
- Do not redesign the architecture.
- Do not rename service folders.
- Do not create a new service unless explicitly requested.
- Logs and error messages must be in English.
- Use camelCase for Go variables where applicable.
- Run gofmt after code changes.
- Run go test ./... before finishing a step.

## Current Implementation Order

1. Bootstrap service skeletons
2. auth_service login and /auth/me
3. curriculum/problem read APIs
4. run_service with Piston
5. submission_service create Pending submission
6. Redis queue and judge worker
7. repo_service mock commit
8. GitHub App implementation
9. api_gateway routing/proxy
