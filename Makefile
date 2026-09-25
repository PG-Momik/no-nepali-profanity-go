# no-nepali-profanity Makefile

.PHONY: all test test-race build build-example run-example fmt vet lint clean deps vuln

# Run all tests
test:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race -v ./...

# Build the package
build:
	go build ./...

# Build the example
build-example:
	go build -o bin/no-nepali-profanity ./example

# Run the example
run-example: build-example
	./bin/no-nepali-profanity

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean -cache

# Download dependencies
deps:
	go mod download
	go mod tidy

# Check for vulnerabilities
vuln:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Default target
all: fmt vet test build