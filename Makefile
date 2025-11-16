build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test-unit:
	go test -v ./tests/unit_tests/...

test-integration: ## Требует запущенный PostgreSQL (docker-compose up -d postgres)
	go test -v ./tests/integration/...

test-all:
	go test -v ./tests/unit_tests/... ./tests/integration/...

test-coverage:
	go test -coverprofile=coverage.out -coverpkg=./internal/... ./tests/unit_tests/... ./tests/integration/...
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет сохранен в coverage.html"

loadtest-burst: ## Требует запущенный API (http://localhost:8080)
	go test -v ./tests/stress_tests/... -run TestBurstLoadTest

loadtest-rampup: ## Требует запущенный API (http://localhost:8080)
	go test -v ./tests/stress_tests/... -run TestRampUpReassignPR

loadtest-all: ## Требует запущенный API (http://localhost:8080)
	go test -v ./tests/stress_tests/...

fmt:
	go fmt ./...

lint:
	golangci-lint run

generate-mocks:
	mockery

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
