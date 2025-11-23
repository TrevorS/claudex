.PHONY: help build run test clean check format vet

.DEFAULT_GOAL := help

help:
	@echo "claudex - Claude Code Conversation History Viewer"
	@echo ""
	@echo "Available targets:"
	@echo "  make build    - Build the claudex binary"
	@echo "  make run      - Build and run claudex"
	@echo "  make test     - Run all tests"
	@echo "  make check    - Format, vet, and test (recommended before commit)"
	@echo "  make format   - Format code with gofmt"
	@echo "  make vet      - Run go vet"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make help     - Show this help message"

build:
	go build -o bin/claudex ./cmd/claudex

run: build
	./bin/claudex

test:
	go test -v ./...

check: format vet test
	@echo "✓ All checks passed"

format:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	go clean
	rm -f bin/claudex
