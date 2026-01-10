.PHONY: build run test clean dev-up dev-down

build:
	go build -o bin/orders-api .

run:
	go run . -env=.env.local

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	@grep -v "test_helpers" coverage.out > coverage.filtered.out || true
	go tool cover -html=coverage.filtered.out -o coverage.html
	@rm -f coverage.filtered.out

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

dev-up:
	docker-compose -f ../.devcontainer/docker-compose.yml up -d postgres rabbitmq

dev-down:
	docker-compose -f ../.devcontainer/docker-compose.yml down

dev-logs:
	docker-compose -f ../.devcontainer/docker-compose.yml logs -f

fmt:
	go fmt ./...

tidy:
	go mod tidy