Output testdata/vhs/detail.gif

# VHS Test: Detail View
# Tests: Enter to open, message display, scrolling, Escape/b to return

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

# Launch the application
Type "./bin/claudex"
Enter
Sleep 2s

# Switch to center pane and select first conversation
Tab
Sleep 300ms

# Test: Open conversation detail with Enter
Enter
Sleep 1s
Screenshot testdata/vhs/detail-opened.png

# Test: Scroll down in detail view with Space
Space
Sleep 500ms
Screenshot testdata/vhs/detail-scroll-down.png

# Test: Scroll down with PageDown
PageDown
Sleep 500ms
Screenshot testdata/vhs/detail-pagedown.png

# Test: Scroll up with PageUp
PageUp
Sleep 500ms
Screenshot testdata/vhs/detail-pageup.png

# Test: Jump to end
End
Sleep 500ms
Screenshot testdata/vhs/detail-end.png

# Test: Jump to start
Home
Sleep 500ms
Screenshot testdata/vhs/detail-home.png

# Test: Return to list with Escape
Escape
Sleep 500ms
Screenshot testdata/vhs/detail-return-escape.png

# Test: Open again and return with 'b' key
Enter
Sleep 1s
Type "b"
Sleep 500ms
Screenshot testdata/vhs/detail-return-b.png

# Clean exit
Ctrl+C
