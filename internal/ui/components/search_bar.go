// ABOUTME: search_bar.go implements a search bar component with debouncing,
// result display, and error handling for the search UI.
package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/TrevorS/claudex/internal/domain"
)

// SearchDebounceMsg is sent after a debounce delay to trigger search execution.
type SearchDebounceMsg struct {
	Query string
	Time  time.Time
}

// SearchResultsMsg contains search results.
type SearchResultsMsg struct {
	Results []*domain.Conversation
	Query   string
}

// SearchErrorMsg contains a search error.
type SearchErrorMsg struct {
	Error string
	Query string
}

// SearchBar represents a search input component with debouncing.
type SearchBar struct {
	input     textinput.Model
	lastQuery string
	queryTime time.Time
	debounce  time.Duration

	// Results state
	resultCount int
	hasResults  bool

	// Error state
	hasError bool
	errorMsg string

	// Styles
	normalStyle lipgloss.Style
	errorStyle  lipgloss.Style
}

// NewSearchBar creates a new search bar component.
func NewSearchBar() *SearchBar {
	input := textinput.New()
	input.Placeholder = "Search conversations... (try: python, title:chat, tokens:>5000)"
	input.Focus()
	input.CharLimit = 200
	input.Width = 50

	return &SearchBar{
		input:       input,
		debounce:    300 * time.Millisecond,
		normalStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		errorStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
	}
}

// Init initializes the search bar.
func (s *SearchBar) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles search bar updates.
func (s *SearchBar) Update(msg tea.Msg) (*SearchBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Immediate search on enter
			if s.input.Value() != "" {
				s.lastQuery = s.input.Value()
				s.queryTime = time.Now()
				return s, s.sendSearchDebounce()
			}

		case "esc":
			// Clear search
			if s.input.Value() != "" {
				s.input.SetValue("")
				s.lastQuery = ""
				s.hasError = false
				s.hasResults = false
				s.resultCount = 0
				return s, nil
			}
		}

	case SearchResultsMsg:
		// Update result count
		s.hasResults = true
		s.resultCount = len(msg.Results)
		s.hasError = false
		return s, nil

	case SearchErrorMsg:
		// Update error state
		s.hasError = true
		s.errorMsg = msg.Error
		s.hasResults = false
		return s, nil

	case SearchDebounceMsg:
		// If the query matches and hasn't changed, execute search
		if msg.Query == s.input.Value() && msg.Time.Equal(s.queryTime) {
			// The debounce period has passed, search should be executed upstream
			return s, nil
		}
	}

	// Update text input
	s.input, cmd = s.input.Update(msg)

	// Check if input changed - start debounce timer
	if s.input.Value() != s.lastQuery {
		s.queryTime = time.Now()
		return s, tea.Batch(cmd, s.debounceCmd())
	}

	return s, cmd
}

// View renders the search bar.
func (s *SearchBar) View() string {
	var status string

	if s.hasError {
		status = s.errorStyle.Render(fmt.Sprintf(" ❌ %s", s.errorMsg))
	} else if s.hasResults {
		status = s.normalStyle.Render(fmt.Sprintf(" ✓ %d results", s.resultCount))
	} else if s.input.Value() != "" {
		status = s.normalStyle.Render(" 🔍 Searching...")
	}

	return s.input.View() + status
}

// Value returns the current search query.
func (s *SearchBar) Value() string {
	return s.input.Value()
}

// SetValue sets the search query.
func (s *SearchBar) SetValue(value string) {
	s.input.SetValue(value)
}

// Focus focuses the search bar.
func (s *SearchBar) Focus() tea.Cmd {
	return s.input.Focus()
}

// Blur removes focus from the search bar.
func (s *SearchBar) Blur() {
	s.input.Blur()
}

// SetWidth sets the width of the search bar.
func (s *SearchBar) SetWidth(width int) {
	s.input.Width = width
}

// debounceCmd returns a command that waits for the debounce period.
func (s *SearchBar) debounceCmd() tea.Cmd {
	query := s.input.Value()
	t := s.queryTime
	return func() tea.Msg {
		time.Sleep(s.debounce)
		return SearchDebounceMsg{Query: query, Time: t}
	}
}

// sendSearchDebounce immediately sends a debounce message.
func (s *SearchBar) sendSearchDebounce() tea.Cmd {
	query := s.input.Value()
	t := s.queryTime
	return func() tea.Msg {
		return SearchDebounceMsg{Query: query, Time: t}
	}
}
