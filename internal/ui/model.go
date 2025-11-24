// ABOUTME: Root Bubble Tea model with global state management
package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/TrevorS/claudex/internal/domain"
	"github.com/TrevorS/claudex/internal/repository"
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

// Model is the root Bubble Tea model managing global application state.
type Model struct {
	// Core dependencies
	repository repository.Repository

	// View state
	state ViewState

	// Window dimensions
	width  int
	height int

	// Conversation list state
	conversations []*domain.Conversation
	selected      string // Selected conversation ID
	scrollPos     int    // Scroll position in list

	// Error state
	err error
}

// NewModel creates and returns a new root Model.
func NewModel(repo repository.Repository) Model {
	return Model{
		repository: repo,
		state:      ViewList,
		width:      0,
		height:     0,
	}
}

// Init returns the initial command for the model.
// This loads the initial conversation list from the repository.
func (m Model) Init() tea.Cmd {
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
