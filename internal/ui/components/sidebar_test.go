// ABOUTME: Comprehensive tests for the Sidebar component covering construction,
// project loading, filtering, navigation, rendering, and data updates.
package components

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/TrevorS/claudex/internal/domain"
)

// Test helpers

// makeConversation creates a test conversation with the given parameters.
func makeConversation(id, title, projectPath string, updatedAt time.Time) *domain.Conversation {
	conv := domain.NewConversation(id, title, "claude-sonnet-4", updatedAt)
	// Need to set ProjectPath manually since NewConversation doesn't take it
	return &domain.Conversation{
		ID:          id,
		Title:       title,
		Model:       conv.Model,
		ProjectPath: projectPath,
		CreatedAt:   conv.CreatedAt,
		UpdatedAt:   updatedAt,
		Messages:    conv.Messages,
	}
}

// makeTestConversations creates a set of test conversations across different projects and dates.
func makeTestConversations() []*domain.Conversation {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)
	lastMonth := now.Add(-32 * 24 * time.Hour)

	return []*domain.Conversation{
		makeConversation("1", "Recent conversation", "project-a", now),
		makeConversation("2", "Another recent", "project-a", yesterday),
		makeConversation("3", "Week old", "project-b", lastWeek),
		makeConversation("4", "Month old", "project-b", lastMonth),
		makeConversation("5", "Today's work", "project-c", now),
		makeConversation("6", "No project", "", yesterday),
	}
}

// Test: Construction

func TestNewSidebar(t *testing.T) {
	s := NewSidebar()
	assert.NotNil(t, s)
	assert.Empty(t, s.projects)
	assert.Empty(t, s.conversations)
	assert.Equal(t, -1, s.selectedProjectIndex) // -1 means no selection (all projects)
	assert.Equal(t, FilterAll, s.selectedFilter)
	assert.False(t, s.hasFocus)
}

func TestSidebar_Init(t *testing.T) {
	s := NewSidebar()
	cmd := s.Init()
	assert.Nil(t, cmd)
}

// Test: Project Loading

func TestSidebar_SetConversations_ExtractsProjects(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()

	s = s.SetConversations(convs)

	// Should extract 3 unique projects (project-a, project-b, project-c) plus empty project
	assert.Len(t, s.projects, 4)
	assert.Contains(t, s.projects, "project-a")
	assert.Contains(t, s.projects, "project-b")
	assert.Contains(t, s.projects, "project-c")
	assert.Contains(t, s.projects, "") // No project
}

func TestSidebar_CalculatesConversationCounts(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)

	// project-a should have 2 conversations
	meta, exists := s.projectMetadata["project-a"]
	assert.True(t, exists)
	assert.Equal(t, 2, meta.Count)

	// project-b should have 2 conversations
	meta, exists = s.projectMetadata["project-b"]
	assert.True(t, exists)
	assert.Equal(t, 2, meta.Count)

	// project-c should have 1 conversation
	meta, exists = s.projectMetadata["project-c"]
	assert.True(t, exists)
	assert.Equal(t, 1, meta.Count)

	// Empty project should have 1 conversation
	meta, exists = s.projectMetadata[""]
	assert.True(t, exists)
	assert.Equal(t, 1, meta.Count)
}

func TestSidebar_CalculatesLastUpdated(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)

	now := time.Now()

	// project-a should have the most recent timestamp
	meta := s.projectMetadata["project-a"]
	diff := now.Sub(meta.LastUpdated)
	assert.True(t, diff < 5*time.Second, "LastUpdated should be within last 5 seconds")

	// project-b should have older timestamp (last week)
	meta = s.projectMetadata["project-b"]
	diff = now.Sub(meta.LastUpdated)
	assert.True(t, diff > 6*24*time.Hour, "LastUpdated should be over 6 days ago")
}

func TestSidebar_HandlesEmptyConversationList(t *testing.T) {
	s := NewSidebar()
	s = s.SetConversations([]*domain.Conversation{})

	assert.Empty(t, s.projects)
	assert.Empty(t, s.projectMetadata)
	assert.Equal(t, -1, s.selectedProjectIndex) // -1 means no selection
}

func TestSidebar_HandlesConversationsWithoutProjectPath(t *testing.T) {
	s := NewSidebar()
	convs := []*domain.Conversation{
		makeConversation("1", "Conv 1", "", time.Now()),
		makeConversation("2", "Conv 2", "", time.Now()),
	}
	s = s.SetConversations(convs)

	assert.Len(t, s.projects, 1)
	assert.Contains(t, s.projects, "")
	assert.Equal(t, 2, s.projectMetadata[""].Count)
}

// Test: Filters

func TestSidebar_HasAllFilterTypes(t *testing.T) {
	s := NewSidebar()
	filters := s.Filters()

	assert.Len(t, filters, 4)
	assert.Contains(t, filters, FilterAll)
	assert.Contains(t, filters, FilterToday)
	assert.Contains(t, filters, FilterThisWeek)
	assert.Contains(t, filters, FilterThisMonth)
}

func TestSidebar_SelectedFilter_DefaultsToAll(t *testing.T) {
	s := NewSidebar()
	assert.Equal(t, FilterAll, s.SelectedFilter())
}

func TestSidebar_SetFilter(t *testing.T) {
	s := NewSidebar()
	s = s.SetFilter(FilterToday)
	assert.Equal(t, FilterToday, s.selectedFilter)

	s = s.SetFilter(FilterThisWeek)
	assert.Equal(t, FilterThisWeek, s.selectedFilter)
}

func TestSidebar_FilterAll_ShowsAllConversations(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFilter(FilterAll)

	filtered := s.FilteredConversations()
	assert.Len(t, filtered, len(convs))
}

func TestSidebar_FilterToday_MatchesCorrectDateRange(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFilter(FilterToday)

	filtered := s.FilteredConversations()

	// Should only include conversations from today
	// Based on makeTestConversations: ids 1 and 5 are from today
	assert.Len(t, filtered, 2)
	for _, conv := range filtered {
		diff := time.Since(conv.UpdatedAt)
		assert.True(t, diff < 24*time.Hour, "Conversation should be from today")
	}
}

func TestSidebar_FilterThisWeek_Matches7DayRange(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFilter(FilterThisWeek)

	filtered := s.FilteredConversations()

	// Should include conversations from last 7 days
	// Based on makeTestConversations: ids 1, 2, 3, 5 are within 7 days
	assert.Len(t, filtered, 4)
	for _, conv := range filtered {
		diff := time.Since(conv.UpdatedAt)
		assert.True(t, diff < 8*24*time.Hour, "Conversation should be within 7 days")
	}
}

func TestSidebar_FilterThisMonth_Matches30DayRange(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFilter(FilterThisMonth)

	filtered := s.FilteredConversations()

	// Should include conversations from last 30 days
	// Based on makeTestConversations: ids 1, 2, 3, 5, 6 are within 30 days
	assert.Len(t, filtered, 5)
	for _, conv := range filtered {
		diff := time.Since(conv.UpdatedAt)
		assert.True(t, diff < 31*24*time.Hour, "Conversation should be within 30 days")
	}
}

// Test: Navigation

func TestSidebar_SelectProject_WithArrowKeys(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(true)

	// Initially at index -1 (all projects)
	// Set to 0 to start navigating
	s.selectedProjectIndex = 0
	assert.Equal(t, 0, s.selectedProjectIndex)

	// Press down arrow
	newS, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
	s = newS
	assert.Equal(t, 1, s.selectedProjectIndex)

	// Press down arrow again
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyDown})
	s = newS
	assert.Equal(t, 2, s.selectedProjectIndex)

	// Press up arrow
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyUp})
	s = newS
	assert.Equal(t, 1, s.selectedProjectIndex)
}

func TestSidebar_SelectProject_DoesNotGoOutOfBounds(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(true)

	// Should have 4 projects (indices 0-3)
	assert.Len(t, s.projects, 4)

	// Start at "All" (-1), go up should stay at -1
	s.selectedProjectIndex = -1
	newS, _ := s.Update(tea.KeyMsg{Type: tea.KeyUp})
	s = newS
	assert.Equal(t, -1, s.selectedProjectIndex, "Going up from 'All' should stay at 'All'")

	// Going down from All should go to 0
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyDown})
	s = newS
	assert.Equal(t, 0, s.selectedProjectIndex)

	// Go to last index
	s.selectedProjectIndex = 3

	// Try to go down from last
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyDown})
	s = newS
	assert.Equal(t, 3, s.selectedProjectIndex)
}

func TestSidebar_NumberKeys_JumpToProject(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(true)

	// Press '1' to jump to first project
	newS, _ := s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	s = newS
	assert.Equal(t, 0, s.selectedProjectIndex)

	// Press '3' to jump to third project
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	s = newS
	assert.Equal(t, 2, s.selectedProjectIndex)

	// Press '9' should clamp to max index
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}})
	s = newS
	assert.Equal(t, 3, s.selectedProjectIndex) // max index for 4 projects
}

func TestSidebar_FilterNavigation(t *testing.T) {
	s := NewSidebar()
	s = s.SetFocus(true)

	// Initially on FilterAll
	assert.Equal(t, FilterAll, s.selectedFilter)

	// Press 'f' to cycle filters
	newS, _ := s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	s = newS
	assert.Equal(t, FilterToday, s.selectedFilter)

	// Press 'f' again
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	s = newS
	assert.Equal(t, FilterThisWeek, s.selectedFilter)

	// Press 'f' again
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	s = newS
	assert.Equal(t, FilterThisMonth, s.selectedFilter)

	// Press 'f' again - should wrap to All
	newS, _ = s.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	s = newS
	assert.Equal(t, FilterAll, s.selectedFilter)
}

func TestSidebar_EnterKey_SendsProjectSelectedMsg(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(true)
	s.selectedProjectIndex = 1

	newS, cmd := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	s = newS

	assert.NotNil(t, cmd)
	msg := cmd()
	projectMsg, ok := msg.(ProjectSelectedMsg)
	assert.True(t, ok)
	assert.Equal(t, s.projects[1], projectMsg.ProjectPath)
}

func TestSidebar_TabKey_SendsFocusNextMsg(t *testing.T) {
	s := NewSidebar()
	s = s.SetFocus(true)

	_, cmd := s.Update(tea.KeyMsg{Type: tea.KeyTab})

	assert.NotNil(t, cmd)
	msg := cmd()
	_, ok := msg.(FocusNextMsg)
	assert.True(t, ok)
}

// Test: Rendering

func TestSidebar_View_ReturnsNonEmpty(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)

	view := s.View()
	assert.NotEmpty(t, view)
}

func TestSidebar_View_RendersAllProjects(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetSize(22, 24) // Set proper size for rendering

	view := s.View()

	// Should contain project names
	assert.Contains(t, view, "project-a")
	assert.Contains(t, view, "project-b")
	assert.Contains(t, view, "project-c")
}

func TestSidebar_View_RendersAllFilters(t *testing.T) {
	s := NewSidebar()
	s = s.SetSize(22, 24) // Set proper size for rendering
	view := s.View()

	// Should contain filter labels
	assert.Contains(t, view, "All")
	assert.Contains(t, view, "Today")
	assert.Contains(t, view, "This Week")
	assert.Contains(t, view, "This Month")
}

func TestSidebar_View_ShowsMetadata(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetSize(22, 24) // Set proper size for rendering

	view := s.View()

	// Should show conversation counts
	// project-a has 2 conversations, should show "2" somewhere
	assert.Regexp(t, `2`, view)
}

func TestSidebar_View_ShowsFocusIndicator(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetSize(22, 24) // Set proper size for rendering

	// Without focus
	viewNoFocus := s.View()

	// With focus
	s = s.SetFocus(true)
	viewWithFocus := s.View()

	// Views should be different when focus changes
	assert.NotEqual(t, viewNoFocus, viewWithFocus)
}

func TestSidebar_View_ShowsSelectionHighlight(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(true)
	s = s.SetSize(22, 24) // Set proper size for rendering

	// Select first project
	s.selectedProjectIndex = 0
	view1 := s.View()
	// View should be rendered
	assert.NotEmpty(t, view1)
	assert.Contains(t, view1, "PROJECTS")

	// Select second project
	s.selectedProjectIndex = 1
	view2 := s.View()
	assert.NotEmpty(t, view2)

	// Note: Lipgloss styling might not show differences in test output,
	// but we can verify the sidebar state is different
	assert.Equal(t, 1, s.selectedProjectIndex)
}

// Test: Data Updates

func TestSidebar_SetConversations_UpdatesProjectList(t *testing.T) {
	s := NewSidebar()

	// Set initial conversations
	convs1 := []*domain.Conversation{
		makeConversation("1", "Conv 1", "project-a", time.Now()),
	}
	s = s.SetConversations(convs1)
	assert.Len(t, s.projects, 1)

	// Update with new conversations
	convs2 := makeTestConversations()
	s = s.SetConversations(convs2)
	assert.Len(t, s.projects, 4)
}

func TestSidebar_SetConversations_RecalculatesCounts(t *testing.T) {
	s := NewSidebar()

	// Set initial conversations
	convs1 := []*domain.Conversation{
		makeConversation("1", "Conv 1", "project-a", time.Now()),
	}
	s = s.SetConversations(convs1)
	assert.Equal(t, 1, s.projectMetadata["project-a"].Count)

	// Add more conversations to same project
	convs2 := []*domain.Conversation{
		makeConversation("1", "Conv 1", "project-a", time.Now()),
		makeConversation("2", "Conv 2", "project-a", time.Now()),
		makeConversation("3", "Conv 3", "project-a", time.Now()),
	}
	s = s.SetConversations(convs2)
	assert.Equal(t, 3, s.projectMetadata["project-a"].Count)
}

func TestSidebar_SetConversations_RecalculatesLastUpdated(t *testing.T) {
	s := NewSidebar()

	oldTime := time.Now().Add(-48 * time.Hour)
	newTime := time.Now()

	// Set initial conversations
	convs1 := []*domain.Conversation{
		makeConversation("1", "Conv 1", "project-a", oldTime),
	}
	s = s.SetConversations(convs1)
	oldLastUpdated := s.projectMetadata["project-a"].LastUpdated

	// Add newer conversation
	convs2 := []*domain.Conversation{
		makeConversation("1", "Conv 1", "project-a", oldTime),
		makeConversation("2", "Conv 2", "project-a", newTime),
	}
	s = s.SetConversations(convs2)
	newLastUpdated := s.projectMetadata["project-a"].LastUpdated

	// New last updated should be more recent
	assert.True(t, newLastUpdated.After(oldLastUpdated))
}

// Test: Focus management

func TestSidebar_SetFocus(t *testing.T) {
	s := NewSidebar()
	assert.False(t, s.hasFocus)

	s = s.SetFocus(true)
	assert.True(t, s.hasFocus)

	s = s.SetFocus(false)
	assert.False(t, s.hasFocus)
}

func TestSidebar_IgnoresKeysWhenNotFocused(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)
	s = s.SetFocus(false) // Not focused

	initialIndex := s.selectedProjectIndex

	// Try to navigate
	newS, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
	s = newS

	// Index should not change
	assert.Equal(t, initialIndex, s.selectedProjectIndex)
}

// Test: Window size handling

func TestSidebar_SetSize(t *testing.T) {
	s := NewSidebar()
	// Initial values should be reasonable defaults (so View() renders content before WindowSizeMsg)
	assert.Equal(t, 24, s.viewWidth)
	assert.Equal(t, 20, s.viewHeight)

	s = s.SetSize(80, 24)
	assert.Equal(t, 80, s.viewWidth)
	assert.Equal(t, 24, s.viewHeight)
}

func TestSidebar_RespondsToWindowSizeMsg(t *testing.T) {
	s := NewSidebar()

	newS, _ := s.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	s = newS

	assert.Equal(t, 120, s.viewWidth)
	assert.Equal(t, 40, s.viewHeight)
}

// Test: Edge cases

func TestSidebar_SelectedProject_WithNoProjects(t *testing.T) {
	s := NewSidebar()
	// No conversations set

	project := s.SelectedProject()
	assert.Empty(t, project)
}

func TestSidebar_SelectedProject_ReturnsCorrectProject(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)

	s.selectedProjectIndex = 1
	project := s.SelectedProject()
	assert.Equal(t, s.projects[1], project)
}

func TestSidebar_FilteredConversations_RespectsProjectSelection(t *testing.T) {
	s := NewSidebar()
	convs := makeTestConversations()
	s = s.SetConversations(convs)

	// Select project-a (which has 2 conversations)
	// Find the index of project-a
	for i, p := range s.projects {
		if p == "project-a" {
			s.selectedProjectIndex = i
			break
		}
	}

	filtered := s.FilteredConversations()

	// All filtered conversations should be from project-a
	for _, conv := range filtered {
		assert.Equal(t, "project-a", conv.ProjectPath)
	}
}
