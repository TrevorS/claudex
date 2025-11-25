Output testdata/vhs/navigation.gif

# VHS Test: List Navigation
# Tests: j/k navigation, arrow keys, Home/End, PageUp/PageDown, Enter to detail

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

# Launch the application
Type "./bin/claudex"
Enter
Sleep 2s

# First, switch focus to center pane for navigation testing
Tab
Sleep 300ms
Screenshot testdata/vhs/nav-center-focus.png

# Test: Navigate down with j key (vim-style)
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 300ms
Screenshot testdata/vhs/nav-j-down.png

# Test: Navigate up with k key (vim-style)
Type "k"
Sleep 200ms
Type "k"
Sleep 300ms
Screenshot testdata/vhs/nav-k-up.png

# Test: Navigate with arrow keys
Down
Sleep 200ms
Down
Sleep 200ms
Screenshot testdata/vhs/nav-arrow-down.png

Up
Sleep 200ms
Screenshot testdata/vhs/nav-arrow-up.png

# Test: Jump to end with End key
End
Sleep 300ms
Screenshot testdata/vhs/nav-end.png

# Test: Jump to start with Home key
Home
Sleep 300ms
Screenshot testdata/vhs/nav-home.png

# Test: Page navigation
PageDown
Sleep 300ms
Screenshot testdata/vhs/nav-pagedown.png

PageUp
Sleep 300ms
Screenshot testdata/vhs/nav-pageup.png

# Clean exit
Ctrl+C
