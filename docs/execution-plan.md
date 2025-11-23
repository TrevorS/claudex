# Claudex Technical Execution Plan

**Date:** 2025-11-23
**Project:** Claudex - Claude Code Conversation History Viewer
**Framework:** Bubble Tea (Go TUI)
**Status:** Planning Phase - Ready for Implementation

## Executive Summary

Claudex is a sophisticated terminal UI application that enables developers to browse, search, and analyze their Claude Code conversation history. The implementation uses Bubble Tea framework with a three-pane layout, advanced search capabilities, statistics dashboard, and rich message rendering with syntax highlighting.

The execution plan spans 4 implementation phases with clear deliverables, technical specifications, and testing requirements at each stage.

---

## Project Structure

```
claudex/
├── cmd/
│   └── claudex/
│       └── main.go                    # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go                     # Application bootstrap
│   │   └── config.go                  # Configuration management
│   ├── ui/
│   │   ├── model.go                   # Root Bubble Tea model
│   │   ├── update.go                  # Global update logic
│   │   ├── view.go                    # View composition
│   │   ├── views/
│   │   │   ├── list.go               # Conversation list view
│   │   │   ├── detail.go             # Message viewer view
│   │   │   ├── search.go             # Search modal view
│   │   │   ├── stats.go              # Statistics dashboard
│   │   │   └── palette.go            # Command palette
│   │   ├── components/
│   │   │   ├── statusbar.go          # Footer status bar
│   │   │   ├── header.go             # View headers
│   │   │   └── sidebar.go            # Project/filter sidebar
│   │   └── styles/
│   │       ├── colors.go             # Color palette
│   │       └── styles.go             # Reusable Lipgloss styles
│   ├── domain/
│   │   ├── conversation.go           # Core domain entities
│   │   ├── message.go
│   │   ├── statistics.go
│   │   └── errors.go                 # Error types
│   ├── repository/
│   │   ├── repository.go             # Interface definition
│   │   ├── jsonl.go                  # JSONL implementation
│   │   ├── cached.go                 # Caching wrapper
│   │   └── paginated.go              # Pagination support
│   ├── search/
│   │   ├── engine.go                 # Search engine
│   │   ├── query.go                  # Query parser (AST)
│   │   ├── filter.go                 # Filter composition
│   │   └── lexer.go                  # Query tokenizer
│   ├── cache/
│   │   ├── lru.go                    # LRU cache implementation
│   │   └── intern.go                 # String interning
│   └── render/
│       ├── markdown.go               # Markdown rendering
│       └── syntax.go                 # Syntax highlighting
├── pkg/
│   └── jsonl/
│       └── parser.go                 # Reusable JSONL parser
├── testdata/
│   └── conversations/
│       └── fixtures.jsonl            # Test fixtures
├── go.mod
├── go.sum
├── Makefile                          # Build automation
├── README.md
└── docs/
    ├── ux-specs.md                   # Complete UX specification
    ├── bubble-tea-arch.md            # Architecture research
    ├── bubble-tea-eco.md             # Ecosystem research
    └── execution-plan.md             # This file
```

## Core Dependencies

### Framework & UI
- **github.com/charmbracelet/bubbletea** v1.x - TUI framework
- **github.com/charmbracelet/lipgloss** v1.x - Styling and layout
- **github.com/charmbracelet/bubbles** - Built-in components (list, textarea, spinner, etc.)

### Advanced Components
- **github.com/evertras/bubble-table** - Sophisticated table with sorting/filtering
- **github.com/NimbleMarkets/ntcharts** - Terminal charts and visualizations
- **github.com/erikgeiser/promptkit** - Command palette and prompts
- **github.com/charmbracelet/harmonica** - Spring-based animations

### Utilities
- **github.com/alecthomas/chroma/v2** - Syntax highlighting (100+ languages)
- **github.com/charmbracelet/log** - Structured logging
- **github.com/hashicorp/golang-lru/v2** - Production LRU cache

### Testing
- **github.com/charmbracelet/x/exp/teatest** - Bubble Tea testing library
- **github.com/stretchr/testify** - Assertions and mocking
- **github.com/spf13/afero** - Virtual filesystem for tests

---

## Phase 1: Foundation & Core Structure

**Duration:** 1-2 weeks
**Goal:** Establish project structure, domain models, and basic UI shell

### Deliverables

1. **Project Setup**
   - Initialize Go module with proper import paths
   - Create directory structure following conventions
   - Set up Makefile with build/test/lint targets
   - Configure dependencies in go.mod

2. **Domain Layer** (internal/domain/)
   - `Conversation` struct with ID, Title, Model, CreatedAt, UpdatedAt, Messages
   - `Message` struct with Role, Content, Timestamp, TokenCount, Cost metadata
   - `Statistics` struct for aggregated metrics
   - `AppError` custom error type with type hierarchy

3. **Data Access Layer** (internal/repository/)
   - `ConversationRepository` interface with List(), GetByID(), Search(), GetStatistics()
   - `JSONLRepository` implementation with streaming parser
   - Support for ~/.claude/projects/ directory structure
   - Project discovery and metadata extraction

4. **UI Foundation** (internal/ui/)
   - Root `Model` struct implementing tea.Model
   - ViewState enum: LIST, DETAIL, SEARCH, STATISTICS, COMMAND_PALETTE
   - Basic three-pane layout composition with Lipgloss
   - Window size handling and responsive breakpoints

5. **Navigation & State Machine**
   - View state transitions with history preservation
   - Global keyboard shortcuts (Ctrl+C, Ctrl+P, etc.)
   - Focus management between panes
   - Message routing from child components

### Testing Strategy
- Unit tests for domain entities
- Repository mock for testing UI
- Basic integration test for state transitions
- Use teatest for basic layout verification

### Success Criteria
- Application launches without errors
- Three-pane layout renders correctly
- Window resize handled gracefully
- Keyboard navigation between views works
- 80%+ test coverage for domain layer

---

## Phase 2: Data Display & User Interaction

**Duration:** 2-3 weeks
**Goal:** Implement conversation browsing and viewing functionality

### Deliverables

1. **Conversation List View** (internal/ui/views/list.go)
   - Using evertras/bubble-table with columns:
     - Timestamp (12 chars) - Last Active, relative date formatting
     - Title (flexible) - Intelligent truncation with tooltips
     - Messages (8 chars) - Count with thousands separator
     - Tokens (10 chars) - K/M suffixes with cache indicator
     - Cost (8 chars) - USD with yellow highlight for >$10
   - Sorting: Primary by timestamp desc, secondary by cost desc
   - Multi-column sort support (Shift+Click)
   - Keyboard navigation (arrows, j/k, Page Up/Down, Home/End)
   - Selection with Enter or mouse
   - Per-project sort persistence

2. **Message Viewer** (internal/ui/views/detail.go)
   - Viewport component for scrollable message content
   - Conversation header with title, timestamp, token summary
   - Message grouping with time-based dividers (>1 hour gaps)
   - User messages: Blue box with 👤 icon
   - Assistant messages: Green box with 🤖 icon
   - Tool use/result display with specialized styling
   - Thinking content collapsible sections
   - Syntax highlighting with chroma (auto-detect language)
   - Code block copy action
   - Token usage per-message breakdown

3. **Sidebar Component** (internal/ui/components/sidebar.go)
   - Project list with folder icons
   - Project metadata (conversation count, last activity)
   - Filters section with quick filters and advanced groups
   - Collapsible filter groups (Time, Cost, Tokens, Tools, Projects)
   - Active filter count indicator
   - Keyboard navigation (arrows, expand/collapse, select)
   - Number keys (1-9) for quick project jumping

4. **Search Implementation** (internal/search/)
   - Lexer tokenizes query strings
   - Parser builds AST from tokens
   - Query node types: Term, And, Or, Not, Field-specific
   - Boolean operator support with proper precedence
   - Field-specific searches: title:, tool:, date:, cost:, tokens:, branch:
   - Date range syntax (today, this-week, 2025-01-15..2025-01-20)
   - Search engine applies query against dataset
   - Results ranked by relevance (phrase > word, title > content, recent)

5. **List View Search Integration**
   - Search input field appearing above conversation list
   - Real-time filtering with 300ms debounce
   - Live result count display
   - Filter tag pills with close buttons
   - "Clear all filters" action
   - Empty state with helpful message

6. **Repository Enhancements**
   - LRU cache wrapper (50-100 conversation capacity)
   - Lazy loading: parse JSONL on-demand
   - Metadata extraction from first/last message
   - Project path decoding (hyphens to slashes)
   - Statistics aggregation (token counts, costs, tool usage)

### Testing Strategy
- Unit tests for search parser with table-driven cases
- Repository tests with mock JSONL files (afero)
- UI component tests with teatest
- Golden file testing for message rendering
- Search relevance ranking tests

### Success Criteria
- Browse and select conversations
- View full message content with proper formatting
- Search conversations with boolean syntax
- Filter by projects and time ranges
- Table sorts correctly with visual indicators
- Code blocks render with syntax highlighting
- Performance acceptable with 1000+ conversations

---

## Phase 3: Advanced Features & Analytics

**Duration:** 2-3 weeks
**Goal:** Add statistics, command palette, and export functionality

### Deliverables

1. **Statistics Dashboard** (internal/ui/views/stats.go)
   - Using NimbleMarkets/ntcharts
   - Overview cards:
     - Total Projects with project directory count
     - Total Conversations with this-month count
     - Total Messages with average per conversation
     - Total Cost with average per conversation
   - Temporal analysis:
     - Line chart: Conversations per day (last 30 days)
     - Stacked area chart: Token usage trends (input/output/cache)
   - Tool usage analysis: Horizontal bar chart with success rates
   - Cost distribution: Pie chart by project
   - Activity heatmap: By hour of day and day of week
   - Date range controls: 7/30/90 days, year, all time, custom
   - Export options: PNG, CSV, JSON

2. **Command Palette** (internal/ui/views/palette.go)
   - Modal overlay with semi-transparent background
   - Fuzzy matching on command names
   - Command categories:
     - Navigation: Go to sidebar/list/detail
     - View: Toggle stats, refresh, compact mode
     - Search: Filter by date/cost/tokens, clear filters
     - Data: Export, copy link, refresh index
     - App: Toggle theme, settings, help, quit
   - Keyboard shortcuts displayed per command
   - Parameter input for commands requiring args
   - Confirmation dialogs for destructive operations
   - Recent commands at top of list

3. **Export Functionality**
   - Markdown export with formatted messages
   - JSON export with complete conversation data
   - CSV export for spreadsheet analysis
   - File save with dialog picker
   - Progress indicator for batch operations
   - Confirmation before overwriting

4. **Help System**
   - Modal with tabbed interface
   - Keyboard shortcuts reference (organized by context)
   - Getting started guide
   - Feature walkthroughs
   - Troubleshooting section
   - Search within help

5. **Settings & Configuration** (internal/app/config.go)
   - YAML configuration file at ~/.config/claudex/settings.yaml
   - Appearance: Theme (auto/light/dark), color accents
   - Behavior: Auto-refresh, animations on/off, confirmation dialogs
   - Data: Cache size limit, retention policy
   - Advanced: Debug mode, performance mode
   - Settings import/export for sync across machines

6. **Animation & Polish**
   - Harmonica spring animations for transitions
   - Spinner for loading states
   - Progress bars with gradient colors
   - Smooth fade transitions between views
   - Status messages with timeout

### Testing Strategy
- Chart rendering tests with sample data
- Command palette fuzzy matching tests
- Export format validation
- Settings persistence tests
- Animation timing tests

### Success Criteria
- Statistics dashboard displays all metrics correctly
- Command palette has 30+ commands with fuzzy search
- Export produces valid formats (JSON parseable, markdown viewable)
- Settings persist across sessions
- All animations smooth and responsive

---

## Phase 4: Optimization, Testing & Polish

**Duration:** 1-2 weeks
**Goal:** Ensure production-quality code with full test coverage

### Deliverables

1. **Performance Optimization**
   - Virtual scrolling for large message lists
   - Lazy loading of conversation details
   - Pagination for conversation lists (50 per page)
   - String interning for common values (role, model)
   - Rendered content caching with dirty flags
   - Memory profiling with pprof
   - Benchmark tests for hot paths

2. **Comprehensive Testing**
   - Unit tests: 85%+ coverage for domain and search
   - Integration tests: Core workflows (browse, search, export)
   - Component tests: All UI views with teatest
   - Error scenario tests: Corrupted files, missing projects
   - End-to-end test suite with VHS recordings

3. **Error Handling & Recovery**
   - Permission denied → suggest chmod or change directory
   - File not found → mark in list with warning icon
   - Corrupted JSONL → skip line, log error, continue
   - Large files → show progress indicator
   - Network errors (future): Retry with exponential backoff
   - Out of memory → suggest filters or archive old data

4. **Documentation**
   - README.md with features and installation
   - Architecture document explaining design decisions
   - Contributing guide for future maintainers
   - Keyboard shortcut reference
   - Configuration guide
   - Troubleshooting section

5. **CI/CD Setup**
   - GitHub Actions workflow for testing
   - Automated linting (golangci-lint)
   - Code coverage reporting
   - Build matrix for multiple OS/architectures
   - Release automation with GoReleaser

6. **Final Polish**
   - Review all keybindings for consistency
   - Test on multiple terminal emulators
   - Color scheme testing for light/dark backgrounds
   - Accessibility review (high contrast, keyboard-only)
   - Load testing with large conversation histories

### Testing Coverage Targets
- Domain layer: 90%+
- Repository layer: 85%+
- Search engine: 95%+
- UI layer: 70%+ (where feasible)
- Overall: 80%+ coverage

### Success Criteria
- All tests pass with clean output
- Code coverage targets met
- Application handles 10,000+ conversations smoothly
- No memory leaks under sustained use
- Response time <100ms for typical operations
- Works on macOS, Linux, Windows
- Professional README and documentation

---

## Technical Specifications

### State Management
- **Single source of truth:** Model struct holds all state
- **Message-driven updates:** All state changes via messages
- **View state machine:** Explicit enum for navigation
- **History preservation:** Save/restore scroll positions when navigating
- **Command pattern:** Async I/O through tea.Cmd functions

### Data Flow Architecture
```
File System (~/.claude/projects/)
       ↓
JSONLRepository (parsing layer)
       ↓
CachedRepository (LRU wrapper)
       ↓
SearchEngine (filter & match)
       ↓
UI Views (render results)
       ↓
Lipgloss (styling)
       ↓
Terminal
```

### Error Handling Strategy
```
Level 1: Parse errors (file → domain)
  → Skip malformed lines, log error, continue

Level 2: Repository errors (data access)
  → Wrapped with context (file path, operation)
  → Retry with exponential backoff for I/O

Level 3: Business logic errors (search, filter)
  → Validation errors for user input
  → Operation errors (division by zero, type mismatch)

Level 4: UI errors (rendering, state)
  → Error modal with recovery options
  → Graceful degradation (show cached data)
```

### Performance Targets
- Startup: <500ms to first render
- Conversation load: <300ms parsing + rendering
- Search: <500ms for typical query on 10k convs
- Scroll: 60 FPS sustained
- Memory: <100MB for 1000 conversations
- Idle CPU: <1%

### Responsive Design Breakpoints
- **Wide (120+ chars):** All three panes visible, all columns
- **Standard (100-119 chars):** All panes, reduced columns
- **Narrow (80-99 chars):** Sidebar hidden by default
- **Minimal (<80 chars):** Single pane with full-screen overlays

---

## Implementation Guidelines

### Code Organization
- One responsibility per file
- Package functions logically
- Interfaces for testability
- Dependency injection throughout
- Minimal global state (only for const colors/styles)

### Naming Conventions
- Views: `*ViewModel` (e.g., ListViewModel)
- Models: `*Model` (e.g., Conversation)
- Messages: `*Msg` (e.g., conversationSelectedMsg)
- Commands: Function returning `tea.Cmd`
- Private helpers: lowercase, short

### Documentation
- Exported functions have comments
- Complex logic has inline explanations
- Structs describe purpose and constraints
- Examples in package-level docs
- Keep ABOUT ME comments in files

### Testing Approach
- Write tests alongside implementation
- Table-driven tests for multiple cases
- Mock files in testdata/
- Golden files for output validation
- Property-based tests where applicable

---

## Key Decision Points

### Library Choices Rationale
1. **bubble-table over bubbles/table:** Sorting, filtering, frozen columns essential
2. **chroma for syntax highlighting:** 100+ languages, Lipgloss integration
3. **ntcharts for visualizations:** Purpose-built for terminal charts
4. **harmonica for animations:** Physics-based, smooth, professional feel

### Architecture Choices
1. **Repository pattern:** Allows easy mocking, future storage backends
2. **Message-driven state:** Testable, predictable, clear data flow
3. **Component composition:** Reusable, maintainable UI components
4. **LRU caching:** Simple, bounded memory, works for typical usage

### Trade-offs
1. **No ORM:** Direct JSONL parsing simpler for this use case
2. **No CLI framework:** Bubble Tea sufficient for all interaction
3. **Single-threaded:** Bubble Tea's async commands sufficient
4. **Limited i18n:** Focus on English initially, architecture supports it

---

## Success Metrics

### Functional Completeness
- [ ] Browse all projects and conversations
- [ ] View full message content with formatting
- [ ] Search with boolean operators and field-specific syntax
- [ ] Filter by multiple criteria simultaneously
- [ ] View statistics and usage patterns
- [ ] Export conversations in multiple formats
- [ ] Navigate efficiently with keyboard only
- [ ] Handle 10,000+ conversations smoothly

### Code Quality
- [ ] 80%+ test coverage
- [ ] Zero panics in normal operation
- [ ] All linting checks pass
- [ ] No memory leaks (detected by pprof)
- [ ] Clean error messages for users

### User Experience
- [ ] Responsive to all inputs (<100ms latency)
- [ ] Intuitive keyboard navigation
- [ ] Clear visual hierarchy with Lipgloss
- [ ] Consistent color scheme and styling
- [ ] Helpful error messages and hints

---

## Next Steps

1. **Immediate:** Set up project structure and dependencies
2. **Week 1:** Implement domain models and repository
3. **Week 2:** Build basic UI shell and navigation
4. **Week 3:** Implement conversation list and detail views
5. **Week 4:** Add search and filtering
6. **Week 5:** Statistics and command palette
7. **Week 6:** Export, settings, help system
8. **Week 7-8:** Optimization, testing, documentation
9. **Week 9:** Final polish and release

---

## References

### Framework Documentation
- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/main/tutorials)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)

### Architecture Resources
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go)
- [Clean Code in Go](https://github.com/golang-standards/project-layout)

### Design Resources
- UX Specification: docs/ux-specs.md
- Architecture Research: docs/bubble-tea-arch.md
- Ecosystem Research: docs/bubble-tea-eco.md

---

**Document Version:** 1.0
**Last Updated:** 2025-11-23
**Status:** Approved - Ready for Implementation
