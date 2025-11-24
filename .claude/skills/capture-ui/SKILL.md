---
name: capture-ui
description: Use VHS to capture automated screenshots and terminal recordings of Bubble Tea applications. Activates when creating documentation, demonstrating features, generating GIFs for README files, capturing bug reproductions, or recording UI interactions.
---

# Capture UI - Automated Terminal Recordings with VHS

## Purpose
Create automated, reproducible terminal recordings and screenshots using VHS tape files for documentation and demonstration.

## When to Use
- Creating documentation screenshots for README or docs
- Recording feature demonstrations
- Generating animated GIFs for pull requests
- Capturing bug reproductions
- Building tutorial materials
- Recording UI interactions for testing/validation

## Instructions

### VHS Basics
VHS uses "tape" files (`.tape`) that script terminal interactions:
1. Start a shell command
2. Send keystrokes
3. Add delays
4. Take screenshots
5. Generate GIF/MP4/WebM output

### Creating a Tape File
Create a file in `testdata/vhs/` with `.tape` extension:
```tape
Output demo.gif
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 800

# Start the application
Type "./bin/claudex"
Enter
Sleep 2s

# Interact with UI
Type "search"
Sleep 500ms
Type "j"  # Navigate down
Sleep 500ms

# Screenshot at key moment
Screenshot demo-search.png

# More interactions
Enter  # Select
Sleep 2s

# Exit
Type "q"
Sleep 500ms
```

### Running VHS
```bash
# Run a single tape
vhs testdata/vhs/demo.tape

# Run all tapes
make vhs-test

# Run and record for docs
make vhs-record
```

### Tape File Commands
- **Output**: Set output filename (`.gif`, `.mp4`, `.webm`, `.png`)
- **Set**: Configure dimensions, fonts, theme, shell
- **Type**: Send keystrokes (`Enter`, `Tab`, arrow keys, etc.)
- **Sleep**: Pause execution (supports `ms`, `s` units)
- **Screenshot**: Capture still image at current frame
- **Enter**: Send Enter key (shorthand for Type Enter)
- **Ctrl+C**: Send interrupt signal

### Project Tape Files
Existing tapes in `testdata/vhs/`:
- **startup.tape**: Initial load and basic navigation
- **search.tape**: Search functionality demo

Create new tapes for:
- Feature demonstrations
- Bug reproductions
- Tutorial walkthroughs
- Documentation screenshots

### Best Practices
1. **Start with clean state**: Fresh terminal session
2. **Use consistent timing**: Human-like delays (300-500ms between actions)
3. **Add screenshots for docs**: Still images for key UI states
4. **Keep tapes focused**: One feature per tape file
5. **Test tapes regularly**: Run with `make vhs-test` before releases
6. **Version control tapes**: Commit `.tape` files, not generated assets (GIFs)

## Examples

**Example 1: Basic feature demo**
```tape
Output demo-search.gif
Set Shell bash
Set Width 1000
Set Height 600

Type "./bin/claudex"
Enter
Sleep 1.5s
Type "/"
Sleep 300ms
Type "model:sonnet"
Sleep 1s
Type "q"
```

**Example 2: Multi-step workflow**
```tape
Output workflow-demo.gif
Set Shell bash
Set Width 1200
Set Height 800

Type "./bin/claudex"
Enter
Sleep 2s

# Navigate list
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms

# Open conversation
Enter
Sleep 1.5s

# View metadata
Type "Tab"
Sleep 1s

# Exit
Type "q"
Sleep 500ms
```

**Example 3: Documentation screenshots (no animation)**
```tape
Output ignore.gif
Set Shell bash
Set Width 1400
Set Height 900

Type "./bin/claudex"
Enter
Sleep 2s

Screenshot docs/main-view.png
Sleep 1s

Type "/"
Sleep 500ms
Type "test"
Sleep 1s
Screenshot docs/search-view.png

Type "q"
```

**Example 4: Bug reproduction**
```tape
Output bug-reproduction.gif
Set Shell bash
Set Width 800
Set Height 600

Type "./bin/claudex"
Enter
Sleep 2s

# Steps to reproduce
Type "Tab"
Sleep 300ms
Type "Tab"
Sleep 300ms

Screenshot bug-state.png

Type "q"
```

## Advanced Features

### Variable Typing Speed
```tape
Set TypingSpeed 100ms  # Slower, more natural typing
Type "This types slowly"
Set TypingSpeed 10ms   # Fast typing
Type "This types quickly"
```

### Custom Themes
```tape
Set Theme "Catppuccin Mocha"
```

## Makefile Integration
```bash
make vhs-test    # Build and run all tapes, capture screenshots
make vhs-record  # Record GIFs for documentation
```

## Troubleshooting
- **VHS not found**: Install with `brew install vhs`
- **Output looks wrong**: Adjust Width/Height/FontSize in tape file
- **App doesn't start**: Ensure `./bin/claudex` binary exists (run `make build` first)
- **Timing issues**: Increase Sleep durations for slower systems
- **GIF too large**: Reduce Width/Height or use MP4 format instead
