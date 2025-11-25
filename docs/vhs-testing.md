# VHS Testing Guide

This document describes the VHS-based UI testing approach for claudex.

## Overview

[VHS](https://github.com/charmbracelet/vhs) is a tool for creating terminal recordings and screenshots. We use it for:

1. **UI Regression Testing** - Capture screenshots at key states to detect visual changes
2. **Documentation** - Generate GIFs showing how to use the application
3. **CI Integration** - Automated visual testing in GitHub Actions

## Directory Structure

```
testdata/vhs/
├── startup.tape      # Application startup tests
├── navigation.tape   # List navigation tests
├── sidebar.tape      # Sidebar and focus management tests
├── detail.tape       # Detail view tests
├── search.tape       # Search functionality tests
├── output/           # Generated screenshots and GIFs
└── expected/         # Expected screenshots for comparison
```

## Running Tests

### Prerequisites

1. Install VHS: https://github.com/charmbracelet/vhs#installation
2. Build the application: `make build`

### Commands

```bash
# Run all VHS tests
make vhs-test

# Run with verbose output
make vhs-test-verbose

# Record all tapes (generate GIFs + screenshots)
make vhs-record

# Update expected screenshots after intentional UI changes
make vhs-update
```

### Running Individual Tests

Use the test script directly:

```bash
# Run a specific tape
./scripts/vhs-test.sh -t startup

# Run with verbose output
./scripts/vhs-test.sh -v -t navigation
```

## Writing VHS Tapes

### Basic Structure

```tape
Output testdata/vhs/example.gif

# Test description comment
# Tests: feature1, feature2, feature3

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

# Launch the application
Type "./bin/claudex"
Enter
Sleep 2s

# Test action
Type "j"
Sleep 300ms
Screenshot testdata/vhs/example-state.png

# Clean exit
Ctrl+C
```

### Key Commands

| Command | Description |
|---------|-------------|
| `Type "text"` | Type text character by character |
| `Type@100ms "text"` | Type with custom delay per character |
| `Enter` | Press Enter key |
| `Escape` | Press Escape key |
| `Tab` | Press Tab key |
| `Ctrl+C` | Press Ctrl+C |
| `Ctrl+F` | Press Ctrl+F |
| `Up`, `Down`, `Left`, `Right` | Arrow keys |
| `Home`, `End` | Jump to start/end |
| `PageUp`, `PageDown` | Page navigation |
| `Space` | Press spacebar |
| `Sleep 500ms` | Wait for specified duration |
| `Screenshot path.png` | Capture screenshot |

### Best Practices

1. **Use descriptive comments** - Document what each section tests
2. **Add Sleep after actions** - Give the UI time to update (300-500ms)
3. **Take screenshots at key states** - Capture before/after for comparisons
4. **Use consistent naming** - `tapename-statename.png`
5. **Test one feature per tape** - Keep tests focused and maintainable

## Test Scenarios

### startup.tape
- Application launch
- Three-pane layout rendering
- Initial focus state (sidebar)
- Header/footer display

### navigation.tape
- j/k vim-style navigation
- Arrow key navigation
- Home/End jumping
- PageUp/PageDown scrolling
- Focus switching with Tab

### sidebar.tape
- Project list navigation
- Project selection with Enter
- Date filter cycling (f key)
- Focus management
- Clear selection with Escape

### detail.tape
- Open conversation with Enter
- Message scrolling (Space, PageUp/Down)
- Viewport navigation (Home/End)
- Return to list (Escape, b key)

### search.tape
- Enter search mode (Ctrl+F)
- Type search query
- Execute search (Enter)
- Clear/exit search (Escape)

## Adding New Tests

1. Create a new `.tape` file in `testdata/vhs/`
2. Follow the structure above
3. Run `make vhs-record` to generate initial screenshots
4. Run `make vhs-update` to save as expected
5. Add to CI workflow if needed

## Troubleshooting

### VHS Not Installed
```
Error: VHS is not installed.
Install VHS from: https://github.com/charmbracelet/vhs
```

Install using one of:
- `brew install vhs` (macOS)
- `go install github.com/charmbracelet/vhs@latest`
- See VHS documentation for other options

### Screenshots Don't Match
1. Check if the UI changed intentionally
2. If intentional, run `make vhs-update`
3. If not intentional, investigate the regression

### Flaky Tests
- Increase `Sleep` durations after actions
- Check for race conditions in async operations
- Ensure deterministic test data

## CI Integration

Add to `.github/workflows/test.yml`:

```yaml
- name: Install VHS
  run: |
    go install github.com/charmbracelet/vhs@latest

- name: Run VHS tests
  run: make vhs-test
```

## Future Improvements

- [ ] Add image comparison using ImageMagick
- [ ] Generate test reports with visual diffs
- [ ] Add more edge case tests (empty states, errors)
- [ ] Test responsive layouts at different widths
- [ ] Add accessibility testing scenarios
