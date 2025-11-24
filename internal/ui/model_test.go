package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TrevorS/claudex/internal/repository"
)

func TestNewModel(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	assert.NotNil(t, model)
	assert.Equal(t, ViewList, model.state)
	assert.Equal(t, 0, model.width)
	assert.Equal(t, 0, model.height)
	assert.Nil(t, model.err)
}

func TestModelInit(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Init can return nil (no initial command needed)
	_ = model.Init()
}

func TestViewStateValues(t *testing.T) {
	// Verify ViewState constants exist and are unique
	states := []ViewState{ViewList, ViewDetail, ViewSearch, ViewStats, ViewPalette}
	seen := make(map[ViewState]bool)

	for _, state := range states {
		assert.False(t, seen[state], "duplicate ViewState value")
		seen[state] = true
	}
}

func TestModelRepository(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	assert.Equal(t, repo, model.repository)
}

func TestModelDimensions(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Initial dimensions should be zero
	assert.Equal(t, 0, model.width)
	assert.Equal(t, 0, model.height)
}

func TestModelConversationState(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Initial conversation state
	assert.Nil(t, model.conversations)
	assert.Equal(t, "", model.selected)
	assert.Equal(t, 0, model.scrollPos)
}
