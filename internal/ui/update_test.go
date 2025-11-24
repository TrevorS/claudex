package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/TrevorS/claudex/internal/repository"
)

func TestUpdateWindowSize(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Send a window size message
	msg := tea.WindowSizeMsg{Width: 120, Height: 30}
	newModel, _ := model.Update(msg)

	updatedModel := newModel.(Model)
	assert.Equal(t, 120, updatedModel.Width())
	assert.Equal(t, 30, updatedModel.Height())
}

func TestUpdateQuitKeyPress(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Send Ctrl+C
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(msg)

	assert.NotNil(t, cmd)
}

func TestUpdateEscapeKey(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetState(ViewDetail)

	// Send Escape key
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.Update(msg)

	updatedModel := newModel.(Model)
	// When in detail view, escape should return to list
	assert.Equal(t, ViewList, updatedModel.State())
}

func TestUpdateTabKey(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Send Tab key - should be handled without error
	msg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := model.Update(msg)

	assert.NotNil(t, newModel)
}

func TestUpdateReturnsModel(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(msg)

	_, ok := newModel.(Model)
	assert.True(t, ok, "Update should return a Model")
}

func TestUpdatePreservesRepositoryOnStateChange(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	msg := tea.WindowSizeMsg{Width: 120, Height: 30}
	newModel, _ := model.Update(msg)

	updatedModel := newModel.(Model)
	assert.Equal(t, repo, updatedModel.Repository())
}
