## Implementation Progress

### ✅ Completed
- [x] STEP 1 — Bootstrap shared infrastructure and service skeletons
- [x] STEP 2 — auth_service login and current-user (JWT + bcrypt)
- [x] STEP 3 — curriculum/problem read APIs (8 endpoints in api_gateway)
- [x] STEP 4 — run_service code execution with Piston (POST /run)
- [x] STEP 5 — submission_service create Pending submissions (POST /submit + read APIs)
- [x] STEP 6 — Redis queue and judge worker (async judging)
- [x] STEP 7 — repo_service mock commit (path builder, repo upsert, commit records)
- [x] STEP 8 — GitHub App implementation (real commits, per-user installations)
- [x] STEP 9 — api_gateway proxy/routing (single entrypoint, JWT auth passthrough)
- [x] STEP 10 — Local seed data and E2E MVP validation
- [x] STEP 11 — Admin content management APIs (CRUD for languages, subjects, sessions, lessons, problems, test cases, tags)
- [x] STEP 12 — GitHub App webhook handling (installation, repo events, signature verification)

- [x] STEP 13 — Local developer test UI (React + Vite + TypeScript)

- [x] STEP 14 — GitHub Commit E2E Visibility (commit URLs in submission responses)
- [x] STEP 15 — Test UI GitHub Connect UX improvements
- [x] STEP 16 — GitHub OAuth Account Linking (commit author identity)
- [x] STEP 16 fix — Stabilize GitHub App installation lifecycle

### 🔜 Pending
- [ ] STEP 17 — Commit attribution fix using OAuth identity
