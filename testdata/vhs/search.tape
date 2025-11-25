Output testdata/vhs/search.gif

Set FontSize 14
Set Width 1200
Set Height 800
Set Shell "bash"

Type "./bin/claudex"
Enter
Sleep 2s

# Open search
Type "/"
Sleep 500ms
Type "test"
Sleep 1s
Screenshot testdata/vhs/search-results.png

Ctrl+C
