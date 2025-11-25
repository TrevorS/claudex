# VHS Tape Examples

## Basic Feature Demo

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

## Multi-Step Workflow

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

## Documentation Screenshots

For static images without animation:

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

## Bug Reproduction

Capture steps to reproduce a bug:

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

## Advanced: Variable Typing Speed

```tape
Set TypingSpeed 100ms  # Slower, more natural
Type "This types slowly"

Set TypingSpeed 10ms   # Fast
Type "This types quickly"
```

## Advanced: Custom Themes

```tape
Set Theme "Catppuccin Mocha"
```
