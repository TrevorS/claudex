.PHONY: help build run test clean check format vet dev dev-debug vhs-test vhs-test-verbose vhs-record vhs-update

.DEFAULT_GOAL := help

help:
	@echo "claudex - Claude Code Conversation History Viewer"
	@echo ""
	@echo "Available targets:"
	@echo "  make build            - Build the claudex binary"
	@echo "  make run              - Build and run claudex"
	@echo "  make dev              - Run with auto-rebuild on file changes (using Air)"
	@echo "  make dev-debug        - Run with auto-rebuild and debug logging enabled"
	@echo "  make test             - Run all tests"
	@echo "  make check            - Format, vet, and test (recommended before commit)"
	@echo "  make format           - Format code with gofmt"
	@echo "  make vet              - Run go vet"
	@echo "  make vhs-test         - Run VHS tape tests for UI regression testing"
	@echo "  make vhs-test-verbose - Run VHS tests with verbose output"
	@echo "  make vhs-record       - Record all VHS tapes (GIFs + screenshots)"
	@echo "  make vhs-update       - Update expected screenshots from current output"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make help             - Show this help message"

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
	@./scripts/vhs-test.sh

vhs-test-verbose: build
	@echo "📹 Running VHS tape tests (verbose)..."
	@./scripts/vhs-test.sh -v

vhs-record: build
	@echo "📹 Recording all VHS tapes..."
	@for tape in testdata/vhs/*.tape; do \
		echo "Recording $$tape..."; \
		vhs "$$tape"; \
	done
	@echo "✓ VHS recording complete - check testdata/vhs/ for output"

vhs-update: build
	@echo "📹 Updating expected screenshots..."
	@./scripts/vhs-test.sh -u
	@echo "✓ Expected screenshots updated"
