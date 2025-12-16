.PHONY: run build test test-coverage test-coverage-html

# Run the API
run:
	go run main.go

# Build the binary
build:
	go build -o main.exe .

# Run all tests
test:
	go test ./...

test-coverage:
	go test ./internal/... ./utils/identity/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-coverage-html:
	go test ./internal/... ./utils/identity/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-full:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
