.PHONY: build clean test run install

# Build the binary
build:
	go build -o remote-claude-mcp

# Build for different architectures
build-all:
	GOOS=linux GOARCH=amd64 go build -o remote-claude-mcp-linux-amd64
	GOOS=linux GOARCH=arm64 go build -o remote-claude-mcp-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -o remote-claude-mcp-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o remote-claude-mcp-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o remote-claude-mcp-windows-amd64.exe

# Clean build artifacts
clean:
	rm -f remote-claude-mcp*
	go clean

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run the server with example config
run: build
	@echo "Create a config.json file based on config.example.json first"
	@echo "Then run: ./remote-claude-mcp config.json"

# Install to system PATH
install: build
	sudo cp remote-claude-mcp /usr/local/bin/

# Development setup
dev-setup:
	go mod tidy
	go install golang.org/x/tools/gopls@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Lint and format
lint:
	gofmt -s -w .
	go vet ./...
	staticcheck ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  build-all   - Build for all architectures"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  deps        - Install dependencies"
	@echo "  run         - Build and show run instructions"
	@echo "  install     - Install to system PATH"
	@echo "  dev-setup   - Setup development environment"
	@echo "  lint        - Run linting and formatting"
	@echo "  help        - Show this help"