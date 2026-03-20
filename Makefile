.PHONY: build install uninstall clean test lint build-all run completion deps tidy

# Binary name
BINARY_NAME=certimate
INSTALL_DIR=$(HOME)/.local/bin

VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/certimate-go/cli/cmd.Version=$(VERSION) \
	-X github.com/certimate-go/cli/cmd.Commit=$(COMMIT) \
	-X github.com/certimate-go/cli/cmd.BuildDate=$(BUILD_DATE)"

# Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Install to ~/.local/bin
install: build
	@mkdir -p $(INSTALL_DIR)
	@cp bin/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME) to $(INSTALL_DIR)/"
	@echo "  Make sure $(INSTALL_DIR) is in your PATH"

# Uninstall from ~/.local/bin
uninstall:
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled $(BINARY_NAME) from $(INSTALL_DIR)/"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf dist/

# Run tests
test:
	go test -v -race ./...

# Run linter
lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	fi

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Run with development settings
run:
	go run . $(ARGS)

# Generate shell completion
completion:
	@echo "# Bash completion:"
	@go run . completion bash
	@echo ""
	@echo "# Zsh completion:"
	@go run . completion zsh
	@echo ""
	@echo "# Fish completion:"
	@go run . completion fish

# Development dependencies
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser/v2@latest

# Tidy dependencies
tidy:
	go mod tidy
