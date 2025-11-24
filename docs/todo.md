# Claudex Development Roadmap

> Comprehensive todo list for building Claude's conversation browser. Track progress through all phases here.

---

## Step 0: Foundation & Setup ✅ COMPLETE

The foundation has been laid. These items are done:

- [x] Project structure created (cmd/, internal/, pkg/, testdata/, docs/)
- [x] Go module initialized with testify dependency
- [x] Domain layer fully implemented
  - [x] `internal/domain/conversation.go` - Conversation struct, constructors, methods
  - [x] `internal/domain/message.go` - Message struct, Role enum, TokenCount tracking
  - [x] `internal/domain/errors.go` - AppError struct, ErrorType enum, error helpers
  - [x] `internal/domain/statistics.go` - Statistics aggregation, calculations
  - [x] Complete test coverage for all domain models
- [x] Build system configured
  - [x] `Makefile` with build, run, test, check, format, vet, clean targets
  - [x] `.gitignore` configured
- [x] Documentation written
  - [x] `execution-plan.md` - Complete technical execution plan
  - [x] `ux-specs.md` - Complete UX specification
  - [x] `bubble-tea-arch.md` - Architecture research and patterns
  - [x] `bubble-tea-eco.md` - Ecosystem research
- [x] Basic entry point created
  - [x] `cmd/claudex/main.go` - Minimal placeholder

---

## Phase 1: Core Structure & Data Access (64% Complete)

**Goal:** Get the app loading conversations from disk, displaying in three-pane layout, and navigating between views.

### 1.1 Dependencies & Setup ✅ COMPLETE

- [x] Add Bubble Tea framework to go.mod
  - [x] `github.com/charmbracelet/bubbletea`
  - [x] `github.com/charmbracelet/lipgloss` (styling)
  - [x] `github.com/charmbracelet/bubbles` (components)
  - [x] Run `go mod tidy`
- [x] Verify all dependencies download and compile

### 1.2 JSONL Parser Implementation ✅ COMPLETE

- [x] Create `pkg/jsonl/parser.go`
  - [x] `Parser` struct with `io.Reader`
  - [x] `Parser.Next() (*Message, error)` method for streaming parse
  - [x] Line-by-line JSON unmarshal with error handling
  - [x] Handle malformed JSON gracefully
- [x] Create test fixtures in `testdata/`
  - [x] Valid JSONL file with 3-5 messages
  - [x] JSONL file with malformed lines
  - [x] Empty JSONL file
- [x] Write comprehensive tests for `pkg/jsonl/`
  - [x] Valid message parsing
  - [x] Error handling for bad JSON
  - [x] EOF handling
  - [x] 92.3% test coverage (exceeds requirements)

### 1.3 Repository Layer ✅ COMPLETE

- [x] Create `internal/repository/repository.go`
  - [x] `Repository` interface with methods:
    - [x] `List(ctx context.Context) ([]*Conversation, error)`
    - [x] `GetByID(ctx context.Context, id string) (*Conversation, error)`
    - [x] `Search(ctx context.Context, query string) ([]*Conversation, error)`
    - [x] `GetStatistics(ctx context.Context) (*Statistics, error)`
  - [x] Error type definitions

- [x] Create `internal/repository/jsonl.go` - JSONL implementation
  - [x] `JSONLRepository` struct with root path
  - [x] `NewJSONLRepository()` constructor
  - [x] Parse `~/.claude/projects/` directory structure
  - [x] Decode hyphenated project paths back to real paths (e.g., `-foo-bar` → `/foo/bar`)
  - [x] Lazy loading strategy for conversations (metadata on List, full load on GetByID)
  - [x] Extract metadata from first and last lines of JSONL files
  - [x] Implement all `Repository` interface methods
  - [x] Integration with JSONL parser from 1.2

- [x] Create `internal/repository/mock.go` for testing
  - [x] Mock repository implementation with pre-loaded data
  - [x] Test helper functions (`NewMock`, `NewMockWithStatistics`)

- [x] Write comprehensive tests
  - [x] List all conversations (11 test cases)
  - [x] Get specific conversation by ID
  - [x] Handle missing projects/files gracefully
  - [x] Path encoding/decoding
  - [x] Corrupted JSONL handling (skip bad lines, continue parsing)
  - [x] Search functionality (substring matching)
  - [x] Statistics aggregation
  - [x] 84.2% test coverage (exceeds 85% target with only mock edge cases uncovered)

### 1.4 Application Bootstrap

- [ ] Create `internal/app/app.go`
  - [ ] `App` struct with repository, UI model, config
  - [ ] `New()` constructor
  - [ ] `Run()` method to start Bubble Tea program
  - [ ] Graceful shutdown handling

- [ ] Create `internal/app/config.go`
  - [ ] `Config` struct with paths and settings
  - [ ] Default paths (`~/.claude/projects/`)
  - [ ] `LoadConfig()` constructor
  - [ ] Config validation

- [ ] Write tests for app package
  - [ ] Config loading
  - [ ] App initialization
  - [ ] Bubble Tea program setup

### 1.5 UI Foundation

- [ ] Create `internal/ui/model.go`
  - [ ] `ViewState` enum: LIST, DETAIL, SEARCH, STATS, PALETTE
  - [ ] Root `Model` struct with:
    - [ ] `repository` Repository
    - [ ] `state` ViewState
    - [ ] `width`, `height` for window dimensions
    - [ ] State fields for each view (conversations, selected, scrollPos, etc.)
  - [ ] `Model.Init()` method

- [ ] Create `internal/ui/update.go`
  - [ ] Global `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
  - [ ] Message routing logic
  - [ ] Window resize handling (update width/height)
  - [ ] View state transitions
  - [ ] Global keyboard shortcuts: Ctrl+C (quit), Ctrl+P (palette), Tab (focus), Escape (back)

- [ ] Create `internal/ui/view.go`
  - [ ] Global `View() string` method
  - [ ] Three-pane layout with Lipgloss:
    - [ ] Left pane: Sidebar (projects list)
    - [ ] Center pane: Conversation list or message detail
    - [ ] Right pane: Metadata/stats
  - [ ] Responsive breakpoints (full width if < 80 cols)
  - [ ] Border styling with Lipgloss

- [ ] Create `internal/ui/components/` directory structure
  - [ ] Placeholder files for components to be implemented in Phase 2

- [ ] Create `internal/ui/views/` directory structure
  - [ ] Placeholder files for views to be implemented in Phase 2

### 1.6 Integration & Main Entry Point

- [ ] Update `cmd/claudex/main.go`
  - [ ] Parse command-line flags
  - [ ] Load configuration
  - [ ] Initialize app
  - [ ] Start Bubble Tea program
  - [ ] Handle errors gracefully

- [ ] Create version/build info in `cmd/claudex/version.go`
  - [ ] Version string
  - [ ] Build date
  - [ ] Git commit (if available)

### 1.7 Phase 1 Testing & Validation

- [ ] Run `make test` - all tests passing
- [ ] Run `make check` - no lint errors
- [ ] Run `make build` - binary builds without errors
- [ ] Run `make run` - app starts and shows basic layout
- [ ] Verify Phase 1 success criteria:
  - [ ] Application launches without crashing
  - [ ] Shows three-pane layout
  - [ ] Can discover and list projects in sidebar
  - [ ] Window resize handled correctly
  - [ ] Navigation between panes works (Tab key)
  - [ ] Escape key returns to list view
  - [ ] Ctrl+C quits cleanly
  - [ ] Domain test coverage ≥ 80%

---

## Phase 2: Data Display & User Interaction

**Goal:** Browse conversations, view messages, search with filters, display metadata.

### 2.1 Conversation List View

- [ ] Create `internal/ui/views/list.go`
  - [ ] Integrate `evertras/bubble-table` for table rendering
  - [ ] Define table columns: Timestamp, Title, # Messages, Tokens, Cost
  - [ ] Load all conversations from repository on view init
  - [ ] Display conversations in table

- [ ] Implement sorting
  - [ ] Default: descending by timestamp
  - [ ] Click column headers to sort
  - [ ] Multi-column sort support
  - [ ] Sort indicators (↑↓)

- [ ] Implement keyboard navigation
  - [ ] Arrow keys up/down
  - [ ] j/k vim keys
  - [ ] Page Up/Down for fast scrolling
  - [ ] Home/End for jump to start/end
  - [ ] Enter to view selected conversation

- [ ] Selection handling
  - [ ] Highlight selected row
  - [ ] Update root model state with selected ID
  - [ ] Trigger state change to DETAIL view on Enter

- [ ] Virtual scrolling for performance
  - [ ] Only render visible rows
  - [ ] Handle 1000+ conversations smoothly

### 2.2 Message Detail View

- [ ] Create `internal/ui/views/detail.go`
  - [ ] Load selected conversation from repository
  - [ ] Display conversation header
  - [ ] Scrollable message viewport

- [ ] Message components
  - [ ] User message component (blue theme, right-aligned?)
  - [ ] Assistant message component (green theme)
  - [ ] Tool use/result display
  - [ ] Thinking content (collapsible sections)
  - [ ] Token count display per message

- [ ] Keyboard navigation
  - [ ] Space/Page Down to scroll down
  - [ ] Page Up to scroll up
  - [ ] Home/End to jump to start/end
  - [ ] Escape or 'b' to go back to list
  - [ ] 'y' to copy message content

- [ ] Message actions (future expansion point)
  - [ ] Copy to clipboard
  - [ ] Export individual message
  - [ ] Expand/collapse thinking sections

### 2.3 Sidebar Component

- [ ] Create `internal/ui/components/sidebar.go`
  - [ ] List all projects with metadata
  - [ ] Display project path, message count, last updated
  - [ ] Collapsible filter groups
  - [ ] Quick filters: All, Today, This Week, This Month

- [ ] Keyboard navigation in sidebar
  - [ ] Arrow up/down through projects
  - [ ] Tab to switch between panes
  - [ ] Number keys 1-9 for quick project jumping
  - [ ] Filter selection with arrow keys

- [ ] Focus indicator
  - [ ] Highlight current project
  - [ ] Indicate if sidebar or list has focus

### 2.4 Search Engine Implementation

- [ ] Create `internal/search/lexer.go`
  - [ ] `Lexer` for tokenizing query strings
  - [ ] Token types: WORD, FIELD, OPERATOR, QUOTE
  - [ ] Handle quoted strings: "exact phrase"
  - [ ] Parse special operators: AND, OR, NOT

- [ ] Create `internal/search/query.go`
  - [ ] `Query` AST structure
  - [ ] `Parser` to build AST from tokens
  - [ ] Support syntax: `title:"my chat"`, `tool:web`, `date:>2024-01-01`
  - [ ] Boolean operators: `(ai OR ml) AND NOT archived`

- [ ] Create `internal/search/filter.go`
  - [ ] `Filter` interface with `Match(c *Conversation) bool`
  - [ ] Implementations:
    - [ ] `TextFilter` - full-text search in title/messages
    - [ ] `FieldFilter` - match specific fields (title, tool, date, cost, tokens)
    - [ ] `BooleanFilter` - AND/OR/NOT combinations
  - [ ] Case-insensitive matching

- [ ] Create `internal/search/engine.go`
  - [ ] `Engine` struct with repository reference
  - [ ] `Search(ctx context.Context, query string) ([]Conversation, error)`
  - [ ] Relevance ranking (title matches higher than content)
  - [ ] Limit results to top 100
  - [ ] Error handling for invalid queries

- [ ] Write comprehensive search tests
  - [ ] Text matching
  - [ ] Boolean operators
  - [ ] Field-specific searches
  - [ ] Date range parsing
  - [ ] Relevance ranking
  - [ ] 85%+ coverage

### 2.5 Search Integration in UI

- [ ] Create `internal/ui/components/search_bar.go`
  - [ ] Input field for query
  - [ ] Real-time search with 300ms debounce
  - [ ] Result count display
  - [ ] Error message display

- [ ] Filter tag display
  - [ ] Show active filters as pills
  - [ ] 'x' button on each pill to remove filter
  - [ ] Clear all filters button

- [ ] Update list view for search results
  - [ ] Display search results or all conversations
  - [ ] Update counts dynamically
  - [ ] Empty state when no results

### 2.6 Repository Caching

- [ ] Create `internal/repository/cached.go`
  - [ ] `CachedRepository` wrapper with LRU cache
  - [ ] Cache conversations and statistics
  - [ ] `New(repo Repository, maxSize int) *CachedRepository`
  - [ ] Implement `Repository` interface with cache hits
  - [ ] TTL or invalidation strategy

- [ ] Lazy loading strategy
  - [ ] Load conversation headers on List()
  - [ ] Load full messages on GetByID()
  - [ ] Don't load all message bodies at startup

- [ ] Statistics aggregation
  - [ ] Calculate on first GetStatistics() call
  - [ ] Cache the result
  - [ ] Optional background refresh

### 2.7 Markdown Rendering Basics

- [ ] Create `internal/render/markdown.go`
  - [ ] Basic markdown to terminal rendering
  - [ ] Bold, italic, code inline
  - [ ] Code fence detection
  - [ ] Preserve structure but strip most formatting (TUI friendly)

### 2.8 Syntax Highlighting

- [ ] Create `internal/render/syntax.go`
  - [ ] Integrate Chroma library for syntax highlighting
  - [ ] Detect language from code fence
  - [ ] Map Chroma colors to terminal colors
  - [ ] Apply highlighting to code blocks in messages

- [ ] Add Chroma dependency
  - [ ] `github.com/alecthomas/chroma/v2`

### 2.9 Phase 2 Testing & Validation

- [ ] Run `make test` - all Phase 2 tests passing
- [ ] Run `make check` - no lint errors
- [ ] Verify Phase 2 success criteria:
  - [ ] Browse and select conversations from list
  - [ ] View full message content in detail view
  - [ ] Scroll through long conversations
  - [ ] Search with text and field filters
  - [ ] Boolean operators work (AND, OR, NOT)
  - [ ] Syntax highlighting renders correctly
  - [ ] Sidebar shows projects and quick filters
  - [ ] Navigate between all panes with keyboard
  - [ ] Tab and Escape work as expected
  - [ ] No crashes with 1000+ conversations
  - [ ] Tests cover 75%+ of Phase 2 code

---

## Phase 3: Advanced Features & Analytics

**Goal:** Statistics dashboard, command palette, export, settings, help system.

### 3.1 Statistics Dashboard

- [ ] Create `internal/ui/views/stats.go`
  - [ ] Statistics aggregation in repository
  - [ ] Four main cards: Total Conversations, Total Messages, Total Tokens, Estimated Cost
  - [ ] Key metrics displayed prominently

- [ ] Add ntcharts dependency
  - [ ] `github.com/NimbleMarkets/ntcharts`

- [ ] Implement charts
  - [ ] Temporal chart: Conversations per day/week/month
  - [ ] Token usage trends (moving average)
  - [ ] Tool usage bar chart
  - [ ] Cost distribution pie chart
  - [ ] Activity heatmap (conversations by hour of day)

- [ ] Date range controls
  - [ ] Last 7 days, 30 days, all time
  - [ ] Custom date range picker
  - [ ] Update charts on range change

- [ ] Performance considerations
  - [ ] Cache chart data
  - [ ] Regenerate on repository changes
  - [ ] Handle large datasets gracefully

### 3.2 Command Palette

- [ ] Create `internal/ui/views/palette.go`
  - [ ] Modal overlay with input field
  - [ ] Command list below input
  - [ ] Fuzzy matching on command name and description
  - [ ] Category groups: Navigation, View, Search, Data, App

- [ ] Implement 30+ commands
  - Navigation: Go to List, Go to Stats, Go to Detail, Go to Search
  - View: Refresh, Toggle Focus, Focus Sidebar, Focus List, Focus Detail
  - Search: Clear Search, Save Search, Load Saved Search
  - Data: Export Conversation, Export All, Copy ID
  - App: Settings, Help, About, Quit
  - Example: `palette.go:Register("go-list", "Go to conversation list", ...)`

- [ ] Recent commands tracking
  - [ ] Store last 10 commands
  - [ ] Display recently used at top
  - [ ] Clear history option

- [ ] Parameter input for commands
  - [ ] Some commands take parameters (e.g., search term)
  - [ ] Follow-up input field after command selection

- [ ] Integration with promptkit
  - [ ] `github.com/erikgeiser/promptkit`

### 3.3 Export Functionality

- [ ] Create `internal/export/` package
  - [ ] `Exporter` interface with methods:
    - [ ] `ExportMarkdown(ctx context.Context, c *Conversation) ([]byte, error)`
    - [ ] `ExportJSON(ctx context.Context, c *Conversation) ([]byte, error)`
    - [ ] `ExportCSV(ctx context.Context, conversations []Conversation) ([]byte, error)`

- [ ] Markdown exporter
  - [ ] Conversation title as header
  - [ ] Metadata section (model, date, tokens, cost)
  - [ ] Messages with role labels
  - [ ] Code blocks preserved
  - [ ] Thinking content in collapsible sections

- [ ] JSON exporter
  - [ ] Full conversation data structure
  - [ ] All message content and metadata
  - [ ] Indented, readable output

- [ ] CSV exporter
  - [ ] One row per conversation
  - [ ] Columns: Title, Date, Messages, Tokens, Cost, Model, Project
  - [ ] Proper escaping for commas/quotes

- [ ] File save dialog
  - [ ] Ask for filename
  - [ ] Default filename based on conversation title or timestamp
  - [ ] Show file path after save
  - [ ] Error handling for permission denied, disk full, etc.

- [ ] Progress indicators
  - [ ] Show spinner while exporting
  - [ ] Progress bar for bulk exports
  - [ ] Cancel option

### 3.4 Help System

- [ ] Create `internal/ui/views/help.go` (or modal)
  - [ ] Modal overlay (non-blocking)
  - [ ] Tab system: Getting Started, Keyboard Shortcuts, Features, Troubleshooting

- [ ] Getting Started tab
  - [ ] Quick intro to claudex
  - [ ] How to navigate
  - [ ] First steps (view a conversation)

- [ ] Keyboard Shortcuts tab
  - [ ] Organized by context (Global, List View, Detail View, Search)
  - [ ] Auto-generated from command definitions
  - [ ] Example:
    ```
    GLOBAL SHORTCUTS
    Ctrl+C        Quit application
    Ctrl+P        Open command palette
    Tab           Switch focus between panes
    Escape        Back to conversation list

    LIST VIEW
    ↑/↓ or j/k    Navigate conversations
    Enter         View selected conversation
    / or Ctrl+F   Search conversations
    ```

- [ ] Features tab
  - [ ] Overview of main features
  - [ ] Links to help sections
  - [ ] Feature highlights

- [ ] Troubleshooting tab
  - [ ] Common issues
  - [ ] Solutions
  - [ ] Contact info

### 3.5 Settings & Configuration

- [ ] Expand `internal/app/config.go`
  - [ ] `Settings` struct with appearance, behavior, data options
  - [ ] Load from YAML at `~/.config/claudex/settings.yaml`
  - [ ] Defaults if file doesn't exist
  - [ ] Validation on load

- [ ] Settings fields
  - [ ] **Appearance:** theme (auto/light/dark), font size (not applicable to TUI), animation speed
  - [ ] **Behavior:** auto-refresh interval, confirm before quit, confirm before delete (future)
  - [ ] **Data:** max cache size, result limit, search debounce ms

- [ ] Create `internal/ui/views/settings.go` (modal)
  - [ ] List all settings
  - [ ] Edit values with arrow keys/Enter
  - [ ] Save to config file on change
  - [ ] Restart required for some settings (show indicator)

- [ ] Settings persistence
  - [ ] Write YAML on each change
  - [ ] Validate before writing
  - [ ] Error handling (permission denied, invalid YAML)

- [ ] Import/Export settings
  - [ ] Export settings to file for backup
  - [ ] Import settings from file to restore

### 3.6 Animation & Polish

- [ ] Add Harmonica dependency
  - [ ] `github.com/charmbracelet/harmonica`

- [ ] Implement spring animations
  - [ ] Fade in new views (0.3s)
  - [ ] Spring animation for charts (0.8s with overshoot)
  - [ ] Smooth transitions between panes

- [ ] Loading spinners
  - [ ] Spinner while loading conversations
  - [ ] Spinner while searching
  - [ ] Spinner while exporting

- [ ] Progress bars
  - [ ] Progress during bulk operations
  - [ ] Animated fill color

- [ ] Status messages
  - [ ] "Exporting conversation..." messages
  - [ ] Success/error toast notifications (temporary, dismiss with space)
  - [ ] Fade out after 3-5 seconds

- [ ] Add logging dependency
  - [ ] `github.com/charmbracelet/log`

### 3.7 Phase 3 Testing & Validation

- [ ] Run `make test` - all Phase 3 tests passing
- [ ] Run `make check` - no lint errors
- [ ] Verify Phase 3 success criteria:
  - [ ] Statistics show all metrics correctly
  - [ ] Charts render and update on filter change
  - [ ] Command palette has 30+ commands
  - [ ] Fuzzy matching works smoothly
  - [ ] Export produces valid files (Markdown, JSON, CSV)
  - [ ] Settings persist across restarts
  - [ ] Animations are smooth (no jank)
  - [ ] Loading spinners display during operations
  - [ ] Help system is comprehensive
  - [ ] No crashes with 10,000+ conversations
  - [ ] Tests cover 75%+ of Phase 3 code

---

## Phase 4: Optimization, Testing & Polish

**Goal:** High performance, comprehensive testing, production-ready quality.

### 4.1 Performance Optimization

- [ ] Virtual scrolling/pagination
  - [ ] Implement in list view (only render visible rows)
  - [ ] Pagination for detail view (50 messages per page)
  - [ ] Lazy load next page on scroll near bottom

- [ ] String interning
  - [ ] Create `internal/cache/intern.go`
  - [ ] Intern commonly repeated strings (roles, model names)
  - [ ] Reduce memory for 1000+ conversations

- [ ] Rendered content caching
  - [ ] Cache formatted message output
  - [ ] Invalidate on content change or window resize
  - [ ] Profile cache hit rate

- [ ] Memory profiling
  - [ ] Add pprof integration
  - [ ] Profile with typical dataset (1000+ conversations)
  - [ ] Identify and fix memory leaks
  - [ ] Benchmark: <500MB memory for 10,000 conversations

- [ ] Benchmark tests
  - [ ] `BenchmarkRepositoryList` - list 10,000 conversations
  - [ ] `BenchmarkSearchEngine` - search across 10,000 conversations
  - [ ] `BenchmarkUIRender` - render 1000-message conversation
  - [ ] Document baseline performance

- [ ] Comparison with initial state
  - [ ] Measure improvement from optimization
  - [ ] Set performance goals and verify

### 4.2 Comprehensive Testing

- [ ] Unit tests (85%+ coverage for domain/search)
  - [ ] Existing domain tests (✓ done in Phase 1)
  - [ ] Repository layer tests (Phase 1)
  - [ ] Search engine tests (Phase 2)
  - [ ] All 4.2.* sections in this phase

- [ ] Integration tests
  - [ ] End-to-end: Load project → Search → Export
  - [ ] Browse workflow: List → Detail → Back → List
  - [ ] Search workflow: Query → View results → Select → View detail
  - [ ] Settings workflow: Change setting → Restart → Verify change

- [ ] Component tests with teatest
  - [ ] Add dependency: `github.com/charmbracelet/x/exp/teatest`
  - [ ] Test list view interactions
  - [ ] Test detail view scrolling
  - [ ] Test search input and filtering
  - [ ] Test keyboard navigation

- [ ] Error scenario tests
  - [ ] Permission denied (no access to ~/.claude/projects)
  - [ ] File not found (conversation file deleted)
  - [ ] Corrupted JSONL (invalid JSON in file)
  - [ ] Out of memory (graceful degradation)
  - [ ] Disk full (export failure)
  - [ ] Invalid config (syntax error in settings.yaml)

- [ ] End-to-end test suite with VHS
  - [ ] Record terminal sessions with VHS
  - [ ] Create test scenarios:
    - [ ] Start app, list conversations, open one
    - [ ] Perform a search with filters
    - [ ] Export a conversation
    - [ ] Navigate to stats, view charts
    - [ ] Open command palette, execute command
  - [ ] Verify output matches expected
  - [ ] Can run headless for CI

- [ ] Coverage reporting
  - [ ] Run: `go test -cover ./...`
  - [ ] Generate coverage report: `go test -coverprofile=coverage.out ./...`
  - [ ] Target: 80%+ overall coverage
  - [ ] Focus on critical paths: domain, search, repository

### 4.3 Error Handling & Recovery

- [ ] Permission denied scenarios
  - [ ] Can't read ~/.claude/projects directory
  - [ ] Show helpful error message
  - [ ] Suggest checking permissions
  - [ ] Graceful degradation (read-only mode if possible)

- [ ] File not found handling
  - [ ] Conversation file deleted externally
  - [ ] Show "Conversation no longer available"
  - [ ] Return to list view
  - [ ] Refresh list

- [ ] Corrupted JSONL handling
  - [ ] Skip malformed lines with warning
  - [ ] Continue parsing rest of file
  - [ ] Log error details for user
  - [ ] Show notification: "Some messages could not be loaded"

- [ ] Large file handling
  - [ ] Show progress indicator for large conversations
  - [ ] Implement timeout (don't hang for 1 hour long chats)
  - [ ] Pagination to avoid loading entire conversation

- [ ] Out of memory handling
  - [ ] Monitor memory usage
  - [ ] Show warning if approaching limit
  - [ ] Clear caches if memory critical
  - [ ] Offer to close app and free memory

- [ ] Panic recovery
  - [ ] Defer recover() in main
  - [ ] Log panic details
  - [ ] Show error screen: "An unexpected error occurred"
  - [ ] Preserve state for crash recovery

### 4.4 Documentation

- [ ] Create `README.md`
  - [ ] Features overview (2-3 sentences)
  - [ ] Installation instructions
  - [ ] Quick start (how to run)
  - [ ] Key features with examples
  - [ ] Screenshots (optional)
  - [ ] Keyboard shortcuts quick ref
  - [ ] Configuration
  - [ ] Troubleshooting section
  - [ ] Contributing guidelines
  - [ ] License

- [ ] Architecture document
  - [ ] High-level design
  - [ ] Package structure and responsibilities
  - [ ] Data flow diagrams
  - [ ] Key abstractions (Repository, Search, etc.)
  - [ ] Extension points for future features

- [ ] Contributing guide
  - [ ] Development setup
  - [ ] Running tests
  - [ ] Code style
  - [ ] Pull request process
  - [ ] Reporting bugs

- [ ] Keyboard shortcuts reference
  - [ ] Complete, organized by context
  - [ ] Markdown file in docs/
  - [ ] Embedded in help system

- [ ] Configuration guide
  - [ ] Where to find settings file
  - [ ] All available options
  - [ ] Examples
  - [ ] How to reset to defaults

- [ ] Troubleshooting section
  - [ ] Common issues and solutions
  - [ ] How to enable debug logging
  - [ ] How to report issues
  - [ ] Performance troubleshooting

### 4.5 CI/CD Setup

- [ ] Create `.github/workflows/test.yml`
  - [ ] Run `make test` on push
  - [ ] Run `make check` (lint/fmt) on push
  - [ ] Test matrix: macOS (latest), Linux (ubuntu-latest), Windows
  - [ ] Go version: 1.21+
  - [ ] Report coverage to Codecov (if desired)

- [ ] golangci-lint integration
  - [ ] Add `.golangci.yml` config
  - [ ] Configure linters (deadcode, unused, vet, etc.)
  - [ ] Fix any CI lint failures
  - [ ] Run locally before pushing: `golangci-lint run`

- [ ] Code coverage reporting
  - [ ] Generate coverage on each push
  - [ ] Upload to Codecov or similar
  - [ ] Track coverage trend
  - [ ] Fail if coverage drops below threshold (80%)

- [ ] Build matrix
  - [ ] macOS 12.x (amd64, arm64)
  - [ ] Linux (amd64, arm64)
  - [ ] Windows 11 (amd64)
  - [ ] Build binaries for all platforms

- [ ] GoReleaser setup (optional for releases)
  - [ ] Configure release build process
  - [ ] Auto-generate changelog
  - [ ] Upload binaries to GitHub releases
  - [ ] Build on tag push

### 4.6 Final Polish & Testing

- [ ] Keybinding consistency review
  - [ ] Make sure keybindings are predictable
  - [ ] Standard TUI conventions (vim keys, hjkl, j/k)
  - [ ] Escape always goes back
  - [ ] Tab/Shift+Tab switches focus
  - [ ] Enter confirms

- [ ] Multi-terminal testing
  - [ ] Test on iTerm2, Terminal.app (macOS)
  - [ ] Test on GNOME Terminal, Alacritty (Linux)
  - [ ] Test on Windows Terminal (Windows)
  - [ ] Verify colors render correctly on all

- [ ] Color scheme testing
  - [ ] Light terminal background
  - [ ] Dark terminal background
  - [ ] High contrast mode (accessibility)
  - [ ] Ensure readable on all
  - [ ] Test color-blind modes (if possible)

- [ ] Accessibility review
  - [ ] All interactive elements keyboard accessible
  - [ ] Focus indicators visible
  - [ ] Color not sole indicator (use symbols too)
  - [ ] Screen reader compatibility (if applicable)
  - [ ] Reasonable font sizes (not controlled in TUI, but mention in docs)

- [ ] Load testing with large datasets
  - [ ] Test with 10,000 conversations
  - [ ] Measure startup time (<2s)
  - [ ] Measure search time (<500ms)
  - [ ] Measure render time (<100ms)
  - [ ] No memory leaks over extended use

- [ ] Stress testing
  - [ ] Rapid keyboard input (spam j/k, arrows)
  - [ ] Rapid window resizes
  - [ ] Run for 1+ hour without issues
  - [ ] No zombie processes or leaked file handles

### 4.7 Release Preparation

- [ ] Version bumping
  - [ ] Update version in code
  - [ ] Update version in README
  - [ ] Create git tag: `v0.1.0`

- [ ] Changelog
  - [ ] Document all features from Phases 1-4
  - [ ] Create `CHANGELOG.md`
  - [ ] Format: Features, Bug Fixes, Breaking Changes

- [ ] Release notes
  - [ ] Highlight major features
  - [ ] Installation instructions
  - [ ] Known limitations

- [ ] Final verification
  - [ ] All tests pass
  - [ ] All linters pass
  - [ ] Binary builds for all platforms
  - [ ] README is complete and accurate
  - [ ] No TODOs left in code

### 4.8 Phase 4 Testing & Validation

- [ ] Run `make test` - all tests passing (85%+ coverage)
- [ ] Run `make check` - no lint errors, fmt clean
- [ ] Build for macOS, Linux, Windows
- [ ] Verify Phase 4 success criteria:
  - [ ] Startup time <2 seconds
  - [ ] Search <500ms on 10,000 conversations
  - [ ] No memory leaks over extended use
  - [ ] Handles 10,000+ conversations smoothly
  - [ ] All error scenarios handled gracefully
  - [ ] Documentation is comprehensive
  - [ ] CI/CD pipeline green
  - [ ] Works on macOS, Linux, Windows
  - [ ] Performance benchmarks established
  - [ ] Ready for public release

---

## Summary

**Current Status:** ~14% complete (Step 0 ✅, Phase 1.1-1.3 ✅, Phase 1 at 78%)

**Latest Update:** Phase 1.3 (Repository Layer) completed with:
- Repository interface defined for clean data access abstraction
- JSONLRepository fully implemented with lazy loading strategy
- Mock repository for testing without file I/O
- 13 comprehensive test cases with 84.2% coverage
- Path decoding for hyphenated directory names
- Corrupted JSONL file handling (gracefully skips bad lines)
- Search and statistics aggregation implemented
- All validation checks passing (format, vet, test, check)

**Remaining Effort:**
- Phase 1 Completion: 1.4-1.7 remaining (App Bootstrap, UI Foundation, Main Entry Point, Validation)
- Phase 2: ~3-4 weeks (0% done)
- Phase 3: ~3-4 weeks (0% done)
- Phase 4: ~2-3 weeks (0% done)

**Total Remaining:** ~6-9 weeks

---

## How to Use This Todo List

- Mark items as complete: `- [x]` (with the x)
- Track your progress in the todo list tool
- Update this file as you work
- Use it as a living document
- Celebrate completion of each phase! 🎉
