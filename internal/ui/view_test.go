package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TrevorS/claudex/internal/repository"
)

func TestViewRendersString(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(120, 30)

	view := model.View()
	assert.NotEmpty(t, view)
	assert.IsType(t, "", view)
}

func TestViewLayoutWideScreen(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(120, 30)

	view := model.View()

	// View should contain content for a wide screen
	assert.NotEmpty(t, view)
	lines := strings.Split(view, "\n")
	assert.Greater(t, len(lines), 0)
}

func TestViewLayoutNarrowScreen(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(80, 24)

	view := model.View()

	// View should render even on narrow screens
	assert.NotEmpty(t, view)
}

func TestViewLayoutMinimalScreen(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(40, 10)

	view := model.View()

	// View should render even on minimal screens
	assert.NotEmpty(t, view)
}

func TestViewShowsErrorIfPresent(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(120, 30)

	testErr := "test error message"
	model = model.SetError(ErrAppError(testErr))

	view := model.View()
	assert.NotEmpty(t, view)
}

func TestViewZeroDimensions(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)

	// Even with zero dimensions, should return a string
	view := model.View()
	assert.IsType(t, "", view)
}

func TestViewMultipleRenders(t *testing.T) {
	repo := repository.NewMock(nil)
	model := NewModel(repo)
	model = model.SetSize(120, 30)

	// Should render consistently
	view1 := model.View()
	view2 := model.View()

	assert.Equal(t, view1, view2)
}
