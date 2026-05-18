DATABASE_URL ?= postgres://innogen:innogen@localhost:5432/innogen?sslmode=disable

.PHONY: run-gateway run-auth run-runner run-submission run-submission-worker run-repo run-all stop-all tidy fmt test compose-up compose-down logs-piston piston-install seed-dev e2e

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

piston-install:
	@echo "Installing Node.js 18.15.0 runtime in Piston..."
	curl -s -X POST http://localhost:2000/api/v2/packages \
		-H "Content-Type: application/json" \
		-d '{"language":"node","version":"18.15.0"}' || true
	@echo "Waiting for runtime installation..."
	@sleep 5
	@echo "Verifying runtime..."
	@curl -s http://localhost:2000/api/v2/runtimes | python3 -c "import sys,json; d=json.load(sys.stdin); [print(f'  {r[\"language\"]} {r[\"version\"]}') for r in d if 'javascript' in r.get('language','').lower()]" || echo "  (verification skipped - python3 not available)"
	@echo "Done."

seed-dev: piston-install
	psql "$(DATABASE_URL)" -f schema.sql
	psql "$(DATABASE_URL)" -f seeds/dev_seed.sql

e2e:
	bash scripts/e2e_mvp.sh

run-all:
	@echo "Starting all services in background..."
	@echo "  api_gateway      :8080"
	@echo "  auth_service     :8081"
	@echo "  run_service      :8082"
	@echo "  submission_service :8083"
	@echo "  repo_service     :8084"
	@echo "  submission_worker (background)"
	@echo ""
	@echo "Press Ctrl+C to stop all services."
	go run ./api_gateway/cmd/main.go &
	go run ./auth_service/cmd/main.go &
	go run ./run_service/cmd/main.go &
	go run ./submission_service/cmd/api/main.go &
	go run ./submission_service/cmd/worker/main.go &
	go run ./repo_service/cmd/main.go &
	wait

stop-all:
	@echo "Stopping all services..."
	@pkill -f "api_gateway/cmd/main.go" 2>/dev/null || true
	@pkill -f "auth_service/cmd/main.go" 2>/dev/null || true
	@pkill -f "run_service/cmd/main.go" 2>/dev/null || true
	@pkill -f "submission_service/cmd/api/main.go" 2>/dev/null || true
	@pkill -f "submission_service/cmd/worker/main.go" 2>/dev/null || true
	@pkill -f "repo_service/cmd/main.go" 2>/dev/null || true
	@echo "All services stopped."
