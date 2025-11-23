# Go TUI Application Architecture Research

## Building Complex Terminal Applications with Bubble Tea



**Date:** 2025-11-23

**Framework:** Bubble Tea (Charmbracelet)



---



## Table of Contents



1. [Project Structure](#1-project-structure)

2. [State Management](#2-state-management)

3. [Data Layer Patterns](#3-data-layer-patterns)

4. [Search and Filtering](#4-search-and-filtering)

5. [Testing Strategy](#5-testing-strategy)

6. [Error Handling](#6-error-handling)

7. [Performance Optimization](#7-performance-optimization)

8. [Complete Example](#8-complete-example)



---



## 1. Project Structure



### Recommended Directory Layout



```

claudex/

├── cmd/

│   └── claudex/

│       └── main.go              # Application entry point

├── internal/

│   ├── app/

│   │   ├── app.go              # Main application setup

│   │   └── config.go           # Configuration management

│   ├── ui/

│   │   ├── views/

│   │   │   ├── list.go         # List view component

│   │   │   ├── detail.go       # Detail view component

│   │   │   ├── search.go       # Search view component

│   │   │   ├── stats.go        # Statistics view component

│   │   │   └── palette.go      # Command palette component

│   │   ├── components/

│   │   │   ├── table.go        # Reusable table component

│   │   │   ├── input.go        # Input field component

│   │   │   └── statusbar.go    # Status bar component

│   │   ├── model.go            # Root Bubble Tea model

│   │   ├── update.go           # Update function and messages

│   │   └── view.go             # View rendering

│   ├── domain/

│   │   ├── conversation.go     # Core domain entities

│   │   ├── message.go

│   │   └── statistics.go

│   ├── repository/

│   │   ├── conversation.go     # Repository interface

│   │   ├── jsonl.go           # JSONL implementation

│   │   └── cached.go          # Cached repository wrapper

│   ├── search/

│   │   ├── engine.go          # Search engine interface

│   │   ├── query.go           # Query parser

│   │   ├── filter.go          # Filter composition

│   │   └── index.go           # Search index

│   └── cache/

│       ├── lru.go             # LRU cache implementation

│       └── cache.go           # Cache interface

├── pkg/                        # Optional: public libraries

│   └── jsonl/

│       └── parser.go          # Reusable JSONL parser

└── testdata/                   # Test fixtures

    └── conversations/

        └── sample.jsonl

```



### Architecture Overview



```

┌─────────────────────────────────────────────────────────────┐

│                         cmd/claudex                          │

│                      (Entry Point)                           │

└────────────────────────┬────────────────────────────────────┘

                         │

                         ▼

┌─────────────────────────────────────────────────────────────┐

│                     internal/app                             │

│              (Application Bootstrap)                         │

│  • Configuration loading                                     │

│  • Dependency injection                                      │

│  • Service initialization                                    │

└────────────────────────┬────────────────────────────────────┘

                         │

        ┌────────────────┼────────────────┐

        │                │                │

        ▼                ▼                ▼

┌──────────────┐  ┌──────────────┐  ┌──────────────┐

│ internal/ui  │  │internal/search│  │internal/repo │

│              │  │               │  │              │

│ • Views      │  │ • Engine      │  │ • Interface  │

│ • Components │  │ • Parser      │  │ • JSONL Impl │

│ • Model      │  │ • Filters     │  │ • Cache Wrap │

└──────┬───────┘  └──────┬────────┘  └──────┬───────┘

       │                 │                  │

       └─────────────────┼──────────────────┘

                         │

                         ▼

               ┌──────────────────┐

               │ internal/domain  │

               │                  │

               │ • Entities       │

               │ • Value Objects  │

               │ • Business Logic │

               └──────────────────┘

```



### Key Principles



1. **internal/** - Private application code that cannot be imported by external projects

2. **cmd/** - Application entry points (one directory per binary)

3. **pkg/** - Optional public libraries that could be reused (use sparingly)

4. **Separation of Concerns:**

   - **UI Layer** (`internal/ui`) - Bubble Tea components, views, rendering

   - **Business Layer** (`internal/domain`) - Core business logic, domain models

   - **Data Layer** (`internal/repository`) - Data access, persistence

   - **Service Layer** (`internal/search`, `internal/cache`) - Cross-cutting concerns



---



## 2. State Management



### The Elm Architecture in Bubble Tea



Bubble Tea follows The Elm Architecture pattern with three core components:



```

┌─────────────────────────────────────────────────────────────┐

│                    THE ELM ARCHITECTURE                      │

│                                                               │

│  ┌─────────┐      ┌─────────┐      ┌─────────┐             │

│  │  MODEL  │─────▶│  VIEW   │─────▶│   MSG   │             │

│  │ (State) │      │(Render) │      │(Events) │             │

│  └────▲────┘      └─────────┘      └────┬────┘             │

│       │                                  │                   │

│       │         ┌─────────┐              │                   │

│       └─────────│ UPDATE  │◀─────────────┘                   │

│                 │(Handler)│                                  │

│                 └─────────┘                                  │

│                                                               │

│  1. Model holds application state                            │

│  2. View renders state to UI                                 │

│  3. User interactions generate Messages                      │

│  4. Update processes messages and returns new state          │

└─────────────────────────────────────────────────────────────┘

```



### Multi-View State Management



#### State Enum Pattern



```go

// internal/ui/model.go

package ui



import (

    "github.com/charmbracelet/bubbletea"

    "github.com/yourusername/claudex/internal/domain"

    "github.com/yourusername/claudex/internal/repository"

    "github.com/yourusername/claudex/internal/search"

)



// ViewState represents the current view

type ViewState int



const (

    ViewList ViewState = iota

    ViewDetail

    ViewSearch

    ViewStatistics

    ViewCommandPalette

)



// Model is the root Bubble Tea model

type Model struct {

    // State management

    currentView ViewState

    previousView ViewState  // For returning from command palette



    // View models (sub-components)

    listView      ListViewModel

    detailView    DetailViewModel

    searchView    SearchViewModel

    statsView     StatsViewModel

    paletteView   PaletteViewModel



    // Services (dependency injection)

    repo          repository.ConversationRepository

    searchEngine  search.Engine



    // Shared state

    conversations []domain.Conversation

    selectedID    string

    errorMessage  string



    // UI state

    width         int

    height        int

    ready         bool

}



// Init initializes the Bubble Tea program

func (m Model) Init() tea.Cmd {

    return tea.Batch(

        m.loadConversations(),

        tea.EnterAltScreen,

    )

}



// Update handles messages and state transitions

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        // Global keybindings

        switch msg.String() {

        case "ctrl+c", "q":

            return m, tea.Quit

        case "ctrl+p":

            return m.showCommandPalette()

        case "esc":

            if m.currentView == ViewCommandPalette {

                return m.hideCommandPalette()

            }

        }



    case tea.WindowSizeMsg:

        m.width = msg.Width

        m.height = msg.Height

        m.ready = true

        return m, nil



    case conversationsLoadedMsg:

        m.conversations = msg.conversations

        m.errorMessage = ""

        return m, nil



    case errMsg:

        m.errorMessage = msg.Error()

        return m, nil



    case changeViewMsg:

        return m.changeView(msg.view, msg.context)

    }



    // Delegate to current view

    return m.updateCurrentView(msg)

}



// View renders the current view

func (m Model) View() string {

    if !m.ready {

        return "Initializing..."

    }



    // Render current view

    switch m.currentView {

    case ViewList:

        return m.listView.View()

    case ViewDetail:

        return m.detailView.View()

    case ViewSearch:

        return m.searchView.View()

    case ViewStatistics:

        return m.statsView.View()

    case ViewCommandPalette:

        // Overlay command palette on current view

        return m.renderWithOverlay(m.paletteView.View())

    default:

        return "Unknown view"

    }

}



// Messages for state transitions

type changeViewMsg struct {

    view    ViewState

    context interface{}

}



type conversationsLoadedMsg struct {

    conversations []domain.Conversation

}



type errMsg struct{ error }



// Helper methods

func (m Model) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd



    switch m.currentView {

    case ViewList:

        m.listView, cmd = m.listView.Update(msg)

    case ViewDetail:

        m.detailView, cmd = m.detailView.Update(msg)

    case ViewSearch:

        m.searchView, cmd = m.searchView.Update(msg)

    case ViewStatistics:

        m.statsView, cmd = m.statsView.Update(msg)

    case ViewCommandPalette:

        m.paletteView, cmd = m.paletteView.Update(msg)

    }



    return m, cmd

}



func (m Model) changeView(view ViewState, context interface{}) (tea.Model, tea.Cmd) {

    m.previousView = m.currentView

    m.currentView = view



    // Initialize new view with context

    var cmd tea.Cmd

    switch view {

    case ViewDetail:

        if id, ok := context.(string); ok {

            m.selectedID = id

            cmd = m.loadConversationDetail(id)

        }

    case ViewSearch:

        if query, ok := context.(string); ok {

            cmd = m.searchView.SetQuery(query)

        }

    }



    return m, cmd

}

```



#### View Component Pattern



```go

// internal/ui/views/list.go

package views



import (

    "github.com/charmbracelet/bubbles/list"

    "github.com/charmbracelet/bubbles/key"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/yourusername/claudex/internal/domain"

)



// ListViewModel manages the conversation list view

type ListViewModel struct {

    list          list.Model

    conversations []domain.Conversation

    keys          listKeyMap

}



type listKeyMap struct {

    Enter  key.Binding

    Detail key.Binding

    Search key.Binding

}



func NewListViewModel() ListViewModel {

    // Create list with default delegate

    l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

    l.Title = "Conversations"

    l.SetShowStatusBar(true)

    l.SetFilteringEnabled(true)



    return ListViewModel{

        list: l,

        keys: listKeyMap{

            Enter: key.NewBinding(

                key.WithKeys("enter"),

                key.WithHelp("enter", "view details"),

            ),

            Detail: key.NewBinding(

                key.WithKeys("d"),

                key.WithHelp("d", "details"),

            ),

            Search: key.NewBinding(

                key.WithKeys("/"),

                key.WithHelp("/", "search"),

            ),

        },

    }

}



func (m ListViewModel) Init() tea.Cmd {

    return nil

}



func (m ListViewModel) Update(msg tea.Msg) (ListViewModel, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch {

        case key.Matches(msg, m.keys.Enter), key.Matches(msg, m.keys.Detail):

            if item := m.list.SelectedItem(); item != nil {

                conv := item.(conversationItem)

                // Send message to change view

                return m, func() tea.Msg {

                    return changeViewMsg{

                        view:    ViewDetail,

                        context: conv.ID,

                    }

                }

            }

        case key.Matches(msg, m.keys.Search):

            return m, func() tea.Msg {

                return changeViewMsg{view: ViewSearch}

            }

        }



    case conversationsUpdatedMsg:

        m.conversations = msg.conversations

        items := make([]list.Item, len(msg.conversations))

        for i, conv := range msg.conversations {

            items[i] = conversationItem{Conversation: conv}

        }

        return m, m.list.SetItems(items)

    }



    var cmd tea.Cmd

    m.list, cmd = m.list.Update(msg)

    return m, cmd

}



func (m ListViewModel) View() string {

    return m.list.View()

}



// conversationItem implements list.Item interface

type conversationItem struct {

    domain.Conversation

}



func (i conversationItem) FilterValue() string {

    return i.Title

}



func (i conversationItem) Title() string {

    return i.Conversation.Title

}



func (i conversationItem) Description() string {

    return i.Conversation.Summary()

}

```



### State Preservation Pattern



```go

// internal/ui/state.go

package ui



import "github.com/yourusername/claudex/internal/domain"



// StateSnapshot preserves view state during navigation

type StateSnapshot struct {

    // List view state

    listScrollPosition int

    listFilterQuery    string



    // Search view state

    searchQuery        string

    searchResults      []domain.Conversation



    // Detail view state

    detailScrollPosition int

    detailConversationID string

}



// SaveState creates a snapshot of current view state

func (m Model) SaveState() StateSnapshot {

    return StateSnapshot{

        listScrollPosition: m.listView.list.Index(),

        listFilterQuery:    m.listView.list.FilterValue(),

        searchQuery:        m.searchView.query,

        searchResults:      m.searchView.results,

        detailScrollPosition: m.detailView.scrollPosition,

        detailConversationID: m.detailView.conversationID,

    }

}



// RestoreState restores view state from snapshot

func (m *Model) RestoreState(snapshot StateSnapshot) {

    // Restore list view state

    m.listView.list.Select(snapshot.listScrollPosition)

    if snapshot.listFilterQuery != "" {

        m.listView.list.SetFilterValue(snapshot.listFilterQuery)

    }



    // Restore search view state

    m.searchView.query = snapshot.searchQuery

    m.searchView.results = snapshot.searchResults



    // Restore detail view state

    m.detailView.scrollPosition = snapshot.detailScrollPosition

    m.detailView.conversationID = snapshot.detailConversationID

}



// Enhanced model with state preservation

type ModelWithHistory struct {

    Model

    history []StateSnapshot

}



func (m *ModelWithHistory) PushState() {

    m.history = append(m.history, m.SaveState())

}



func (m *ModelWithHistory) PopState() {

    if len(m.history) > 0 {

        snapshot := m.history[len(m.history)-1]

        m.history = m.history[:len(m.history)-1]

        m.RestoreState(snapshot)

    }

}

```



---



## 3. Data Layer Patterns



### Repository Pattern



```go

// internal/repository/conversation.go

package repository



import (

    "context"

    "github.com/yourusername/claudex/internal/domain"

)



// ConversationRepository defines data access interface

type ConversationRepository interface {

    // List returns all conversations

    List(ctx context.Context) ([]domain.Conversation, error)



    // GetByID returns a conversation by ID

    GetByID(ctx context.Context, id string) (*domain.Conversation, error)



    // Search finds conversations matching criteria

    Search(ctx context.Context, query string) ([]domain.Conversation, error)



    // GetStatistics returns conversation statistics

    GetStatistics(ctx context.Context) (*domain.Statistics, error)

}



// JSONL Implementation

// internal/repository/jsonl.go

package repository



import (

    "bufio"

    "context"

    "encoding/json"

    "fmt"

    "os"

    "path/filepath"



    "github.com/yourusername/claudex/internal/domain"

)



type JSONLRepository struct {

    dataDir string

}



func NewJSONLRepository(dataDir string) *JSONLRepository {

    return &JSONLRepository{

        dataDir: dataDir,

    }

}



func (r *JSONLRepository) List(ctx context.Context) ([]domain.Conversation, error) {

    files, err := filepath.Glob(filepath.Join(r.dataDir, "*.jsonl"))

    if err != nil {

        return nil, fmt.Errorf("failed to list files: %w", err)

    }



    var conversations []domain.Conversation

    for _, file := range files {

        conv, err := r.parseConversation(ctx, file)

        if err != nil {

            // Log error but continue processing other files

            continue

        }

        conversations = append(conversations, *conv)

    }



    return conversations, nil

}



func (r *JSONLRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {

    filename := filepath.Join(r.dataDir, id+".jsonl")

    return r.parseConversation(ctx, filename)

}



// parseConversation uses streaming JSON parsing for memory efficiency

func (r *JSONLRepository) parseConversation(ctx context.Context, filename string) (*domain.Conversation, error) {

    file, err := os.Open(filename)

    if err != nil {

        return nil, fmt.Errorf("failed to open file: %w", err)

    }

    defer file.Close()



    var conv domain.Conversation

    scanner := bufio.NewScanner(file)



    // Increase buffer size for large lines

    buf := make([]byte, 0, 64*1024)

    scanner.Buffer(buf, 1024*1024) // 1MB max token size



    lineNum := 0

    for scanner.Scan() {

        lineNum++



        // Check context cancellation

        select {

        case <-ctx.Done():

            return nil, ctx.Err()

        default:

        }



        var msg domain.Message

        if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {

            return nil, fmt.Errorf("failed to parse line %d: %w", lineNum, err)

        }



        conv.Messages = append(conv.Messages, msg)

    }



    if err := scanner.Err(); err != nil {

        return nil, fmt.Errorf("scanner error: %w", err)

    }



    // Extract metadata from first message

    if len(conv.Messages) > 0 {

        conv.ID = extractID(filename)

        conv.Title = extractTitle(conv.Messages)

        conv.CreatedAt = conv.Messages[0].Timestamp

        conv.UpdatedAt = conv.Messages[len(conv.Messages)-1].Timestamp

    }



    return &conv, nil

}

```



### Cached Repository Pattern (Decorator)



```go

// internal/repository/cached.go

package repository



import (

    "context"

    "fmt"

    "sync"



    "github.com/yourusername/claudex/internal/cache"

    "github.com/yourusername/claudex/internal/domain"

)



// CachedRepository wraps a repository with LRU caching

type CachedRepository struct {

    repo  ConversationRepository

    cache cache.Cache[string, *domain.Conversation]

    mu    sync.RWMutex

}



func NewCachedRepository(repo ConversationRepository, cacheSize int) *CachedRepository {

    return &CachedRepository{

        repo:  repo,

        cache: cache.NewLRU[string, *domain.Conversation](cacheSize),

    }

}



func (r *CachedRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {

    // Try cache first

    r.mu.RLock()

    if conv, ok := r.cache.Get(id); ok {

        r.mu.RUnlock()

        return conv, nil

    }

    r.mu.RUnlock()



    // Cache miss - fetch from underlying repository

    conv, err := r.repo.GetByID(ctx, id)

    if err != nil {

        return nil, err

    }



    // Store in cache

    r.mu.Lock()

    r.cache.Put(id, conv)

    r.mu.Unlock()



    return conv, nil

}



func (r *CachedRepository) List(ctx context.Context) ([]domain.Conversation, error) {

    // List operations typically don't benefit from caching

    // as they need fresh data

    return r.repo.List(ctx)

}



func (r *CachedRepository) Search(ctx context.Context, query string) ([]domain.Conversation, error) {

    // Search results could be cached with query as key

    // but requires careful cache invalidation

    return r.repo.Search(ctx, query)

}



func (r *CachedRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {

    return r.repo.GetStatistics(ctx)

}



// Invalidate removes an entry from cache

func (r *CachedRepository) Invalidate(id string) {

    r.mu.Lock()

    defer r.mu.Unlock()

    r.cache.Remove(id)

}



// Clear removes all entries from cache

func (r *CachedRepository) Clear() {

    r.mu.Lock()

    defer r.mu.Unlock()

    r.cache.Clear()

}

```



### LRU Cache Implementation



```go

// internal/cache/lru.go

package cache



import (

    "container/list"

    "sync"

)



// Cache is a generic cache interface

type Cache[K comparable, V any] interface {

    Get(key K) (V, bool)

    Put(key K, value V)

    Remove(key K)

    Clear()

    Len() int

}



// LRU is a thread-safe LRU cache implementation

type LRU[K comparable, V any] struct {

    capacity int

    items    map[K]*list.Element

    evictList *list.List

    mu       sync.RWMutex

}



type entry[K comparable, V any] struct {

    key   K

    value V

}



// NewLRU creates a new LRU cache with the given capacity

func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {

    return &LRU[K, V]{

        capacity:  capacity,

        items:     make(map[K]*list.Element),

        evictList: list.New(),

    }

}



func (c *LRU[K, V]) Get(key K) (V, bool) {

    c.mu.Lock()

    defer c.mu.Unlock()



    if elem, ok := c.items[key]; ok {

        // Move to front (most recently used)

        c.evictList.MoveToFront(elem)

        return elem.Value.(*entry[K, V]).value, true

    }



    var zero V

    return zero, false

}



func (c *LRU[K, V]) Put(key K, value V) {

    c.mu.Lock()

    defer c.mu.Unlock()



    // If key exists, update value and move to front

    if elem, ok := c.items[key]; ok {

        c.evictList.MoveToFront(elem)

        elem.Value.(*entry[K, V]).value = value

        return

    }



    // Add new entry

    elem := c.evictList.PushFront(&entry[K, V]{key, value})

    c.items[key] = elem



    // Evict oldest if over capacity

    if c.evictList.Len() > c.capacity {

        c.evictOldest()

    }

}



func (c *LRU[K, V]) Remove(key K) {

    c.mu.Lock()

    defer c.mu.Unlock()



    if elem, ok := c.items[key]; ok {

        c.removeElement(elem)

    }

}



func (c *LRU[K, V]) Clear() {

    c.mu.Lock()

    defer c.mu.Unlock()



    c.items = make(map[K]*list.Element)

    c.evictList.Init()

}



func (c *LRU[K, V]) Len() int {

    c.mu.RLock()

    defer c.mu.RUnlock()

    return c.evictList.Len()

}



func (c *LRU[K, V]) evictOldest() {

    elem := c.evictList.Back()

    if elem != nil {

        c.removeElement(elem)

    }

}



func (c *LRU[K, V]) removeElement(elem *list.Element) {

    c.evictList.Remove(elem)

    kv := elem.Value.(*entry[K, V])

    delete(c.items, kv.key)

    // Set to nil to prevent memory leaks

    kv.value = *new(V)

}

```



---



## 4. Search and Filtering



### Search Engine Architecture



```

┌─────────────────────────────────────────────────────────────┐

│                      SEARCH ARCHITECTURE                     │

│                                                               │

│  ┌─────────────┐      ┌──────────────┐     ┌──────────────┐│

│  │Query String │─────▶│Query Parser  │────▶│Query Object  ││

│  │"foo AND bar"│      │(Lexer+Parser)│     │(AST)         ││

│  └─────────────┘      └──────────────┘     └──────┬───────┘│

│                                                     │         │

│                                                     ▼         │

│  ┌─────────────┐      ┌──────────────┐     ┌──────────────┐│

│  │   Results   │◀─────│Search Engine │◀────│Filter Builder││

│  │   []Conv    │      │(Matcher)     │     │(Composition) ││

│  └─────────────┘      └──────────────┘     └──────────────┘│

│                              ▲                               │

│                              │                               │

│                       ┌──────┴───────┐                      │

│                       │Index/Dataset │                       │

│                       │[]Conversation│                       │

│                       └──────────────┘                       │

└─────────────────────────────────────────────────────────────┘

```



### Query Parser Implementation



```go

// internal/search/query.go

package search



import (

    "fmt"

    "strings"

    "unicode"

)



// Query represents a parsed search query

type Query struct {

    Root QueryNode

}



// QueryNode is a node in the query AST

type QueryNode interface {

    String() string

}



// TermNode represents a search term

type TermNode struct {

    Field string // empty for any field

    Value string

}



func (n TermNode) String() string {

    if n.Field != "" {

        return fmt.Sprintf("%s:%s", n.Field, n.Value)

    }

    return n.Value

}



// AndNode represents AND operation

type AndNode struct {

    Left  QueryNode

    Right QueryNode

}



func (n AndNode) String() string {

    return fmt.Sprintf("(%s AND %s)", n.Left, n.Right)

}



// OrNode represents OR operation

type OrNode struct {

    Left  QueryNode

    Right QueryNode

}



func (n OrNode) String() string {

    return fmt.Sprintf("(%s OR %s)", n.Left, n.Right)

}



// NotNode represents NOT operation

type NotNode struct {

    Child QueryNode

}



func (n NotNode) String() string {

    return fmt.Sprintf("NOT %s", n.Child)

}



// Token types

type tokenType int



const (

    tokenTerm tokenType = iota

    tokenAnd

    tokenOr

    tokenNot

    tokenLParen

    tokenRParen

    tokenField

    tokenEOF

)



type token struct {

    typ   tokenType

    value string

}



// Lexer tokenizes the query string

type lexer struct {

    input string

    pos   int

}



func newLexer(input string) *lexer {

    return &lexer{input: input}

}



func (l *lexer) nextToken() token {

    // Skip whitespace

    for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {

        l.pos++

    }



    if l.pos >= len(l.input) {

        return token{typ: tokenEOF}

    }



    // Check for special characters

    ch := l.input[l.pos]

    switch ch {

    case '(':

        l.pos++

        return token{typ: tokenLParen}

    case ')':

        l.pos++

        return token{typ: tokenRParen}

    }



    // Read word

    start := l.pos

    for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) &&

        l.input[l.pos] != '(' && l.input[l.pos] != ')' {

        l.pos++

    }



    word := l.input[start:l.pos]



    // Check for keywords

    upper := strings.ToUpper(word)

    switch upper {

    case "AND":

        return token{typ: tokenAnd}

    case "OR":

        return token{typ: tokenOr}

    case "NOT":

        return token{typ: tokenNot}

    }



    // Check for field:value syntax

    if strings.Contains(word, ":") {

        parts := strings.SplitN(word, ":", 2)

        return token{typ: tokenField, value: word}

    }



    return token{typ: tokenTerm, value: word}

}



// Parser builds AST from tokens

type parser struct {

    lexer   *lexer

    current token

}



func newParser(input string) *parser {

    p := &parser{lexer: newLexer(input)}

    p.advance()

    return p

}



func (p *parser) advance() {

    p.current = p.lexer.nextToken()

}



func (p *parser) Parse() (*Query, error) {

    root, err := p.parseExpression()

    if err != nil {

        return nil, err

    }

    return &Query{Root: root}, nil

}



func (p *parser) parseExpression() (QueryNode, error) {

    return p.parseOr()

}



func (p *parser) parseOr() (QueryNode, error) {

    left, err := p.parseAnd()

    if err != nil {

        return nil, err

    }



    for p.current.typ == tokenOr {

        p.advance()

        right, err := p.parseAnd()

        if err != nil {

            return nil, err

        }

        left = OrNode{Left: left, Right: right}

    }



    return left, nil

}



func (p *parser) parseAnd() (QueryNode, error) {

    left, err := p.parseUnary()

    if err != nil {

        return nil, err

    }



    for p.current.typ == tokenAnd ||

        (p.current.typ == tokenTerm || p.current.typ == tokenField || p.current.typ == tokenLParen) {



        // Implicit AND

        if p.current.typ != tokenAnd {

            right, err := p.parseUnary()

            if err != nil {

                return nil, err

            }

            left = AndNode{Left: left, Right: right}

        } else {

            p.advance()

            right, err := p.parseUnary()

            if err != nil {

                return nil, err

            }

            left = AndNode{Left: left, Right: right}

        }

    }



    return left, nil

}



func (p *parser) parseUnary() (QueryNode, error) {

    if p.current.typ == tokenNot {

        p.advance()

        child, err := p.parsePrimary()

        if err != nil {

            return nil, err

        }

        return NotNode{Child: child}, nil

    }



    return p.parsePrimary()

}



func (p *parser) parsePrimary() (QueryNode, error) {

    switch p.current.typ {

    case tokenLParen:

        p.advance()

        expr, err := p.parseExpression()

        if err != nil {

            return nil, err

        }

        if p.current.typ != tokenRParen {

            return nil, fmt.Errorf("expected ')'")

        }

        p.advance()

        return expr, nil



    case tokenField:

        parts := strings.SplitN(p.current.value, ":", 2)

        node := TermNode{

            Field: parts[0],

            Value: parts[1],

        }

        p.advance()

        return node, nil



    case tokenTerm:

        node := TermNode{Value: p.current.value}

        p.advance()

        return node, nil



    default:

        return nil, fmt.Errorf("unexpected token: %v", p.current)

    }

}

```



### Search Engine Implementation



```go

// internal/search/engine.go

package search



import (

    "context"

    "strings"



    "github.com/yourusername/claudex/internal/domain"

)



// Engine performs searches on conversations

type Engine interface {

    Search(ctx context.Context, query string, conversations []domain.Conversation) ([]domain.Conversation, error)

}



// DefaultEngine implements full-text search with boolean operators

type DefaultEngine struct{}



func NewEngine() *DefaultEngine {

    return &DefaultEngine{}

}



func (e *DefaultEngine) Search(ctx context.Context, queryStr string, conversations []domain.Conversation) ([]domain.Conversation, error) {

    // Parse query

    parser := newParser(queryStr)

    query, err := parser.Parse()

    if err != nil {

        // Fallback to simple text search

        return e.simpleSearch(queryStr, conversations), nil

    }



    // Filter conversations using query

    var results []domain.Conversation

    for _, conv := range conversations {

        if e.matches(query.Root, &conv) {

            results = append(results, conv)

        }

    }



    return results, nil

}



func (e *DefaultEngine) matches(node QueryNode, conv *domain.Conversation) bool {

    switch n := node.(type) {

    case TermNode:

        return e.matchesTerm(n, conv)

    case AndNode:

        return e.matches(n.Left, conv) && e.matches(n.Right, conv)

    case OrNode:

        return e.matches(n.Left, conv) || e.matches(n.Right, conv)

    case NotNode:

        return !e.matches(n.Child, conv)

    default:

        return false

    }

}



func (e *DefaultEngine) matchesTerm(term TermNode, conv *domain.Conversation) bool {

    if term.Field != "" {

        // Field-specific search

        return e.matchesField(term.Field, term.Value, conv)

    }



    // Full-text search across all fields

    return e.matchesFullText(term.Value, conv)

}



func (e *DefaultEngine) matchesField(field, value string, conv *domain.Conversation) bool {

    value = strings.ToLower(value)



    switch strings.ToLower(field) {

    case "title":

        return strings.Contains(strings.ToLower(conv.Title), value)

    case "model":

        return strings.Contains(strings.ToLower(conv.Model), value)

    case "role":

        for _, msg := range conv.Messages {

            if strings.Contains(strings.ToLower(msg.Role), value) {

                return true

            }

        }

    case "content":

        for _, msg := range conv.Messages {

            if strings.Contains(strings.ToLower(msg.Content), value) {

                return true

            }

        }

    }



    return false

}



func (e *DefaultEngine) matchesFullText(value string, conv *domain.Conversation) bool {

    value = strings.ToLower(value)



    // Search in title

    if strings.Contains(strings.ToLower(conv.Title), value) {

        return true

    }



    // Search in messages

    for _, msg := range conv.Messages {

        if strings.Contains(strings.ToLower(msg.Content), value) {

            return true

        }

    }



    return false

}



func (e *DefaultEngine) simpleSearch(query string, conversations []domain.Conversation) []domain.Conversation {

    query = strings.ToLower(query)

    var results []domain.Conversation



    for _, conv := range conversations {

        if e.matchesFullText(query, &conv) {

            results = append(results, conv)

        }

    }



    return results

}

```



### Filter Composition



```go

// internal/search/filter.go

package search



import "github.com/yourusername/claudex/internal/domain"



// Filter is a function that filters conversations

type Filter func(*domain.Conversation) bool



// And combines filters with AND logic

func And(filters ...Filter) Filter {

    return func(conv *domain.Conversation) bool {

        for _, f := range filters {

            if !f(conv) {

                return false

            }

        }

        return true

    }

}



// Or combines filters with OR logic

func Or(filters ...Filter) Filter {

    return func(conv *domain.Conversation) bool {

        for _, f := range filters {

            if f(conv) {

                return true

            }

        }

        return false

    }

}



// Not negates a filter

func Not(filter Filter) Filter {

    return func(conv *domain.Conversation) bool {

        return !filter(conv)

    }

}



// Predefined filters

func ModelFilter(model string) Filter {

    return func(conv *domain.Conversation) bool {

        return conv.Model == model

    }

}



func HasToolUseFilter() Filter {

    return func(conv *domain.Conversation) bool {

        for _, msg := range conv.Messages {

            if msg.HasToolUse {

                return true

            }

        }

        return false

    }

}



func MinMessagesFilter(min int) Filter {

    return func(conv *domain.Conversation) bool {

        return len(conv.Messages) >= min

    }

}



// Example: Complex filter composition

// (model:sonnet OR model:opus) AND has_tool_use AND min_messages:10

func ExampleComplexFilter() Filter {

    return And(

        Or(

            ModelFilter("claude-sonnet-4"),

            ModelFilter("claude-opus-4"),

        ),

        HasToolUseFilter(),

        MinMessagesFilter(10),

    )

}

```



---



## 5. Testing Strategy



### Test Structure



```

internal/

├── ui/

│   ├── model_test.go

│   └── views/

│       └── list_test.go

├── repository/

│   ├── jsonl_test.go

│   └── cached_test.go

├── search/

│   ├── query_test.go

│   ├── engine_test.go

│   └── filter_test.go

└── testutil/

    ├── fixtures.go

    ├── mock_repo.go

    └── golden.go

```



### Unit Testing Bubble Tea Components



```go

// internal/ui/views/list_test.go

package views



import (

    "testing"



    tea "github.com/charmbracelet/bubbletea"

    "github.com/charmbracelet/x/exp/teatest"

    "github.com/yourusername/claudex/internal/domain"

)



func TestListViewModel_Update(t *testing.T) {

    tests := []struct {

        name     string

        initial  ListViewModel

        msg      tea.Msg

        wantView ViewState

        wantCmd  bool

    }{

        {

            name:    "enter key changes to detail view",

            initial: NewListViewModel(),

            msg:     tea.KeyMsg{Type: tea.KeyEnter},

            wantView: ViewDetail,

            wantCmd:  true,

        },

        {

            name:    "slash key changes to search view",

            initial: NewListViewModel(),

            msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}},

            wantView: ViewSearch,

            wantCmd:  true,

        },

    }



    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            model, cmd := tt.initial.Update(tt.msg)



            if (cmd != nil) != tt.wantCmd {

                t.Errorf("Update() cmd = %v, wantCmd %v", cmd != nil, tt.wantCmd)

            }



            if cmd != nil {

                msg := cmd()

                if changeMsg, ok := msg.(changeViewMsg); ok {

                    if changeMsg.view != tt.wantView {

                        t.Errorf("Update() view = %v, want %v", changeMsg.view, tt.wantView)

                    }

                }

            }



            _ = model // Use model if needed

        })

    }

}



func TestListViewModel_View(t *testing.T) {

    model := NewListViewModel()



    // Add test data

    conversations := []domain.Conversation{

        {ID: "1", Title: "Test Conversation 1"},

        {ID: "2", Title: "Test Conversation 2"},

    }



    items := make([]list.Item, len(conversations))

    for i, conv := range conversations {

        items[i] = conversationItem{Conversation: conv}

    }

    model.list.SetItems(items)



    // Test view output

    view := model.View()

    if view == "" {

        t.Error("View() returned empty string")

    }



    // Check that titles appear in view

    if !strings.Contains(view, "Test Conversation 1") {

        t.Error("View() doesn't contain first conversation title")

    }

}



// Integration test using teatest

func TestListViewModel_Integration(t *testing.T) {

    model := NewListViewModel()



    tm := teatest.NewTestModel(

        t, model,

        teatest.WithInitialTermSize(80, 24),

    )



    // Send key messages

    tm.Send(tea.KeyMsg{Type: tea.KeyDown})

    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})



    // Wait for final model

    fm := tm.FinalModel(t)



    // Assert final state

    finalModel := fm.(ListViewModel)

    _ = finalModel // Add assertions



    // Optionally compare against golden file

    teatest.RequireEqualOutput(t, tm.Output())

}

```



### Mocking File System



```go

// internal/testutil/fixtures.go

package testutil



import (

    "encoding/json"

    "testing"



    "github.com/spf13/afero"

    "github.com/yourusername/claudex/internal/domain"

)



// SetupTestFS creates an in-memory filesystem with test data

func SetupTestFS(t *testing.T) afero.Fs {

    fs := afero.NewMemMapFs()



    // Create test directory structure

    err := fs.MkdirAll("/data/conversations", 0755)

    if err != nil {

        t.Fatalf("Failed to create test directory: %v", err)

    }



    // Create test JSONL files

    conv1 := []domain.Message{

        {Role: "user", Content: "Hello"},

        {Role: "assistant", Content: "Hi there!"},

    }

    WriteJSONL(t, fs, "/data/conversations/conv1.jsonl", conv1)



    conv2 := []domain.Message{

        {Role: "user", Content: "How are you?"},

        {Role: "assistant", Content: "I'm doing well!"},

    }

    WriteJSONL(t, fs, "/data/conversations/conv2.jsonl", conv2)



    return fs

}



// WriteJSONL writes messages to a JSONL file

func WriteJSONL(t *testing.T, fs afero.Fs, path string, messages []domain.Message) {

    file, err := fs.Create(path)

    if err != nil {

        t.Fatalf("Failed to create file: %v", err)

    }

    defer file.Close()



    encoder := json.NewEncoder(file)

    for _, msg := range messages {

        if err := encoder.Encode(msg); err != nil {

            t.Fatalf("Failed to encode message: %v", err)

        }

    }

}

```



### Repository Testing with Mocks



```go

// internal/repository/jsonl_test.go

package repository



import (

    "context"

    "testing"



    "github.com/spf13/afero"

    "github.com/yourusername/claudex/internal/testutil"

)



func TestJSONLRepository_List(t *testing.T) {

    // Setup in-memory filesystem

    fs := testutil.SetupTestFS(t)



    // Create repository with mock filesystem

    repo := NewJSONLRepositoryWithFS(fs, "/data/conversations")



    // Test List

    conversations, err := repo.List(context.Background())

    if err != nil {

        t.Fatalf("List() error = %v", err)

    }



    if len(conversations) != 2 {

        t.Errorf("List() returned %d conversations, want 2", len(conversations))

    }

}



func TestJSONLRepository_GetByID(t *testing.T) {

    fs := testutil.SetupTestFS(t)

    repo := NewJSONLRepositoryWithFS(fs, "/data/conversations")



    conv, err := repo.GetByID(context.Background(), "conv1")

    if err != nil {

        t.Fatalf("GetByID() error = %v", err)

    }



    if len(conv.Messages) != 2 {

        t.Errorf("GetByID() returned %d messages, want 2", len(conv.Messages))

    }



    if conv.Messages[0].Content != "Hello" {

        t.Errorf("GetByID() first message = %q, want %q", conv.Messages[0].Content, "Hello")

    }

}

```



### Mock Repository for UI Testing



```go

// internal/testutil/mock_repo.go

package testutil



import (

    "context"

    "github.com/yourusername/claudex/internal/domain"

)



type MockRepository struct {

    ListFunc         func(ctx context.Context) ([]domain.Conversation, error)

    GetByIDFunc      func(ctx context.Context, id string) (*domain.Conversation, error)

    SearchFunc       func(ctx context.Context, query string) ([]domain.Conversation, error)

    GetStatisticsFunc func(ctx context.Context) (*domain.Statistics, error)

}



func (m *MockRepository) List(ctx context.Context) ([]domain.Conversation, error) {

    if m.ListFunc != nil {

        return m.ListFunc(ctx)

    }

    return []domain.Conversation{}, nil

}



func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {

    if m.GetByIDFunc != nil {

        return m.GetByIDFunc(ctx, id)

    }

    return &domain.Conversation{}, nil

}



func (m *MockRepository) Search(ctx context.Context, query string) ([]domain.Conversation, error) {

    if m.SearchFunc != nil {

        return m.SearchFunc(ctx, query)

    }

    return []domain.Conversation{}, nil

}



func (m *MockRepository) GetStatistics(ctx context.Context) (*domain.Statistics, error) {

    if m.GetStatisticsFunc != nil {

        return m.GetStatisticsFunc(ctx)

    }

    return &domain.Statistics{}, nil

}



// NewMockRepoWithData creates a mock repository with test data

func NewMockRepoWithData() *MockRepository {

    conversations := []domain.Conversation{

        {ID: "1", Title: "Test 1", Messages: []domain.Message{{Content: "Hello"}}},

        {ID: "2", Title: "Test 2", Messages: []domain.Message{{Content: "World"}}},

    }



    return &MockRepository{

        ListFunc: func(ctx context.Context) ([]domain.Conversation, error) {

            return conversations, nil

        },

        GetByIDFunc: func(ctx context.Context, id string) (*domain.Conversation, error) {

            for _, conv := range conversations {

                if conv.ID == id {

                    return &conv, nil

                }

            }

            return nil, fmt.Errorf("conversation not found")

        },

    }

}

```



### Table-Driven Tests



```go

// internal/search/query_test.go

package search



import (

    "testing"

)



func TestParser_Parse(t *testing.T) {

    tests := []struct {

        name    string

        input   string

        want    string

        wantErr bool

    }{

        {

            name:  "simple term",

            input: "hello",

            want:  "hello",

        },

        {

            name:  "AND operator",

            input: "hello AND world",

            want:  "(hello AND world)",

        },

        {

            name:  "OR operator",

            input: "hello OR world",

            want:  "(hello OR world)",

        },

        {

            name:  "NOT operator",

            input: "NOT hello",

            want:  "NOT hello",

        },

        {

            name:  "complex expression",

            input: "(hello OR world) AND NOT goodbye",

            want:  "((hello OR world) AND NOT goodbye)",

        },

        {

            name:  "field search",

            input: "title:hello",

            want:  "title:hello",

        },

        {

            name:  "implicit AND",

            input: "hello world",

            want:  "(hello AND world)",

        },

    }



    for _, tt := range tests {

        t.Run(tt.name, func(t *testing.T) {

            parser := newParser(tt.input)

            query, err := parser.Parse()



            if (err != nil) != tt.wantErr {

                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)

                return

            }



            if err == nil && query.Root.String() != tt.want {

                t.Errorf("Parse() = %v, want %v", query.Root.String(), tt.want)

            }

        })

    }

}

```



---



## 6. Error Handling



### Error Hierarchy



```go

// internal/domain/errors.go

package domain



import "fmt"



// Error types

type ErrorType int



const (

    ErrorTypeValidation ErrorType = iota

    ErrorTypeNotFound

    ErrorTypeIO

    ErrorTypeParsing

    ErrorTypeInternal

)



// AppError represents application-specific errors

type AppError struct {

    Type    ErrorType

    Err     error

    Message string // User-friendly message

    Context map[string]interface{}

}



func (e *AppError) Error() string {

    if e.Err != nil {

        return fmt.Sprintf("%s: %v", e.Message, e.Err)

    }

    return e.Message

}



func (e *AppError) Unwrap() error {

    return e.Err

}



func (e *AppError) UserMessage() string {

    return e.Message

}



// Constructor functions

func NewValidationError(msg string, err error) *AppError {

    return &AppError{

        Type:    ErrorTypeValidation,

        Err:     err,

        Message: msg,

    }

}



func NewNotFoundError(resource string) *AppError {

    return &AppError{

        Type:    ErrorTypeNotFound,

        Message: fmt.Sprintf("%s not found", resource),

    }

}



func NewIOError(msg string, err error) *AppError {

    return &AppError{

        Type:    ErrorTypeIO,

        Err:     err,

        Message: msg,

        Context: make(map[string]interface{}),

    }

}



func NewParsingError(file string, line int, err error) *AppError {

    return &AppError{

        Type:    ErrorTypeParsing,

        Err:     err,

        Message: fmt.Sprintf("Failed to parse file %s", file),

        Context: map[string]interface{}{

            "file": file,

            "line": line,

        },

    }

}

```



### Error Handling in Repository



```go

// internal/repository/jsonl.go (error handling additions)

package repository



import (

    "context"

    "fmt"

    "os"



    "github.com/yourusername/claudex/internal/domain"

)



func (r *JSONLRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {

    filename := filepath.Join(r.dataDir, id+".jsonl")



    // Check file existence first

    if _, err := os.Stat(filename); err != nil {

        if os.IsNotExist(err) {

            return nil, domain.NewNotFoundError("conversation")

        }

        return nil, domain.NewIOError("failed to access file", err)

    }



    conv, err := r.parseConversation(ctx, filename)

    if err != nil {

        // Wrap with context

        if appErr, ok := err.(*domain.AppError); ok {

            appErr.Context["conversation_id"] = id

            return nil, appErr

        }

        return nil, domain.NewIOError("failed to load conversation", err)

    }



    return conv, nil

}



func (r *JSONLRepository) parseConversation(ctx context.Context, filename string) (*domain.Conversation, error) {

    file, err := os.Open(filename)

    if err != nil {

        return nil, domain.NewIOError(

            fmt.Sprintf("failed to open file: %s", filename),

            err,

        )

    }

    defer file.Close()



    var conv domain.Conversation

    scanner := bufio.NewScanner(file)



    buf := make([]byte, 0, 64*1024)

    scanner.Buffer(buf, 1024*1024)



    lineNum := 0

    for scanner.Scan() {

        lineNum++



        select {

        case <-ctx.Done():

            return nil, ctx.Err()

        default:

        }



        var msg domain.Message

        if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {

            return nil, domain.NewParsingError(filename, lineNum, err)

        }



        conv.Messages = append(conv.Messages, msg)

    }



    if err := scanner.Err(); err != nil {

        return nil, domain.NewIOError("scanner error", err)

    }



    if len(conv.Messages) == 0 {

        return nil, domain.NewValidationError("conversation has no messages", nil)

    }



    return &conv, nil

}

```



### Error Handling in UI



```go

// internal/ui/error.go

package ui



import (

    "fmt"

    "github.com/charmbracelet/lipgloss"

    "github.com/yourusername/claudex/internal/domain"

)



// ErrorView renders user-friendly error messages

func ErrorView(err error) string {

    if err == nil {

        return ""

    }



    var style = lipgloss.NewStyle().

        Foreground(lipgloss.Color("196")). // Red

        Border(lipgloss.RoundedBorder()).

        BorderForeground(lipgloss.Color("196")).

        Padding(1, 2)



    // Check if it's an AppError

    if appErr, ok := err.(*domain.AppError); ok {

        return style.Render(renderAppError(appErr))

    }



    // Generic error

    return style.Render(fmt.Sprintf("Error: %v", err))

}



func renderAppError(err *domain.AppError) string {

    var msg string



    switch err.Type {

    case domain.ErrorTypeNotFound:

        msg = fmt.Sprintf("❌ %s\n\nThe requested resource was not found.", err.Message)



    case domain.ErrorTypeIO:

        msg = fmt.Sprintf("❌ %s\n\nThere was a problem accessing the file system.", err.Message)

        if file, ok := err.Context["file"].(string); ok {

            msg += fmt.Sprintf("\nFile: %s", file)

        }



    case domain.ErrorTypeParsing:

        msg = fmt.Sprintf("❌ %s\n\nThe file format appears to be invalid.", err.Message)

        if file, ok := err.Context["file"].(string); ok {

            msg += fmt.Sprintf("\nFile: %s", file)

        }

        if line, ok := err.Context["line"].(int); ok {

            msg += fmt.Sprintf("\nLine: %d", line)

        }



    case domain.ErrorTypeValidation:

        msg = fmt.Sprintf("❌ %s", err.Message)



    default:

        msg = fmt.Sprintf("❌ An unexpected error occurred\n\n%s", err.Message)

    }



    return msg

}



// Graceful degradation in Update

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case errMsg:

        // Log the full error for debugging

        if m.logger != nil {

            m.logger.Error("error occurred", "error", msg.error)

        }



        // Show user-friendly message

        m.errorMessage = msg.Error()



        // Attempt recovery based on error type

        if appErr, ok := msg.error.(*domain.AppError); ok {

            switch appErr.Type {

            case domain.ErrorTypeNotFound:

                // Return to list view

                return m.changeView(ViewList, nil)



            case domain.ErrorTypeIO, domain.ErrorTypeParsing:

                // Stay on current view but show error

                // User can retry or navigate away

                return m, nil

            }

        }



        return m, nil

    }



    return m.updateCurrentView(msg)

}

```



### Retry Logic with Exponential Backoff



```go

// internal/repository/retry.go

package repository



import (

    "context"

    "errors"

    "math"

    "time"

)



// RetryConfig configures retry behavior

type RetryConfig struct {

    MaxAttempts int

    InitialDelay time.Duration

    MaxDelay     time.Duration

    Multiplier   float64

}



var DefaultRetryConfig = RetryConfig{

    MaxAttempts:  3,

    InitialDelay: 100 * time.Millisecond,

    MaxDelay:     5 * time.Second,

    Multiplier:   2.0,

}



// WithRetry wraps a function with retry logic

func WithRetry[T any](ctx context.Context, config RetryConfig, fn func() (T, error)) (T, error) {

    var result T

    var err error



    for attempt := 0; attempt < config.MaxAttempts; attempt++ {

        result, err = fn()



        if err == nil {

            return result, nil

        }



        // Check if error is retryable

        if !isRetryable(err) {

            return result, err

        }



        // Check context cancellation

        if ctx.Err() != nil {

            return result, ctx.Err()

        }



        // Don't sleep on last attempt

        if attempt < config.MaxAttempts-1 {

            delay := calculateDelay(attempt, config)

            select {

            case <-time.After(delay):

            case <-ctx.Done():

                return result, ctx.Err()

            }

        }

    }



    return result, err

}



func calculateDelay(attempt int, config RetryConfig) time.Duration {

    delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))

    if delay > float64(config.MaxDelay) {

        delay = float64(config.MaxDelay)

    }

    return time.Duration(delay)

}



func isRetryable(err error) bool {

    // Determine if error is temporary

    var appErr *domain.AppError

    if errors.As(err, &appErr) {

        switch appErr.Type {

        case domain.ErrorTypeIO:

            return true // IO errors might be temporary

        case domain.ErrorTypeNotFound, domain.ErrorTypeValidation:

            return false // These are permanent

        }

    }

    return false

}

```



---



## 7. Performance Optimization



### Virtual Scrolling for Large Lists



```go

// internal/ui/components/virtuallist.go

package components



import (

    "github.com/charmbracelet/bubbles/viewport"

    tea "github.com/charmbracelet/bubbletea"

    "github.com/charmbracelet/lipgloss"

)



// VirtualList renders only visible items for performance

type VirtualList struct {

    viewport    viewport.Model

    items       []Item

    itemHeight  int

    cursor      int

    bufferSize  int // Items to render above/below viewport

}



type Item interface {

    Render(width int, selected bool) string

    Height() int

}



func NewVirtualList(width, height int) VirtualList {

    vp := viewport.New(width, height)



    return VirtualList{

        viewport:   vp,

        items:      []Item{},

        itemHeight: 1,

        bufferSize: 5, // Render 5 extra items above/below

    }

}



func (v VirtualList) Update(msg tea.Msg) (VirtualList, tea.Cmd) {

    var cmd tea.Cmd



    switch msg := msg.(type) {

    case tea.KeyMsg:

        switch msg.String() {

        case "up", "k":

            v.cursor--

            if v.cursor < 0 {

                v.cursor = 0

            }

            v.scrollToCursor()



        case "down", "j":

            v.cursor++

            if v.cursor >= len(v.items) {

                v.cursor = len(v.items) - 1

            }

            v.scrollToCursor()



        case "pgup":

            v.cursor -= v.viewport.Height

            if v.cursor < 0 {

                v.cursor = 0

            }

            v.scrollToCursor()



        case "pgdown":

            v.cursor += v.viewport.Height

            if v.cursor >= len(v.items) {

                v.cursor = len(v.items) - 1

            }

            v.scrollToCursor()

        }

    }



    v.viewport, cmd = v.viewport.Update(msg)

    return v, cmd

}



func (v VirtualList) View() string {

    // Calculate visible range

    firstVisible := v.viewport.YOffset / v.itemHeight

    lastVisible := (v.viewport.YOffset + v.viewport.Height) / v.itemHeight



    // Add buffer

    start := max(0, firstVisible-v.bufferSize)

    end := min(len(v.items), lastVisible+v.bufferSize)



    // Render only visible items

    var content string



    // Add padding for items above viewport

    if start > 0 {

        paddingHeight := start * v.itemHeight

        content += lipgloss.NewStyle().

            Height(paddingHeight).

            Render("")

    }



    // Render visible items

    for i := start; i < end; i++ {

        selected := i == v.cursor

        content += v.items[i].Render(v.viewport.Width, selected) + "\n"

    }



    // Add padding for items below viewport

    if end < len(v.items) {

        remaining := len(v.items) - end

        paddingHeight := remaining * v.itemHeight

        content += lipgloss.NewStyle().

            Height(paddingHeight).

            Render("")

    }



    v.viewport.SetContent(content)

    return v.viewport.View()

}



func (v *VirtualList) scrollToCursor() {

    cursorY := v.cursor * v.itemHeight



    // Scroll down if cursor below viewport

    if cursorY >= v.viewport.YOffset+v.viewport.Height {

        v.viewport.YOffset = cursorY - v.viewport.Height + v.itemHeight

    }



    // Scroll up if cursor above viewport

    if cursorY < v.viewport.YOffset {

        v.viewport.YOffset = cursorY

    }

}



func (v *VirtualList) SetItems(items []Item) {

    v.items = items

    if v.cursor >= len(items) {

        v.cursor = len(items) - 1

    }

    if v.cursor < 0 {

        v.cursor = 0

    }

}



func max(a, b int) int {

    if a > b {

        return a

    }

    return b

}



func min(a, b int) int {

    if a < b {

        return a

    }

    return b

}

```



### Lazy Loading with Pagination



```go

// internal/repository/paginated.go

package repository



import (

    "context"

    "github.com/yourusername/claudex/internal/domain"

)



// Page represents a page of results

type Page struct {

    Items      []domain.Conversation

    TotalCount int

    Page       int

    PageSize   int

    HasMore    bool

}



// PaginatedRepository supports pagination for large datasets

type PaginatedRepository interface {

    ListPage(ctx context.Context, page, pageSize int) (*Page, error)

}



type paginatedJSONLRepo struct {

    *JSONLRepository

}



func NewPaginatedRepository(repo *JSONLRepository) PaginatedRepository {

    return &paginatedJSONLRepo{JSONLRepository: repo}

}



func (r *paginatedJSONLRepo) ListPage(ctx context.Context, page, pageSize int) (*Page, error) {

    // Get all conversations (in real app, this would be optimized)

    allConvs, err := r.List(ctx)

    if err != nil {

        return nil, err

    }



    totalCount := len(allConvs)

    start := page * pageSize

    end := start + pageSize



    if start >= totalCount {

        return &Page{

            Items:      []domain.Conversation{},

            TotalCount: totalCount,

            Page:       page,

            PageSize:   pageSize,

            HasMore:    false,

        }, nil

    }



    if end > totalCount {

        end = totalCount

    }



    return &Page{

        Items:      allConvs[start:end],

        TotalCount: totalCount,

        Page:       page,

        PageSize:   pageSize,

        HasMore:    end < totalCount,

    }, nil

}



// UI Integration

// internal/ui/views/list.go (pagination additions)

func (m ListViewModel) loadPage(page int) tea.Cmd {

    return func() tea.Msg {

        ctx := context.Background()

        result, err := m.paginatedRepo.ListPage(ctx, page, m.pageSize)

        if err != nil {

            return errMsg{err}

        }

        return pageLoadedMsg{result}

    }

}



type pageLoadedMsg struct {

    page *Page

}



func (m ListViewModel) Update(msg tea.Msg) (ListViewModel, tea.Cmd) {

    switch msg := msg.(type) {

    case pageLoadedMsg:

        m.currentPage = msg.page

        // Update list items

        items := make([]list.Item, len(msg.page.Items))

        for i, conv := range msg.page.Items {

            items[i] = conversationItem{Conversation: conv}

        }

        return m, m.list.SetItems(items)



    case tea.KeyMsg:

        switch msg.String() {

        case "n": // Next page

            if m.currentPage.HasMore {

                return m, m.loadPage(m.currentPage.Page + 1)

            }

        case "p": // Previous page

            if m.currentPage.Page > 0 {

                return m, m.loadPage(m.currentPage.Page - 1)

            }

        }

    }



    var cmd tea.Cmd

    m.list, cmd = m.list.Update(msg)

    return m, cmd

}

```



### Memory Profiling



```go

// cmd/claudex/main.go

package main



import (

    "flag"

    "log"

    "os"

    "runtime"

    "runtime/pprof"



    tea "github.com/charmbracelet/bubbletea"

    "github.com/yourusername/claudex/internal/app"

)



var (

    cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

    memprofile = flag.String("memprofile", "", "write memory profile to file")

)



func main() {

    flag.Parse()



    // CPU profiling

    if *cpuprofile != "" {

        f, err := os.Create(*cpuprofile)

        if err != nil {

            log.Fatal(err)

        }

        defer f.Close()



        if err := pprof.StartCPUProfile(f); err != nil {

            log.Fatal(err)

        }

        defer pprof.StopCPUProfile()

    }



    // Run application

    model := app.InitialModel()

    p := tea.NewProgram(model, tea.WithAltScreen())



    if _, err := p.Run(); err != nil {

        log.Fatal(err)

    }



    // Memory profiling

    if *memprofile != "" {

        f, err := os.Create(*memprofile)

        if err != nil {

            log.Fatal(err)

        }

        defer f.Close()



        runtime.GC() // Get up-to-date statistics

        if err := pprof.WriteHeapProfile(f); err != nil {

            log.Fatal(err)

        }

    }

}



// Run profiling:

// go run cmd/claudex/main.go -cpuprofile=cpu.prof -memprofile=mem.prof

//

// Analyze results:

// go tool pprof -http=:8080 cpu.prof

// go tool pprof -http=:8080 mem.prof

```



### Optimization Techniques



```go

// internal/domain/conversation.go (memory optimizations)

package domain



import "strings"



type Conversation struct {

    ID        string

    Title     string

    Model     string

    CreatedAt time.Time

    UpdatedAt time.Time

    Messages  []Message



    // Cached computations (lazy evaluation)

    summary       string

    summaryOnce   sync.Once

    messageCount  int

    countOnce     sync.Once

}



// Summary returns a cached summary

func (c *Conversation) Summary() string {

    c.summaryOnce.Do(func() {

        c.summary = c.computeSummary()

    })

    return c.summary

}



func (c *Conversation) computeSummary() string {

    if len(c.Messages) == 0 {

        return ""

    }



    // Use strings.Builder for efficient string concatenation

    var sb strings.Builder

    sb.WriteString(c.Model)

    sb.WriteString(" • ")

    sb.WriteString(fmt.Sprintf("%d messages", len(c.Messages)))



    return sb.String()

}



// MessageCount returns cached count

func (c *Conversation) MessageCount() int {

    c.countOnce.Do(func() {

        c.messageCount = len(c.Messages)

    })

    return c.messageCount

}



// Release clears cached data to free memory

func (c *Conversation) Release() {

    c.Messages = nil // Allow GC to collect messages

}

```



### String Interning for Memory Savings



```go

// internal/cache/intern.go

package cache



import "sync"



// StringInterner reduces memory usage by reusing common strings

type StringInterner struct {

    mu      sync.RWMutex

    strings map[string]string

}



func NewStringInterner() *StringInterner {

    return &StringInterner{

        strings: make(map[string]string),

    }

}



func (si *StringInterner) Intern(s string) string {

    si.mu.RLock()

    if interned, ok := si.strings[s]; ok {

        si.mu.RUnlock()

        return interned

    }

    si.mu.RUnlock()



    si.mu.Lock()

    defer si.mu.Unlock()



    // Double-check after acquiring write lock

    if interned, ok := si.strings[s]; ok {

        return interned

    }



    si.strings[s] = s

    return s

}



// Usage in repository

func (r *JSONLRepository) parseConversation(filename string) (*domain.Conversation, error) {

    // ... parsing code ...



    // Intern common strings to save memory

    for i := range conv.Messages {

        conv.Messages[i].Role = r.interner.Intern(conv.Messages[i].Role)

        conv.Messages[i].Model = r.interner.Intern(conv.Messages[i].Model)

    }



    return &conv, nil

}

```



---



## 8. Complete Example



### Application Bootstrap



```go

// cmd/claudex/main.go

package main



import (

    "log"

    "os"

    "path/filepath"



    tea "github.com/charmbracelet/bubbletea"

    "github.com/yourusername/claudex/internal/app"

)



func main() {

    // Get data directory

    homeDir, err := os.UserHomeDir()

    if err != nil {

        log.Fatal(err)

    }



    dataDir := filepath.Join(homeDir, ".claude", "conversations")



    // Initialize application

    cfg := app.Config{

        DataDir:   dataDir,

        CacheSize: 100,

        PageSize:  50,

    }



    model, err := app.New(cfg)

    if err != nil {

        log.Fatal(err)

    }



    // Run Bubble Tea program

    p := tea.NewProgram(

        model,

        tea.WithAltScreen(),

        tea.WithMouseCellMotion(),

    )



    if _, err := p.Run(); err != nil {

        log.Fatal(err)

    }

}

```



### Application Setup with Dependency Injection



```go

// internal/app/app.go

package app



import (

    "github.com/yourusername/claudex/internal/cache"

    "github.com/yourusername/claudex/internal/repository"

    "github.com/yourusername/claudex/internal/search"

    "github.com/yourusername/claudex/internal/ui"

)



type Config struct {

    DataDir   string

    CacheSize int

    PageSize  int

}



func New(cfg Config) (ui.Model, error) {

    // Initialize string interner

    interner := cache.NewStringInterner()



    // Initialize repository

    baseRepo := repository.NewJSONLRepository(cfg.DataDir)

    baseRepo.SetInterner(interner)



    // Wrap with cache

    cachedRepo := repository.NewCachedRepository(baseRepo, cfg.CacheSize)



    // Initialize search engine

    searchEngine := search.NewEngine()



    // Initialize UI model with dependencies

    model := ui.NewModel(ui.ModelConfig{

        Repository:   cachedRepo,

        SearchEngine: searchEngine,

        PageSize:     cfg.PageSize,

    })



    return model, nil

}

```



### Complete Model Implementation



```go

// internal/ui/model.go (complete version)

package ui



import (

    tea "github.com/charmbracelet/bubbletea"

    "github.com/yourusername/claudex/internal/domain"

    "github.com/yourusername/claudex/internal/repository"

    "github.com/yourusername/claudex/internal/search"

    "github.com/yourusername/claudex/internal/ui/views"

)



type ModelConfig struct {

    Repository   repository.ConversationRepository

    SearchEngine search.Engine

    PageSize     int

}



type Model struct {

    // Configuration

    config ModelConfig



    // State

    currentView  ViewState

    previousView ViewState



    // Views

    listView    views.ListViewModel

    detailView  views.DetailViewModel

    searchView  views.SearchViewModel

    statsView   views.StatsViewModel

    paletteView views.PaletteViewModel



    // Data

    conversations []domain.Conversation



    // UI State

    width  int

    height int

    ready  bool

    error  error

}



func NewModel(cfg ModelConfig) Model {

    return Model{

        config:      cfg,

        currentView: ViewList,

        listView:    views.NewListViewModel(cfg.Repository, cfg.PageSize),

        detailView:  views.NewDetailViewModel(cfg.Repository),

        searchView:  views.NewSearchViewModel(cfg.SearchEngine),

        statsView:   views.NewStatsViewModel(cfg.Repository),

        paletteView: views.NewPaletteViewModel(),

    }

}



func (m Model) Init() tea.Cmd {

    return tea.Batch(

        tea.EnterAltScreen,

        m.listView.Init(),

    )

}



func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:

        m.width = msg.Width

        m.height = msg.Height

        m.ready = true



        // Update all views with new size

        m.listView.SetSize(msg.Width, msg.Height)

        m.detailView.SetSize(msg.Width, msg.Height)

        m.searchView.SetSize(msg.Width, msg.Height)

        m.statsView.SetSize(msg.Width, msg.Height)

        m.paletteView.SetSize(msg.Width, msg.Height)



        return m, nil



    case tea.KeyMsg:

        // Global keybindings

        switch msg.String() {

        case "ctrl+c", "q":

            if m.currentView != ViewCommandPalette {

                return m, tea.Quit

            }

        case "ctrl+p":

            m.previousView = m.currentView

            m.currentView = ViewCommandPalette

            return m, nil

        case "esc":

            if m.currentView == ViewCommandPalette {

                m.currentView = m.previousView

                return m, nil

            }

        }



    case changeViewMsg:

        m.previousView = m.currentView

        m.currentView = msg.view

        return m.initializeView(msg.view, msg.context)



    case errMsg:

        m.error = msg.error

        return m, nil

    }



    // Delegate to current view

    return m.updateCurrentView(msg)

}



func (m Model) View() string {

    if !m.ready {

        return "Initializing..."

    }



    var view string

    switch m.currentView {

    case ViewList:

        view = m.listView.View()

    case ViewDetail:

        view = m.detailView.View()

    case ViewSearch:

        view = m.searchView.View()

    case ViewStatistics:

        view = m.statsView.View()

    case ViewCommandPalette:

        // Overlay on previous view

        var baseView string

        switch m.previousView {

        case ViewList:

            baseView = m.listView.View()

        case ViewDetail:

            baseView = m.detailView.View()

        case ViewSearch:

            baseView = m.searchView.View()

        case ViewStatistics:

            baseView = m.statsView.View()

        }

        return m.renderOverlay(baseView, m.paletteView.View())

    }



    // Show error if present

    if m.error != nil {

        view += "\n" + ErrorView(m.error)

    }



    return view

}



func (m Model) updateCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd



    switch m.currentView {

    case ViewList:

        m.listView, cmd = m.listView.Update(msg)

    case ViewDetail:

        m.detailView, cmd = m.detailView.Update(msg)

    case ViewSearch:

        m.searchView, cmd = m.searchView.Update(msg)

    case ViewStatistics:

        m.statsView, cmd = m.statsView.Update(msg)

    case ViewCommandPalette:

        m.paletteView, cmd = m.paletteView.Update(msg)

    }



    return m, cmd

}



func (m Model) initializeView(view ViewState, context interface{}) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd



    switch view {

    case ViewDetail:

        if id, ok := context.(string); ok {

            cmd = m.detailView.LoadConversation(id)

        }

    case ViewSearch:

        if query, ok := context.(string); ok {

            cmd = m.searchView.SetQuery(query)

        }

    }



    return m, cmd

}



func (m Model) renderOverlay(base, overlay string) string {

    // Center overlay on base view

    // Implementation depends on lipgloss for positioning

    return base + "\n" + overlay

}

```



---



## Key Takeaways



### Architecture Principles



1. **Separation of Concerns**: Keep UI, business logic, and data access separate

2. **Dependency Injection**: Pass dependencies through constructors

3. **Interface-Based Design**: Program to interfaces, not implementations

4. **Single Responsibility**: Each package/type should have one clear purpose

5. **Composition Over Inheritance**: Use embedding and delegation



### Performance Guidelines



1. **Lazy Loading**: Don't load all data upfront

2. **Caching**: Use LRU cache for frequently accessed data

3. **Virtual Scrolling**: Render only visible items

4. **String Interning**: Reuse common strings

5. **Profiling**: Use pprof to identify bottlenecks



### Testing Strategy



1. **Unit Tests**: Test each component in isolation

2. **Mocking**: Use interfaces and dependency injection for testability

3. **Table-Driven Tests**: Use for testing multiple scenarios

4. **Integration Tests**: Test component interactions

5. **Golden Files**: Use for snapshot testing UI output



### Error Handling



1. **Custom Error Types**: Create application-specific errors

2. **Error Wrapping**: Add context while preserving original error

3. **User-Friendly Messages**: Show helpful messages to users

4. **Graceful Degradation**: Allow app to continue after recoverable errors

5. **Retry Logic**: Implement exponential backoff for temporary failures



### State Management



1. **Single Source of Truth**: Root model holds canonical state

2. **Message-Driven Updates**: All state changes via messages

3. **View Delegation**: Delegate updates to child view models

4. **State Preservation**: Save/restore state during navigation

5. **Unidirectional Data Flow**: Follow The Elm Architecture



---



## Recommended Packages



### Core TUI Framework

- `github.com/charmbracelet/bubbletea` - TUI framework

- `github.com/charmbracelet/bubbles` - Reusable components

- `github.com/charmbracelet/lipgloss` - Styling and layout



### Testing

- `github.com/charmbracelet/x/exp/teatest` - Bubble Tea testing

- `github.com/spf13/afero` - Mock filesystem

- `github.com/stretchr/testify` - Assertion library



### Data & Caching

- `github.com/hashicorp/golang-lru/v2` - LRU cache

- Standard library `encoding/json` - JSON parsing



### Profiling & Debugging

- `runtime/pprof` - Built-in profiling

- `net/http/pprof` - HTTP profiling endpoint

- `github.com/pkg/profile` - Simple profiling



### Logging

- `log/slog` - Structured logging (Go 1.21+)

- `github.com/charmbracelet/log` - Styled logging



---



## Additional Resources



### Official Documentation

- [Bubble Tea GitHub](https://github.com/charmbracelet/bubbletea)

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/main/tutorials)

- [Go Project Layout](https://github.com/golang-standards/project-layout)



### Articles & Tutorials

- [Building TUI apps with Bubble Tea](https://charm.sh/blog/)

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

- [Profiling Go Programs](https://go.dev/blog/pprof)
