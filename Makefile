DATABASE_URL ?= postgres://innogen:innogen@localhost:5432/innogen?sslmode=disable

.PHONY: run-gateway run-auth run-runner run-submission run-submission-worker run-repo run-all stop-all tidy fmt test compose-up compose-down logs-piston piston-install seed-dev e2e test-ui-install test-ui test-ui-build

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
	@echo "Installing Node.js 18.15.0 runtime in Piston (timeout 120s)..."
	@curl -s --max-time 120 -X POST http://localhost:2000/api/v2/packages \
		-H "Content-Type: application/json" \
		-d '{"language":"node","version":"18.15.0"}' || echo "  (install skipped or failed - may already be installed)"
	@echo "Done."

seed-dev:
	psql "$(DATABASE_URL)" -f schema.sql
	psql "$(DATABASE_URL)" -f seeds/dev_seed.sql

e2e: piston-install
	bash scripts/e2e_mvp.sh

run-all:
	@echo "Starting all services in background..."
	@echo "  api_gateway         :8080"
	@echo "  auth_service        :8081"
	@echo "  run_service         :8082"
	@echo "  submission_service  :8083"
	@echo "  repo_service        :8084"
	@echo "  submission_worker   (background)"
	@echo ""
	@echo "Press Ctrl+C to stop all services."
	@bash -c '\
		trap "echo; echo Stopping all services...; kill 0; exit 0" INT TERM EXIT; \
		go run ./api_gateway/cmd/main.go & \
		go run ./auth_service/cmd/main.go & \
		go run ./run_service/cmd/main.go & \
		go run ./submission_service/cmd/api/main.go & \
		go run ./submission_service/cmd/worker/main.go & \
		go run ./repo_service/cmd/main.go & \
		wait \
	'

stop-all:
	@echo "Stopping all services by ports..."
	@for port in 8080 8081 8082 8083 8084; do \
		pids=$$(lsof -ti tcp:$$port 2>/dev/null); \
		if [ -n "$$pids" ]; then \
			echo "  killing port $$port: $$pids"; \
			kill $$pids 2>/dev/null || true; \
		else \
			echo "  port $$port: no process"; \
		fi; \
	done
	@pkill -f "submission_service/cmd/worker" 2>/dev/null || true
	@pkill -f "go run ./submission_service/cmd/worker/main.go" 2>/dev/null || true
	@echo "All services stopped."

test-ui-install:
	cd tools/test_ui && npm install

test-ui:
	cd tools/test_ui && npm run dev

test-ui-build:
	cd tools/test_ui && npm run build
