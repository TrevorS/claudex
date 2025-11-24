---
name: capture-tui
description: Capture terminal UI recordings and screenshots with VHS. Use when: creating documentation GIFs, recording feature demos, capturing bug reproductions, generating README screenshots, or debugging visual state transitions in Bubble Tea apps.
---

# Capture TUI - Terminal Recordings with VHS

Create automated, reproducible terminal recordings and screenshots using VHS tape files.

## When to Use

- Documentation screenshots for README/docs
- Feature demonstration GIFs
- Bug reproduction recordings
- Visual regression testing

## Quick Start

Create `testdata/vhs/demo.tape`:
```tape
Output demo.gif
Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 800

Type "./bin/claudex"
Enter
Sleep 2s

Screenshot demo.png

Type "q"
```

Run it:
```bash
vhs testdata/vhs/demo.tape
```

## Tape Commands

| Command | Description |
|---------|-------------|
| `Output file.gif` | Set output file (.gif, .mp4, .webm, .png) |
| `Set Width/Height N` | Terminal dimensions |
| `Set FontSize N` | Font size in pixels |
| `Set Shell bash` | Shell to use |
| `Type "text"` | Send keystrokes |
| `Enter` | Send Enter key |
| `Sleep 500ms` | Pause execution |
| `Screenshot file.png` | Capture still image |
| `Ctrl+C` | Send interrupt |

## Running VHS

```bash
# Single tape
vhs testdata/vhs/demo.tape

# All tapes (build + run)
make vhs-test

# Record for documentation
make vhs-record
```

## Project Tapes

Existing tapes in `testdata/vhs/`:
- `startup.tape` - Initial load and navigation
- `search.tape` - Search functionality

## Best Practices

1. Use human-like delays (300-500ms between actions)
2. Keep tapes focused - one feature per file
3. Commit `.tape` files, not generated GIFs
4. Run `make vhs-test` before releases

## Troubleshooting

- **VHS not found**: `brew install vhs`
- **Output looks wrong**: Adjust Width/Height/FontSize
- **App doesn't start**: Run `make build` first
- **Timing issues**: Increase Sleep durations
- **GIF too large**: Reduce dimensions or use MP4

## More Help

For detailed examples and advanced usage:
- **Examples**: See `examples.md` for tape file templates (feature demos, screenshots, bug reproduction)
- **Debugging**: See `debugging.md` for ffmpeg frame extraction and state transition debugging
