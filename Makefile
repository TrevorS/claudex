.PHONY: help build run test lint format clean coverage

# Variables
GO := go
GOFLAGS := -v
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the claudex binary
	$(GO) build $(GOFLAGS) -o bin/claudex ./cmd/claudex

run: build ## Build and run claudex
	./bin/claudex

test: ## Run all tests
	$(GO) test $(GOFLAGS) -race -cover ./...

coverage: ## Run tests with coverage reporting
	$(GO) test -race -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

coverage-check: ## Run tests and show coverage percentage
	$(GO) test -race -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -func=$(COVERAGE_FILE) | tail -1
	@rm $(COVERAGE_FILE)

lint: ## Run golangci-lint
	golangci-lint run ./...

format: ## Format code with gofmt and goimports
	gofmt -s -w .
	goimports -w .

fmt: format ## Alias for format

vet: ## Run go vet
	$(GO) vet ./...

clean: ## Clean build artifacts and test files
	$(GO) clean
	rm -f bin/claudex $(COVERAGE_FILE) $(COVERAGE_HTML)

deps: ## Download and verify dependencies
	$(GO) mod download
	$(GO) mod verify

tidy: ## Tidy up go.mod
	$(GO) mod tidy

all: fmt vet lint test build ## Run fmt, vet, lint, test, and build
