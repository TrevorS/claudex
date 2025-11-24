.PHONY: help build run test clean check format vet dev dev-debug vhs-test vhs-record

.DEFAULT_GOAL := help

help:
	@echo "claudex - Claude Code Conversation History Viewer"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build the claudex binary"
	@echo "  make run          - Build and run claudex"
	@echo "  make dev          - Run with auto-rebuild on file changes (using Air)"
	@echo "  make dev-debug    - Run with auto-rebuild and debug logging enabled"
	@echo "  make test         - Run all tests"
	@echo "  make check        - Format, vet, and test (recommended before commit)"
	@echo "  make format       - Format code with gofmt"
	@echo "  make vet          - Run go vet"
	@echo "  make vhs-test     - Run VHS tape tests for UI regression testing"
	@echo "  make vhs-record   - Record VHS demos for documentation"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make help         - Show this help message"

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

dev:
	@echo "🚀 Starting development mode with auto-reload..."
	air

dev-debug:
	@echo "🚀 Starting development mode with auto-reload and debug logging..."
	DEBUG=1 air

vhs-test: build
	@echo "📹 Running VHS tape tests..."
	vhs testdata/vhs/startup.tape
	vhs testdata/vhs/search.tape
	@echo "✓ VHS tests complete - check testdata/vhs/ for screenshots"

vhs-record: build
	@echo "📹 Recording VHS demos for documentation..."
	for tape in testdata/vhs/*.tape; do \
		echo "Recording $$tape..."; \
		vhs "$$tape"; \
	done
	@echo "✓ VHS recording complete"
