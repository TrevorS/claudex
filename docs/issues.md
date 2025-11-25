# Claudex TUI Comprehensive Issues Audit

**Generated:** 2024-11-24
**Total Issues Found:** 44
**Status:** In progress - 9 fixed, 2 deferred (not bugs), 33 remaining

---

## Executive Summary

Three parallel explore agents audited the entire UI layer:
1. **List View & Sidebar** - 15 issues found
2. **Detail View & Metadata** - 14 issues found
3. **State/Keys/Layout** - 15 issues found

### Priority Distribution
| Severity | Count | Description |
|----------|-------|-------------|
| Critical | 5 | Crashes, data corruption, broken core functionality |
| High | 8 | Major UX problems, incorrect data display |
| Medium | 18 | Minor UX issues, inconsistencies, fragile code |
| Low | 13 | Code quality, edge cases, future-proofing |

---

## CRITICAL Issues (Fix Immediately)

### C1. Empty List KeyEnd Crash
- **File:** `internal/ui/views/list.go:120`
- **Issue:** If conversation list is empty, pressing End key calls `SetSelectedIndex(-1)`, which clamps to 0 but leaves invalid state
- **Impact:** Potential crash or undefined behavior with empty conversation list
- **Fix:** Add guard: `if len(v.conversations) == 0 { return v, nil }`
- **Status:** [x] Fixed - Added empty list guard for KeyHome and KeyEnd

### C2. UTF-8 Truncation Breaks Multi-byte Characters
- **Files:** `list.go:250-251`, `metadata.go:155-159`
- **Issue:** Title truncation uses byte length `len(s)` instead of rune count. Truncating mid-character produces invalid UTF-8
- **Impact:** Corrupted display, terminal rendering glitches with unicode titles
- **Fix:** Use `[]rune` conversion: `string([]rune(title)[:maxLen-3]) + "..."`
- **Status:** [x] Fixed - Added truncateString() helper using runes

### C3. Markdown Renderer Width Mismatch
- **File:** `internal/ui/views/detail.go:96, 421`
- **Issue:** Renderer initialized with hardcoded width 80, then recreated with different width on resize. Messages pre-rendered before resize have wrong wrapping
- **Impact:** Code blocks and text misaligned after window resize
- **Fix:** Defer markdown rendering until after SetSize, or always use calculated width
- **Status:** [x] Fixed - Now uses consistent width calculation (defaultWidth - 4)

### C4. Selection Highlighting Wrong Item After Sort
- **File:** `internal/ui/views/list.go:244`
- **Issue:** `renderTable()` compares `i == v.selectedIndex` against sorted array, but selectedIndex was set before sorting. Different items get highlighted
- **Impact:** User sees wrong conversation highlighted, confusing navigation
- **Fix:** Track selected conversation by ID, not index
- **Status:** [~] Deferred - Not currently triggerable (sort order is fixed), would need ID tracking for future sort cycling

### C5. ListView Not Reset on Filter Change
- **File:** `internal/ui/update.go:100-131`
- **Issue:** When project/filter selection changes, ListView.SetConversations() doesn't reset selectedIndex. Index 5 stays selected even if filtered list only has 3 items
- **Impact:** Out-of-bounds selection, wrong item highlighted or crash
- **Fix:** Call `m.listView.SetSelectedIndex(0)` after SetConversations()
- **Status:** [x] Fixed - Added SetSelectedIndex(0) after filter changes

---

## HIGH Issues (Fix Soon)

### H1. Repeated O(n log n) Sorting on Every Cursor Move
- **File:** `internal/ui/views/list.go:454`
- **Issue:** `SelectedConversation()` calls `getSortedConversations()` which re-sorts entire list. Called on every cursor move
- **Impact:** Performance degradation with large conversation lists (1000+)
- **Fix:** Cache sorted list, invalidate only when sort params change
- **Status:** [ ] Not started

### H2. Sidebar Navigation Wrapping Unexpectedly
- **File:** `internal/ui/components/sidebar.go:167-171`
- **Issue:** Pressing Up at "All" (-1) wraps to last project instead of staying at "All"
- **Impact:** Confusing navigation, users overshoot
- **Fix:** Clamp to -1 minimum, don't wrap
- **Status:** [~] Not a bug - Code already clamps correctly (if newIndex < -1 { newIndex = -1 })

### H3. Empty Message Content Renders Blank Lines
- **File:** `internal/ui/views/detail.go:284-285`
- **Issue:** Messages with empty content still render wrapped in messageContentStyle, creating visual gaps
- **Impact:** Wasted viewport space, confusing layout
- **Fix:** Skip rendering or show "(empty)" placeholder
- **Status:** [x] Fixed - Skip rendering empty message content

### H4. Metadata Shows Wrong Stats for Lazy-Loaded Conversations
- **File:** `internal/ui/components/metadata.go:72, 79-80, 142-151`
- **Issue:** Token breakdown and user message count only work with fully-loaded conversations. In list view, shows 0/0
- **Impact:** Misleading statistics in metadata panel
- **Fix:** Use cached token counts, or show "—" until loaded
- **Status:** [ ] Not started

### H5. Viewport Scroll Position Resets on Window Resize
- **File:** `internal/ui/views/detail.go:199`
- **Issue:** `initViewport()` always sets YPosition=0. Resize triggers re-init, jumping to top
- **Impact:** Users lose their place in long conversations
- **Fix:** Preserve Y position during resize
- **Status:** [x] Fixed - Preserve and restore YOffset during resize

### H6. SearchBar Focus Not Released on Exit
- **File:** `internal/ui/update.go:273-288`
- **Issue:** Exiting search mode doesn't call `m.searchBar.Blur()`. TextInput may continue consuming keystrokes
- **Impact:** Keyboard input trapped, navigation broken
- **Fix:** Call `m.searchBar.Blur()` before state transition
- **Status:** [x] Fixed - Call Blur() before exiting search mode

### H7. DetailView Cleanup Incomplete
- **File:** `internal/ui/update.go:67-72`
- **Issue:** BackToListMsg only sets `m.detailView = nil`. Pending async commands from old DetailView may still arrive
- **Impact:** Stale messages cause state corruption on rapid navigation
- **Fix:** Use generation counter to invalidate old messages
- **Status:** [ ] Not started

### H8. Selected Index Not Validated After Sort
- **File:** `internal/ui/views/list.go:413-423`
- **Issue:** SetConversations() clamps index but doesn't account for sort order change. Index now points to different conversation
- **Impact:** Metadata shows wrong conversation info
- **Fix:** Send CursorMovedMsg after SetConversations if selection changed
- **Status:** [ ] Not started

---

## MEDIUM Issues (Fix When Convenient)

### M1. FocusRight Never Receives Focus
- **File:** `internal/ui/update.go:313-322`
- **Issue:** Tab cycling only toggles Sidebar↔Center. FocusRight is defined but unreachable
- **Impact:** Metadata panel is read-only (may be intentional)
- **Fix:** Either remove FocusRight or add to Tab cycle
- **Status:** [ ] Not started

### M2. Global 'f' Filter Hotkey Conflicts with Sidebar
- **File:** `internal/ui/update.go:290-303`, `sidebar.go:144`
- **Issue:** Both global handler and sidebar handle 'f' key. Unclear which takes priority
- **Impact:** Inconsistent behavior depending on focus state
- **Fix:** Decide ownership: global-only or sidebar-only
- **Status:** [ ] Not started

### M3. Border Width Calculation Inconsistent
- **Files:** `view.go:214-223`, `update.go:214-223`
- **Issue:** Magic number 4 (border=2 + padding=2) hardcoded in multiple places
- **Impact:** Changing border style requires hunting through code
- **Fix:** Define constant `PANE_CHROME = 4` and use consistently
- **Status:** [ ] Not started

### M4. Separator Overflow in Narrow Terminals
- **File:** `internal/ui/views/detail.go:267-268`
- **Issue:** `strings.Repeat("─", v.viewWidth-4)` produces negative length if width < 4
- **Impact:** Visual glitch in very narrow terminals
- **Fix:** Guard: `max(1, v.viewWidth-4)`
- **Status:** [x] Fixed - Added width guard for separator

### M5. Model Name Truncation Not Handled
- **Files:** `list.go:369-372`, `detail.go:322-329`
- **Issue:** `formatModel()` splits on "-" but doesn't truncate result. Long model names overflow
- **Impact:** Metadata lines break with future long model names
- **Fix:** Truncate to max length after formatting
- **Status:** [ ] Not started

### M6. Token Count Precision Edge Cases
- **File:** `internal/ui/components/metadata.go:121-128`
- **Issue:** Float division produces awkward values: 999999 → "1000.0K"
- **Impact:** Ugly display at token count boundaries
- **Fix:** Round to sensible precision
- **Status:** [ ] Not started

### M7. Metadata Unreadable on Narrow Panels
- **File:** `internal/ui/components/metadata.go:95-109`
- **Issue:** Fixed 12-char key width leaves only 1 char for values at minimum panel width
- **Impact:** Metadata truncated to "..." in narrow terminals
- **Fix:** Make key width responsive
- **Status:** [ ] Not started

### M8. Duplicated Title Truncation Logic
- **File:** `internal/ui/views/list.go:306-308, 250-251`
- **Issue:** Same truncation in `buildRow()` and `renderTable()` with different max lengths
- **Impact:** Maintenance burden, inconsistent behavior
- **Fix:** Extract to `truncateTitle(title, maxLen)` helper
- **Status:** [ ] Not started

### M9. Pagination Bounds Check Incomplete
- **File:** `internal/ui/views/list.go:226-230`
- **Issue:** Edge case where startIdx > 0 but sorted is empty not handled
- **Impact:** Potential out-of-bounds access
- **Fix:** Add `if len(sorted) == 0 { return "No conversations" }`
- **Status:** [ ] Not started

### M10. SearchBar State Not Cleared on Exit
- **File:** `internal/ui/update.go:277-280`
- **Issue:** Error/result flags (hasError, errorMsg) not cleared when exiting search
- **Impact:** Stale error messages on re-entry
- **Fix:** Add `SearchBar.Clear()` method
- **Status:** [ ] Not started

### M11. Metadata Inconsistency During Detail Load
- **File:** `internal/ui/update.go:56-65`
- **Issue:** Metadata may show list cursor's conversation while DetailView loads different one
- **Impact:** Brief confusing display
- **Fix:** Show "Loading..." in metadata until DetailView loaded
- **Status:** [ ] Not started

### M12. SetConversations Index Clamping Issue
- **File:** `internal/ui/views/list.go:417-422`
- **Issue:** Clamps to len-1 but doesn't handle case where list shrinks
- **Impact:** Index may reference wrong conversation
- **Fix:** Always clamp and validate
- **Status:** [ ] Not started

### M13. Multi-line Project Styling Issue
- **File:** `internal/ui/components/sidebar.go:331`
- **Issue:** Single Lipgloss style applied to multi-line project entry. Selection background covers metadata line too
- **Impact:** Visual inconsistency
- **Fix:** Style each line separately
- **Status:** [ ] Not started

### M14. Fragile Model Name Parsing
- **File:** `internal/ui/views/list.go:369-372`
- **Issue:** Assumes "claude-" prefix. "gpt-4o" becomes "4o"
- **Impact:** Non-Claude models display incorrectly
- **Fix:** Hardcode known models or use regex
- **Status:** [ ] Not started

### M15. Date Filter Logic Flaw
- **File:** `internal/ui/components/sidebar.go:542`
- **Issue:** Doesn't handle negative durations (future-dated conversations)
- **Impact:** Future timestamps excluded from filters
- **Fix:** Include if diff < 0
- **Status:** [ ] Not started

### M16. Window Resize Doesn't Update Sidebar
- **File:** `internal/ui/update.go:140-150`
- **Issue:** WindowSizeMsg forwarded to ListView/DetailView but not Sidebar
- **Impact:** Sidebar dimensions stale after resize
- **Fix:** Add `m.sidebar = m.sidebar.SetSize(...)` in handleWindowSize
- **Status:** [ ] Not started

### M17. Dead Code: renderMinimalView == renderNarrowView
- **File:** `internal/ui/view.go:38-54`
- **Issue:** Both functions have identical implementation
- **Impact:** Code confusion
- **Fix:** Merge or implement distinct layouts
- **Status:** [ ] Not started

### M18. isSearchActive Semantic Confusion
- **File:** `internal/ui/update.go:104, 116`
- **Issue:** `SetSearchActive(true)` called for filter results, not actual searches
- **Impact:** Future features expecting search semantics may break
- **Fix:** Use different flag for filter state
- **Status:** [ ] Not started

---

## LOW Issues (Fix Eventually)

### L1. Hardcoded Page Size Fallback
- **File:** `internal/ui/views/list.go:192-195`
- **Issue:** Returns 10 if height < 5, unexpectedly large
- **Fix:** Use `max(3, v.height/2)`
- **Status:** [x] Fixed - Now uses max(3, height/2) fallback

### L2. Token Formatting Loss of Precision
- **File:** `internal/ui/views/list.go:356-363`
- **Issue:** 999 tokens shows as "1.0K"
- **Fix:** Add tooltip for future
- **Status:** [ ] Not started

### L3. extractUniqueProjects Nil Handling
- **File:** `internal/ui/components/sidebar.go:465`
- **Issue:** Empty string "" handled but no documentation
- **Fix:** Add comment explaining behavior
- **Status:** [ ] Not started

### L4. Unused Text Wrapping Code
- **File:** `internal/ui/views/detail.go:353-412`
- **Issue:** `wrapText()` and `wrapLine()` never called
- **Fix:** Remove dead code
- **Status:** [ ] Not started

### L5. No Control Character Sanitization
- **File:** `internal/ui/views/detail.go:284`
- **Issue:** Message content rendered without sanitizing control chars
- **Fix:** Add sanitization before render
- **Status:** [ ] Not started

### L6. Tab Characters in Code Blocks
- **File:** `internal/render/markdown.go:256-263`
- **Issue:** Tabs not expanded to spaces, inconsistent across terminals
- **Fix:** `strings.ReplaceAll("\t", "    ")`
- **Status:** [ ] Not started

### L7. Can't Distinguish 0 Tokens from Not Loaded
- **File:** `internal/ui/components/metadata.go:83-89`
- **Issue:** Both cases return 0, "Token Breakdown" section skipped
- **Fix:** Check if conversation fully loaded
- **Status:** [ ] Not started

### L8. Tab Key Not Forwarded in Detail View
- **File:** `internal/ui/update.go:325-328`
- **Issue:** Tab consumed but not forwarded to DetailView
- **Fix:** Forward to DetailView.Update()
- **Status:** [ ] Not started

### L9. Three-Pane Calculates Sidebar Width When Hidden
- **File:** `internal/ui/view.go:87-122`
- **Issue:** Dead calculations when sidebarVisible=false
- **Fix:** Early return in renderThreePane
- **Status:** [ ] Not started

### L10-L13. Minor Code Quality Issues
- Various files: Dead code, inconsistent naming, missing comments
- **Status:** [ ] Not started

---

## Files to Modify (By Priority)

### Critical/High Priority
| File | Issues |
|------|--------|
| `internal/ui/views/list.go` | C1, C2, C4, H1, H8, M8, M9, M12, M14 |
| `internal/ui/update.go` | C5, H6, H7, M1, M2, M10, M11, M16, M18 |
| `internal/ui/views/detail.go` | C3, H3, H5, M4, L4, L5 |
| `internal/ui/components/metadata.go` | C2, H4, M6, M7, L7 |
| `internal/ui/components/sidebar.go` | H2, M13, M15, L3 |

### Medium/Low Priority
| File | Issues |
|------|--------|
| `internal/ui/view.go` | M3, M17, L9 |
| `internal/render/markdown.go` | L6 |

---

## Recommended Fix Order

**Phase 1: Critical Stability (5 issues)**
1. C1 - Empty list crash guard
2. C5 - Filter selection reset
3. C2 - UTF-8 safe truncation
4. C4 - Track selection by ID
5. C3 - Markdown width fix

**Phase 2: High UX Issues (8 issues)**
1. H1 - Cache sorted conversations
2. H5 - Preserve scroll position
3. H6 - SearchBar blur
4. H2 - Sidebar navigation
5. H3 - Empty message handling
6. H4 - Lazy-load metadata
7. H7 - DetailView cleanup
8. H8 - Validate selection after sort

**Phase 3: Medium Polish (18 issues)**
- Group by file and fix together

**Phase 4: Low/Code Quality (13 issues)**
- Clean up as time permits

---

## Success Criteria

- [ ] All Critical issues resolved
- [ ] All High issues resolved
- [ ] No crashes on empty lists, filters, or rapid navigation
- [ ] UTF-8 titles display correctly
- [ ] Scroll position preserved on resize
- [ ] Performance acceptable with 1000+ conversations
- [ ] Keyboard navigation consistent across all views
