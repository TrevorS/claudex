// ABOUTME: markdown.go implements a markdown-to-terminal renderer that converts
// markdown syntax to terminal-friendly styled text using Lipgloss. Supports text
// formatting, code blocks, lists, headings, block quotes, and links.
package render

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Renderer converts markdown to terminal-friendly text.
type Renderer struct {
	maxWidth int

	// Pre-compiled regexes for performance
	codeBlockRegex  *regexp.Regexp
	boldRegex       *regexp.Regexp
	italicRegex     *regexp.Regexp
	inlineCodeRegex *regexp.Regexp
	strikeRegex     *regexp.Regexp
	linkRegex       *regexp.Regexp
	imageRegex      *regexp.Regexp
	headingRegex    *regexp.Regexp
	hrRegex         *regexp.Regexp
	blockQuoteRegex *regexp.Regexp
	unorderedRegex  *regexp.Regexp
	orderedRegex    *regexp.Regexp

	// Styles
	codeBlockStyle  lipgloss.Style
	inlineCodeStyle lipgloss.Style
	boldStyle       lipgloss.Style
	italicStyle     lipgloss.Style
	strikeStyle     lipgloss.Style
	headingStyle    lipgloss.Style
	quoteStyle      lipgloss.Style
	hrStyle         lipgloss.Style
}

// CodeBlock represents an extracted code block with language and content.
type CodeBlock struct {
	Language string
	Content  string
	StartPos int
	EndPos   int
}

// NewRenderer creates a new markdown renderer with the specified max width.
func NewRenderer(maxWidth int) *Renderer {
	return &Renderer{
		maxWidth: maxWidth,

		// Compile regexes once
		codeBlockRegex:  regexp.MustCompile("```([a-zA-Z0-9]*)\n([\\s\\S]*?)```"),
		boldRegex:       regexp.MustCompile(`\*\*(.+?)\*\*|__(.+?)__`),
		italicRegex:     regexp.MustCompile(`\*(.+?)\*|_(.+?)_`),
		inlineCodeRegex: regexp.MustCompile("`([^`]+?)`"),
		strikeRegex:     regexp.MustCompile(`~~(.+?)~~`),
		linkRegex:       regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`),
		imageRegex:      regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`),
		headingRegex:    regexp.MustCompile(`^(#{1,6})\s+(.+)$`),
		hrRegex:         regexp.MustCompile(`^(---+|\*\*\*+|___+)\s*$`),
		blockQuoteRegex: regexp.MustCompile(`^>\s+(.+)$`),
		unorderedRegex:  regexp.MustCompile(`^(\s*)([-*+])\s+(.+)$`),
		orderedRegex:    regexp.MustCompile(`^(\s*)(\d+)\.\s+(.+)$`),

		// Initialize styles
		codeBlockStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1),
		inlineCodeStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("229")).
			Padding(0, 1),
		boldStyle: lipgloss.NewStyle().
			Bold(true),
		italicStyle: lipgloss.NewStyle().
			Italic(true),
		strikeStyle: lipgloss.NewStyle().
			Strikethrough(true),
		headingStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			MarginTop(1).
			MarginBottom(1),
		quoteStyle: lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("8")).
			PaddingLeft(2).
			Faint(true),
		hrStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1).
			MarginBottom(1),
	}
}

// Render converts markdown text to terminal-friendly styled output.
func (r *Renderer) Render(markdown string) string {
	if markdown == "" {
		return ""
	}

	// Step 1: Extract code blocks and replace with placeholders
	codeBlocks := r.ParseCodeBlocks(markdown)
	text := markdown
	placeholders := make(map[string]string)

	for i, block := range codeBlocks {
		placeholder := "\x00CODE_BLOCK_" + string(rune(i)) + "\x00"
		placeholders[placeholder] = r.RenderCodeBlock(block)
		// Replace the entire code block with placeholder
		text = strings.Replace(text, markdown[block.StartPos:block.EndPos], placeholder, 1)
	}

	// Step 2: Process line by line for block-level elements
	lines := strings.Split(text, "\n")
	var result []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check for code block placeholder
		if strings.Contains(line, "\x00CODE_BLOCK_") {
			for placeholder, rendered := range placeholders {
				if strings.Contains(line, placeholder) {
					result = append(result, rendered)
					break
				}
			}
			continue
		}

		// Horizontal rule
		if r.hrRegex.MatchString(line) {
			separator := strings.Repeat("─", r.maxWidth-4)
			result = append(result, r.hrStyle.Render(separator))
			continue
		}

		// Heading
		if matches := r.headingRegex.FindStringSubmatch(line); matches != nil {
			level := len(matches[1])
			text := matches[2]
			text = r.RenderInlineFormatting(text)

			var styled string
			if level == 1 {
				// H1: bold + larger
				styled = r.headingStyle.Render(text)
			} else {
				// H2-H6: just bold
				styled = r.boldStyle.Render(text)
			}
			result = append(result, styled)
			continue
		}

		// Block quote
		if matches := r.blockQuoteRegex.FindStringSubmatch(line); matches != nil {
			quoted := r.RenderInlineFormatting(matches[1])
			result = append(result, r.quoteStyle.Render(quoted))
			continue
		}

		// Unordered list
		if matches := r.unorderedRegex.FindStringSubmatch(line); matches != nil {
			indent := matches[1]
			content := r.RenderInlineFormatting(matches[3])
			bullet := indent + "  • " + content
			result = append(result, bullet)
			continue
		}

		// Ordered list
		if matches := r.orderedRegex.FindStringSubmatch(line); matches != nil {
			indent := matches[1]
			number := matches[2]
			content := r.RenderInlineFormatting(matches[3])
			numbered := indent + "  " + number + ". " + content
			result = append(result, numbered)
			continue
		}

		// Regular paragraph with inline formatting
		if line != "" {
			formatted := r.RenderInlineFormatting(line)
			result = append(result, formatted)
		} else {
			// Preserve empty lines for paragraph breaks
			result = append(result, "")
		}
	}

	return strings.Join(result, "\n")
}

// ParseCodeBlocks extracts all fenced code blocks from the text.
func (r *Renderer) ParseCodeBlocks(text string) []CodeBlock {
	matches := r.codeBlockRegex.FindAllStringSubmatchIndex(text, -1)
	var blocks []CodeBlock

	for _, match := range matches {
		// match[0], match[1] = full match start/end
		// match[2], match[3] = language group start/end
		// match[4], match[5] = content group start/end

		var language string
		if match[2] >= 0 && match[3] >= 0 {
			language = text[match[2]:match[3]]
		}

		var content string
		if match[4] >= 0 && match[5] >= 0 {
			content = text[match[4]:match[5]]
		}

		blocks = append(blocks, CodeBlock{
			Language: language,
			Content:  strings.TrimSpace(content),
			StartPos: match[0],
			EndPos:   match[1],
		})
	}

	return blocks
}

// RenderCodeBlock renders a code block with styling and optional syntax highlighting.
func (r *Renderer) RenderCodeBlock(block CodeBlock) string {
	return r.RenderCodeBlockWithHighlight(block, true)
}

// RenderCodeBlockWithHighlight renders a code block with optional syntax highlighting.
func (r *Renderer) RenderCodeBlockWithHighlight(block CodeBlock, highlight bool) string {
	var parts []string

	// Add language label if present
	if block.Language != "" {
		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Faint(true).
			Render(block.Language)
		parts = append(parts, label)
	}

	// Prepare content
	content := block.Content

	// Apply syntax highlighting if requested and language is known
	if highlight && block.Language != "" {
		highlighted := HighlightQuick(content, block.Language)
		if highlighted != "" {
			content = highlighted
		}
		// If highlighting fails or returns empty, use original content
	}

	// Render code content with styling
	// Note: Lipgloss preserves ANSI codes from Chroma, so highlighted text works correctly
	styled := r.codeBlockStyle.Render(content)
	parts = append(parts, styled)

	return strings.Join(parts, "\n")
}

// RenderInlineFormatting applies inline formatting (bold, italic, code, links).
func (r *Renderer) RenderInlineFormatting(text string) string {
	// Process in order: code (to protect it), then links/images, then text formatting

	// Step 1: Inline code (protect from other formatting)
	codeMatches := r.inlineCodeRegex.FindAllStringSubmatch(text, -1)
	codeReplacements := make([]string, 0)
	for i, match := range codeMatches {
		placeholder := "\x00INLINECODE" + string(rune('A'+i)) + "\x00"
		rendered := r.inlineCodeStyle.Render(match[1])
		codeReplacements = append(codeReplacements, rendered)
		text = strings.Replace(text, match[0], placeholder, 1)
	}

	// Step 2: Links - extract text only
	text = r.linkRegex.ReplaceAllString(text, "$1")

	// Step 3: Images - extract alt text only
	text = r.imageRegex.ReplaceAllString(text, "$1")

	// Step 4: Strikethrough
	for {
		match := r.strikeRegex.FindStringSubmatch(text)
		if match == nil {
			break
		}
		styled := r.strikeStyle.Render(match[1])
		text = strings.Replace(text, match[0], styled, 1)
	}

	// Step 5: Bold (before italic to handle nested cases)
	for {
		match := r.boldRegex.FindStringSubmatch(text)
		if match == nil {
			break
		}
		// Get the captured content (either group 1 or group 2)
		content := match[1]
		if content == "" {
			content = match[2]
		}
		styled := r.boldStyle.Render(content)
		text = strings.Replace(text, match[0], styled, 1)
	}

	// Step 6: Italic
	// Be more conservative - only match when surrounded by whitespace or punctuation
	// to avoid matching underscores in URLs/filenames
	conservativeItalicRegex := regexp.MustCompile(`(^|\s|\()\*([^*]+?)\*($|\s|\)|,|\.)|(^|\s|\()_([^_]+?)_($|\s|\)|,|\.)`)
	for {
		match := conservativeItalicRegex.FindStringSubmatch(text)
		if match == nil {
			break
		}
		// Get the captured content (either group 2 or group 5)
		content := match[2]
		if content == "" {
			content = match[5]
		}
		before := match[1]
		if before == "" {
			before = match[4]
		}
		after := match[3]
		if after == "" {
			after = match[6]
		}
		styled := before + r.italicStyle.Render(content) + after
		text = strings.Replace(text, match[0], styled, 1)
	}

	// Step 7: Restore inline code
	for i, rendered := range codeReplacements {
		placeholder := "\x00INLINECODE" + string(rune('A'+i)) + "\x00"
		text = strings.Replace(text, placeholder, rendered, 1)
	}

	return text
}
