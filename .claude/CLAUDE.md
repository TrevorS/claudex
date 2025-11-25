# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Claudex** is a terminal UI application for browsing, searching, and analyzing Claude Code conversation history. It's built with Go using the Bubble Tea framework.

See `docs/todo.md` for the complete development roadmap and current implementation status.

## Architecture

### High-Level Design

The application follows a clean architecture with clear separation of concerns:

1. **Domain Layer** (`internal/domain/`) - Pure business logic, completely decoupled from frameworks
   - `Conversation`, `Message`, `TokenCount` - Core entities
   - `Statistics` - Aggregation calculations
   - `AppError` and `ErrorType` - Typed error hierarchy for domain errors
   - All domain types are immutable after creation

2. **Data Access Layer** (`internal/repository/`) - Loads and manages conversations
   - `Repository` interface - Abstract data access contract
   - `JSONLRepository` implementation - Parses `~/.claude/projects/` directory structure
   - Will include `CachedRepository` wrapper for performance
   - Handles hyphenated path decoding (e.g., `foo-bar` → `foo/bar`)

3. **Search Layer** (`internal/search/`) - Query parsing and filtering
   - `Lexer` - Tokenizes query strings
   - `Query` (AST) - Parses tokens into abstract syntax tree
   - `Filter` interface - Composable filtering logic
   - `Engine` - Executes searches with relevance ranking

4. **UI Layer** (`internal/ui/`) - Bubble Tea models and views
   - `Model` - Root Bubble Tea model with global state management
   - `update.go` - Centralized message routing and state transitions
   - `view.go` - Three-pane layout composition with Lipgloss
   - `views/` - Individual view implementations (list, detail, search, stats, palette)
   - `components/` - Reusable UI components (sidebar, header, status bar)

5. **Application Bootstrap** (`internal/app/`) - Initialization and configuration
   - `App` struct - Coordinates repository, UI, and Bubble Tea
   - `Config` - Settings, paths, and runtime configuration

### Data Flow

```
Main → App (bootstrap) → Repository (load conversations) → UI Model
                                                          ↓
                                                     Views render
                                                          ↓
                                              User input → Update message
                                                          ↓
                                                     State changes
```

### Key Design Decisions

1. **Immutable Domain Entities** - Domain types (Conversation, Message) don't have setters; they're created once and not modified
2. **Typed Errors** - `AppError` with `ErrorType` enum provides context-aware error handling
3. **Repository Pattern** - Abstract repository interface allows easy swapping implementations (JSONL, cached, paginated)
4. **Lazy Loading** - Full message content only loaded when viewing a conversation, not at startup
5. **Three-Pane Layout** - Sidebar (projects), center (conversations/messages), right (metadata) with responsive handling

### Known Gotchas

1. **Lipgloss Width() after Border()** - When composing multi-pane layouts with lipgloss, apply `Width()` and `Height()` constraints BEFORE adding borders/padding, not after. If you render content with borders first, then try to constrain with `Width()`, the already-bordered content may exceed the constraint and break layout calculations. See `view.go` `renderXxxWithSize()` functions for the correct pattern.

2. **bubble-table breaks lipgloss layouts** - The `bubble-table` library's ANSI escape sequences confuse lipgloss width calculations. We replaced it with simple string rendering in `ListView`. If you need tables, use plain text formatting instead.

3. **Emoji width calculation** - Emojis can cause width miscalculations in terminal layouts. Prefer ASCII characters for UI elements (e.g., `>` instead of 📁).

## Development Workflow

### Common Commands

```bash
# Build the application
make build

# Build and run
make run

# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run a single test file
go test -v ./internal/domain -run TestConversation

# Run a specific test
go test -v ./internal/domain -run TestConversationAddMessage

# Format code
make format

# Lint code
make vet

# Run all checks (format, vet, test)
make check

# Clean build artifacts
make clean
```

### UI Development (Rapid Iteration)

```bash
# Auto-rebuild on file changes (sub-second feedback loop)
make dev

# Same as dev, but also logs all Bubble Tea messages to debug.log
make dev-debug

# Run VHS tapes and capture UI screenshots
make vhs-test

# Record VHS demos for documentation
make vhs-record
```

### Project Structure Guidelines

- **`cmd/claudex/`** - Entry point only, minimal logic
- **`internal/domain/`** - Domain entities and business logic (no external dependencies beyond stdlib)
- **`internal/repository/`** - Data access implementations
- **`internal/search/`** - Search logic
- **`internal/ui/`** - All Bubble Tea UI code
- **`internal/app/`** - Application bootstrap
- **`internal/cache/`** - Caching utilities
- **`internal/render/`** - Message rendering (markdown, syntax highlighting)
- **`pkg/jsonl/`** - Public JSONL parser library (reusable beyond this project)
- **`testdata/`** - Test fixtures and sample data

### Code Style

- Start every file with `// ABOUTME: ` comment explaining what the file does
- Use `gofmt -s` (simple format)
- Follow Go standard library naming conventions
- Keep functions small and focused
- No circular dependencies between packages
- Domain layer must not import from other layers (pure business logic)

## Claude Code Skills

Three project skills in `.claude/skills/` provide guidance for common development tasks:

- **dev-mode** - Live reload with Air (`make dev`) for rapid UI iteration
- **debug-ui** - Bubble Tea debug logging (`make dev-debug`) for message flow inspection
- **capture-ui** - VHS terminal recordings for documentation and demos

These skills are automatically available when working with Claude Code and can be invoked by asking about UI development, debugging, or creating demos.

## Testing Strategy

- **Test Coverage Target:** 85%+ for domain/search, 75%+ for other layers
- **Test Location:** `*_test.go` files in the same package as the code
- **Test Data:** Use `testdata/` directory for fixtures
- **Assertions:** Use `github.com/stretchr/testify` (already in go.mod)

### Writing Tests

```go
// Use the testify assertion functions
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    result := DoSomething()
    assert.Equal(t, expected, result)
    assert.NoError(t, err)
}
```

## Dependencies

See `go.mod` for current dependencies. See `docs/execution-plan.md` for planned framework and library additions by phase.

## Git Workflow

- Use feature branches for work (created with `/feature-branch create`)
- All commits should reference the todo list items being completed
- Use `make check` before committing to ensure code quality
- Commit messages should be descriptive and reference phase/task

## Key Documentation Files

- **`docs/execution-plan.md`** - Technical architecture and specifications for all planned phases
- **`docs/ux-specs.md`** - User interaction patterns and UI design details
- **`docs/todo.md`** - Complete development roadmap with current implementation status and checkable tasks (the source of truth for what's done/in progress)

## Important Notes

- The domain layer is complete and tested - treat it as stable
- Path handling: Claude project paths use hyphens in filesystem (e.g., `foo-bar`) but represent nested paths (e.g., `foo/bar`)
- Repository interface should hide complexity; consumers just call `List()`, `GetByID()`, `Search()`
- All UI should be keyboard-driven; mouse support is secondary
- Bubble Tea models should be immutable - create new models in update handlers rather than mutating state
