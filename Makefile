.PHONY: run-gateway run-auth run-runner run-submission run-submission-worker run-repo tidy fmt test compose-up compose-down logs-piston

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
