DATABASE_URL ?= postgres://innogen:innogen@localhost:5432/innogen?sslmode=disable

.PHONY: run-gateway run-auth run-runner run-submission run-submission-worker run-repo tidy fmt test compose-up compose-down logs-piston seed-dev e2e

run-gateway:
	go run ./api_gateway/cmd/main.go

run-auth:
	go run ./auth_service/cmd/main.go

run-runner:
	go run ./run_service/cmd/main.go

run-submission:
	go run ./submission_service/cmd/api/main.go

run-submission-worker:
	go run ./submission_service/cmd/worker/main.go

run-repo:
	go run ./repo_service/cmd/main.go

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

test:
	go test ./...

compose-up:
	docker compose up -d postgres redis piston

compose-down:
	docker compose down

logs-piston:
	docker compose logs -f piston

seed-dev:
	psql "$(DATABASE_URL)" -f schema.sql
	psql "$(DATABASE_URL)" -f seeds/dev_seed.sql

e2e:
	bash scripts/e2e_mvp.sh
