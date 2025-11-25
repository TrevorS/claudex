// ABOUTME: Root Bubble Tea model with global state management
package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
	"github.com/TrevorS/claudex/internal/search"
	"github.com/TrevorS/claudex/internal/ui/components"
	"github.com/TrevorS/claudex/internal/ui/views"
)

// ViewState represents the current view being displayed.
type ViewState int

const (
	ViewList ViewState = iota
	ViewDetail
	ViewSearch
	ViewStats
	ViewPalette
)

// FocusPane represents which pane has keyboard focus.
type FocusPane int

const (
	FocusSidebar FocusPane = iota
	FocusCenter
	FocusRight
)

// Model is the root Bubble Tea model managing global application state.
type Model struct {
	// Core dependencies
	repository   repository.Repository
	searchEngine *search.Engine

	// View state
	state     ViewState
	focusPane FocusPane

	// Window dimensions
	width  int
	height int

	// Sub-views
	sidebar    *components.Sidebar
	listView   *views.ListView
	detailView *views.DetailView
	searchBar  *components.SearchBar
	metadata   *components.Metadata

	// Conversation list state
	conversations  []*domain.Conversation
	searchResults  []*domain.Conversation
	selected       string // Selected conversation ID
	scrollPos      int    // Scroll position in list
	isSearchActive bool   // True when showing search results
	currentQuery   string // Current search query

	// UI state
	loading        bool // True until conversations loaded
	sidebarVisible bool // False when viewing conversation detail

	// Error state
	err error
}

// NewModel creates and returns a new root Model.
func NewModel(repo repository.Repository) Model {
	sidebar := components.NewSidebar()
	sidebar = sidebar.SetFocus(true) // Start focused since focusPane is FocusSidebar
	listView := views.NewListView(repo)
	searchEngine := search.NewEngine(repo)
	searchBar := components.NewSearchBar()
	metadata := components.NewMetadata()
	return Model{
		repository:     repo,
		searchEngine:   searchEngine,
		state:          ViewList,
		focusPane:      FocusSidebar,
		width:          0,
		height:         0,
		sidebar:        sidebar,
		listView:       listView,
		searchBar:      searchBar,
		metadata:       metadata,
		loading:        true, // Start in loading state
		sidebarVisible: true, // Visible by default
	}
}

// Init returns the initial command for the model.
// This loads the initial conversation list from the repository.
func (m Model) Init() tea.Cmd {
	if m.listView != nil {
		return m.listView.Init()
	}
	return tea.Batch()
}

// Repository returns the repository instance.
func (m Model) Repository() repository.Repository {
	return m.repository
}

// State returns the current view state.
func (m Model) State() ViewState {
	return m.state
}

// SetState sets the current view state.
func (m Model) SetState(state ViewState) Model {
	m.state = state
	return m
}

// FocusPane returns the current focused pane.
func (m Model) FocusPane() FocusPane {
	return m.focusPane
}

// SetFocusPane sets the current focused pane.
func (m Model) SetFocusPane(pane FocusPane) Model {
	m.focusPane = pane
	return m
}

// Width returns the terminal width.
func (m Model) Width() int {
	return m.width
}

// Height returns the terminal height.
func (m Model) Height() int {
	return m.height
}

// SetSize sets the terminal width and height.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Conversations returns the loaded conversations.
func (m Model) Conversations() []*domain.Conversation {
	return m.conversations
}

// SetConversations sets the conversations list.
func (m Model) SetConversations(conversations []*domain.Conversation) Model {
	m.conversations = conversations
	return m
}

// Selected returns the selected conversation ID.
func (m Model) Selected() string {
	return m.selected
}

// SetSelected sets the selected conversation ID.
func (m Model) SetSelected(id string) Model {
	m.selected = id
	return m
}

// ScrollPos returns the current scroll position.
func (m Model) ScrollPos() int {
	return m.scrollPos
}

// SetScrollPos sets the scroll position.
func (m Model) SetScrollPos(pos int) Model {
	m.scrollPos = pos
	return m
}

// Error returns any error state.
func (m Model) Error() error {
	return m.err
}

// SetError sets the error state.
func (m Model) SetError(err error) Model {
	m.err = err
	return m
}

// Sidebar returns the sidebar component.
func (m Model) Sidebar() *components.Sidebar {
	return m.sidebar
}

// SearchBar returns the search bar component.
func (m Model) SearchBar() *components.SearchBar {
	return m.searchBar
}

// Metadata returns the metadata component.
func (m Model) Metadata() *components.Metadata {
	return m.metadata
}

// SearchResults returns the current search results.
func (m Model) SearchResults() []*domain.Conversation {
	return m.searchResults
}

// SetSearchResults sets the search results.
func (m Model) SetSearchResults(results []*domain.Conversation) Model {
	m.searchResults = results
	return m
}

// IsSearchActive returns true if search is active.
func (m Model) IsSearchActive() bool {
	return m.isSearchActive
}

// SetSearchActive sets the search active flag.
func (m Model) SetSearchActive(active bool) Model {
	m.isSearchActive = active
	return m
}

// CurrentQuery returns the current search query.
func (m Model) CurrentQuery() string {
	return m.currentQuery
}

// SetCurrentQuery sets the current search query.
func (m Model) SetCurrentQuery(query string) Model {
	m.currentQuery = query
	return m
}
