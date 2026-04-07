.PHONY: build clean test install

BINARY_NAME=tools-decision
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
LDFLAGS=-ldflags "-X github.com/yudgnahk/tools-decision/cmd/tools-decision/commands.version=$(VERSION) -X github.com/yudgnahk/tools-decision/cmd/tools-decision/commands.commit=$(COMMIT)"

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/tools-decision

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/tools-decision

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Build for all platforms
build-all: clean
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/tools-decision
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/tools-decision
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/tools-decision
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/tools-decision
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/tools-decision

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Update dependencies
deps:
	go mod tidy
	go mod download

# Run the tool
run:
	go run ./cmd/tools-decision

# Development: build and run
dev: build
	./$(BINARY_NAME)
