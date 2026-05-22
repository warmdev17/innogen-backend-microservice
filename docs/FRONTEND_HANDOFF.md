# Frontend Integration Guide

## Base URL
- Development: `http://localhost:8080`

## API Documentation
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI Spec: `http://localhost:8080/docs/openapi.yaml`

## Auth Flow
1. POST /auth/register or POST /auth/login → get accessToken
2. Authorization: Bearer <accessToken> on protected routes

## Response Envelope
All responses use: `{ status, code, message, data }` for success, `{ status, code, message, error, details }` for errors.

## Business Statuses
- submission.status: Pending, Running, Accepted, Wrong Answer, etc.
- run.status: Accepted, Wrong Answer, Compilation Error, etc.
- githubConnection.status: active, suspended, deleted, disconnected

## Common Routes
See Swagger UI for full API reference.
