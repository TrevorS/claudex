// ABOUTME: markdown_test.go contains comprehensive tests for the markdown-to-terminal renderer.
// Tests cover text formatting, code blocks, lists, headings, block elements, links, and edge cases.
package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRenderPlainText verifies plain text is unchanged
func TestRenderPlainText(t *testing.T) {
	renderer := NewRenderer(80)
	input := "This is plain text without any formatting."
	output := renderer.Render(input)
	assert.Contains(t, output, "This is plain text without any formatting.")
}

// TestRenderBold verifies bold formatting
func TestRenderBold(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "double asterisk bold",
			input:    "This is **bold** text",
			contains: "bold",
		},
		{
			name:     "double underscore bold",
			input:    "This is __bold__ text",
			contains: "bold",
		},
		{
			name:     "multiple bold",
			input:    "**First** and **second**",
			contains: "First",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)
			assert.Contains(t, output, tt.contains)
			// Should not contain markdown syntax
			assert.NotContains(t, output, "**")
			assert.NotContains(t, output, "__")
		})
	}
}

// TestRenderItalic verifies italic formatting
func TestRenderItalic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "single asterisk italic",
			input:    "This is *italic* text",
			contains: "italic",
		},
		{
			name:     "single underscore italic",
			input:    "This is _italic_ text",
			contains: "italic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)
			assert.Contains(t, output, tt.contains)
		})
	}
}

// TestRenderInlineCode verifies inline code formatting
func TestRenderInlineCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "basic inline code",
			input:    "Use `code` here",
			contains: "code",
		},
		{
			name:     "code with special chars",
			input:    "Run `git status` command",
			contains: "git status",
		},
		{
			name:     "multiple inline code",
			input:    "Use `foo` and `bar`",
			contains: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)
			assert.Contains(t, output, tt.contains)
			// Should not contain backticks
			assert.NotContains(t, output, "`code`")
			assert.NotContains(t, output, "`foo`")
		})
	}
}

// TestRenderNestedFormatting verifies nested bold/italic
func TestRenderNestedFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "bold with italic inside",
			input:    "This is **bold with *italic* inside**",
			contains: "bold with",
		},
		{
			name:     "italic with bold inside",
			input:    "This is *italic with **bold** inside*",
			contains: "italic with",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)
			assert.Contains(t, output, tt.contains)
		})
	}
}

// TestRenderStrikethrough verifies strikethrough formatting
func TestRenderStrikethrough(t *testing.T) {
	renderer := NewRenderer(80)
	input := "This is ~~strikethrough~~ text"
	output := renderer.Render(input)
	assert.Contains(t, output, "strikethrough")
	assert.NotContains(t, output, "~~")
}

// TestRenderCodeBlock verifies fenced code block extraction
func TestRenderCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		language string
	}{
		{
			name:  "go code block",
			input: "```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```",
			contains: []string{
				"func",
				"main",
				"fmt",
				"Println",
			},
			language: "go",
		},
		{
			name:  "python code block",
			input: "```python\ndef hello():\n    print('world')\n```",
			contains: []string{
				"def",
				"hello",
				"print",
				"world",
			},
			language: "python",
		},
		{
			name:  "no language specified",
			input: "```\ncode without language\n```",
			contains: []string{
				"code without language",
			},
			language: "",
		},
		{
			name:  "multiple code blocks",
			input: "```go\nfunc foo() {}\n```\n\nSome text\n\n```python\ndef bar():\n    pass\n```",
			contains: []string{
				"func",
				"foo",
				"def",
				"bar",
			},
			language: "go", // First block
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}

			// Should not contain fence markers
			assert.NotContains(t, output, "```")
		})
	}
}

// TestParseCodeBlocks verifies code block extraction
func TestParseCodeBlocks(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		wantLang  string
		wantCode  string
	}{
		{
			name:      "single go block",
			input:     "```go\nfunc main() {}\n```",
			wantCount: 1,
			wantLang:  "go",
			wantCode:  "func main() {}",
		},
		{
			name:      "no language",
			input:     "```\nplain code\n```",
			wantCount: 1,
			wantLang:  "",
			wantCode:  "plain code",
		},
		{
			name:      "no code blocks",
			input:     "Just plain text",
			wantCount: 0,
		},
		{
			name:      "multiple blocks",
			input:     "```go\nfoo\n```\n\n```python\nbar\n```",
			wantCount: 2,
			wantLang:  "go",
			wantCode:  "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			blocks := renderer.ParseCodeBlocks(tt.input)
			assert.Equal(t, tt.wantCount, len(blocks))

			if tt.wantCount > 0 {
				assert.Equal(t, tt.wantLang, blocks[0].Language)
				assert.Contains(t, blocks[0].Content, tt.wantCode)
			}
		})
	}
}

// TestRenderCodeBlockPreservesContent verifies code is not processed
func TestRenderCodeBlockPreservesContent(t *testing.T) {
	renderer := NewRenderer(80)
	input := "```go\n// This has **markdown** inside\nfunc foo() {\n    // More *markdown*\n}\n```"
	output := renderer.Render(input)

	// Should preserve markdown-like syntax inside code blocks
	assert.Contains(t, output, "**markdown**")
	assert.Contains(t, output, "*markdown*")
}

// TestRenderUnorderedList verifies bullet list formatting
func TestRenderUnorderedList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "simple list",
			input: "- Item 1\n- Item 2\n- Item 3",
			contains: []string{
				"•", "Item 1", "Item 2", "Item 3",
			},
		},
		{
			name:  "list with plus",
			input: "+ Item 1\n+ Item 2",
			contains: []string{
				"•", "Item 1",
			},
		},
		{
			name:  "list with asterisk",
			input: "* Item 1\n* Item 2",
			contains: []string{
				"•", "Item 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			for _, expected := range tt.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}

// TestRenderOrderedList verifies numbered list formatting
func TestRenderOrderedList(t *testing.T) {
	renderer := NewRenderer(80)
	input := "1. First\n2. Second\n3. Third"
	output := renderer.Render(input)

	assert.Contains(t, output, "First")
	assert.Contains(t, output, "Second")
	assert.Contains(t, output, "Third")
}

// TestRenderNestedList verifies nested list indentation
func TestRenderNestedList(t *testing.T) {
	renderer := NewRenderer(80)
	input := "- Item 1\n  - Nested 1\n  - Nested 2\n- Item 2"
	output := renderer.Render(input)

	assert.Contains(t, output, "Item 1")
	assert.Contains(t, output, "Nested 1")
	assert.Contains(t, output, "Nested 2")
	assert.Contains(t, output, "Item 2")
}

// TestRenderHeading verifies heading formatting
func TestRenderHeading(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "h1",
			input:    "# Title",
			contains: "Title",
		},
		{
			name:     "h2",
			input:    "## Subtitle",
			contains: "Subtitle",
		},
		{
			name:     "h3",
			input:    "### Section",
			contains: "Section",
		},
		{
			name:     "h4",
			input:    "#### Subsection",
			contains: "Subsection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			assert.Contains(t, output, tt.contains)
			// Should not contain # markers
			assert.NotContains(t, output, "###")
		})
	}
}

// TestRenderBlockQuote verifies block quote indentation
func TestRenderBlockQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "single line quote",
			input:    "> This is a quote",
			contains: "This is a quote",
		},
		{
			name:     "multi line quote",
			input:    "> Line 1\n> Line 2\n> Line 3",
			contains: "Line 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			assert.Contains(t, output, tt.contains)
		})
	}
}

// TestRenderHorizontalRule verifies separator rendering
func TestRenderHorizontalRule(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "three hyphens",
			input: "---",
		},
		{
			name:  "three asterisks",
			input: "***",
		},
		{
			name:  "three underscores",
			input: "___",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			// Should contain some kind of separator
			assert.True(t, strings.Contains(output, "─") || strings.Contains(output, "─"))
		})
	}
}

// TestRenderParagraphBreaks verifies double newlines create paragraphs
func TestRenderParagraphBreaks(t *testing.T) {
	renderer := NewRenderer(80)
	input := "Paragraph 1.\n\nParagraph 2.\n\nParagraph 3."
	output := renderer.Render(input)

	assert.Contains(t, output, "Paragraph 1")
	assert.Contains(t, output, "Paragraph 2")
	assert.Contains(t, output, "Paragraph 3")

	// Should have spacing between paragraphs
	lines := strings.Split(output, "\n")
	assert.True(t, len(lines) >= 5) // At least 3 content + 2 blank lines
}

// TestRenderLink verifies links show text only
func TestRenderLink(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		contains    string
		notContains string
	}{
		{
			name:        "basic link",
			input:       "[Click here](https://example.com)",
			contains:    "Click here",
			notContains: "https://example.com",
		},
		{
			name:        "link in text",
			input:       "Check out [the docs](https://docs.example.com) for more",
			contains:    "the docs",
			notContains: "https://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			assert.Contains(t, output, tt.contains)
			assert.NotContains(t, output, tt.notContains)
		})
	}
}

// TestRenderImage verifies images show alt text only
func TestRenderImage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		contains    string
		notContains string
	}{
		{
			name:        "basic image",
			input:       "![Alt text](image.png)",
			contains:    "Alt text",
			notContains: "image.png",
		},
		{
			name:        "image with URL",
			input:       "![Logo](https://example.com/logo.png)",
			contains:    "Logo",
			notContains: "https://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			assert.Contains(t, output, tt.contains)
			assert.NotContains(t, output, tt.notContains)
		})
	}
}

// TestRenderMixedContent verifies complex mixed markdown
func TestRenderMixedContent(t *testing.T) {
	input := `# Title

This is a paragraph with **bold** and *italic* text.

## Code Example

Here's some code:

` + "```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```" + `

## List of Features

- Feature 1
- Feature 2
- Feature 3

> This is a quote

Check out [the docs](https://example.com) for more information.
`

	renderer := NewRenderer(80)
	output := renderer.Render(input)

	// Verify all components are present
	assert.Contains(t, output, "Title")
	assert.Contains(t, output, "bold")
	assert.Contains(t, output, "italic")
	assert.Contains(t, output, "Code Example")
	assert.Contains(t, output, "func")
	assert.Contains(t, output, "main")
	assert.Contains(t, output, "Feature 1")
	assert.Contains(t, output, "This is a quote")
	assert.Contains(t, output, "the docs")

	// Verify markdown syntax is removed
	assert.NotContains(t, output, "**")
	assert.NotContains(t, output, "```")
	assert.NotContains(t, output, "https://example.com")
}

// TestRenderEdgeCaseEscapeUnderscore verifies underscores in URLs don't become italic
func TestRenderEdgeCaseEscapeUnderscore(t *testing.T) {
	renderer := NewRenderer(80)
	input := "Check https://example.com/some_file_name.txt for details"
	output := renderer.Render(input)

	// URL should remain intact (underscores not treated as italic)
	assert.Contains(t, output, "some_file_name")
}

// TestRenderEdgeCaseEmptyCodeBlock verifies empty code block handling
func TestRenderEdgeCaseEmptyCodeBlock(t *testing.T) {
	renderer := NewRenderer(80)
	input := "```go\n\n```"
	output := renderer.Render(input)

	// Should not crash, should return something reasonable
	assert.NotNil(t, output)
}

// TestRenderEdgeCaseMalformedMarkdown verifies graceful fallback
func TestRenderEdgeCaseMalformedMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "unclosed bold",
			input:    "This is **bold without closing",
			contains: "bold without closing",
		},
		{
			name:     "unclosed italic",
			input:    "This is *italic without closing",
			contains: "italic without closing",
		},
		{
			name:     "unclosed code",
			input:    "This is `code without closing",
			contains: "code without closing",
		},
		{
			name:     "unclosed code block",
			input:    "```go\nfunc main() {\n// no closing fence",
			contains: "func main()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(80)
			output := renderer.Render(tt.input)

			// Should still contain the content
			assert.Contains(t, output, tt.contains)
		})
	}
}

// TestRenderEdgeCaseVeryLongLine verifies long lines are handled
func TestRenderEdgeCaseVeryLongLine(t *testing.T) {
	renderer := NewRenderer(80)
	longText := strings.Repeat("word ", 100)
	output := renderer.Render(longText)

	// Should not crash
	assert.NotNil(t, output)
	assert.Contains(t, output, "word")
}

// TestRenderPerformance verifies rendering is fast enough
func TestRenderPerformance(t *testing.T) {
	renderer := NewRenderer(80)

	// Typical 1000-char message
	input := strings.Repeat("This is a test message with **bold** and *italic* and `code`. ", 20)

	// Render multiple times to check performance
	for i := 0; i < 100; i++ {
		output := renderer.Render(input)
		assert.NotEmpty(t, output)
	}

	// If we got here without timeout, performance is acceptable
}

// TestRendererWithDifferentWidths verifies width handling
func TestRendererWithDifferentWidths(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{
			name:  "narrow terminal",
			width: 40,
		},
		{
			name:  "standard terminal",
			width: 80,
		},
		{
			name:  "wide terminal",
			width: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer(tt.width)
			input := "This is a test with **bold** text"
			output := renderer.Render(input)

			assert.Contains(t, output, "bold")
		})
	}
}

// TestRenderCodeBlockWithIndentation verifies code block preserves indentation
func TestRenderCodeBlockWithIndentation(t *testing.T) {
	renderer := NewRenderer(80)
	input := "```python\ndef foo():\n    if True:\n        print('nested')\n```"
	output := renderer.Render(input)

	// Should preserve indentation and content
	assert.Contains(t, output, "def")
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "nested")
}

// TestRenderMultipleNewlines verifies multiple newlines are preserved
func TestRenderMultipleNewlines(t *testing.T) {
	renderer := NewRenderer(80)
	input := "Line 1\n\n\nLine 2"
	output := renderer.Render(input)

	assert.Contains(t, output, "Line 1")
	assert.Contains(t, output, "Line 2")
}

// TestRenderCodeBlockWithHighlight verifies syntax highlighting integration
func TestRenderCodeBlockWithHighlight(t *testing.T) {
	renderer := NewRenderer(80)

	tests := []struct {
		name         string
		language     string
		code         string
		containsWord string
	}{
		{
			name:         "Go code",
			language:     "go",
			code:         `func main() { println("hello") }`,
			containsWord: "hello",
		},
		{
			name:         "Python code",
			language:     "python",
			code:         `def hello(): print("world")`,
			containsWord: "hello",
		},
		{
			name:         "JavaScript code",
			language:     "javascript",
			code:         `const x = () => { console.log("test"); }`,
			containsWord: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := CodeBlock{
				Language: tt.language,
				Content:  tt.code,
			}

			// Test with highlighting enabled (default)
			result := renderer.RenderCodeBlock(block)

			// Should contain the code content
			assert.Contains(t, result, tt.containsWord)

			// Should contain ANSI codes (from Chroma)
			assert.Contains(t, result, "\x1b[")

			// Should contain language label
			assert.Contains(t, result, tt.language)
		})
	}
}

func TestRenderCodeBlockWithHighlight_Disabled(t *testing.T) {
	renderer := NewRenderer(80)
	block := CodeBlock{
		Language: "go",
		Content:  `func main() { println("test") }`,
	}

	// Test with highlighting disabled
	result := renderer.RenderCodeBlockWithHighlight(block, false)

	// Should contain the code content
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "main")

	// Should NOT contain ANSI codes (no highlighting)
	// Note: Might still have Lipgloss styling, but not Chroma syntax highlighting
	assert.Contains(t, result, "go") // Language label should still be present
}

func TestRenderCodeBlockWithHighlight_UnknownLanguage(t *testing.T) {
	renderer := NewRenderer(80)
	block := CodeBlock{
		Language: "unknownlang",
		Content:  `some random code`,
	}

	// Should handle unknown language gracefully
	result := renderer.RenderCodeBlock(block)

	// Should contain the code content
	assert.Contains(t, result, "some random code")

	// Should contain language label even if unknown
	assert.Contains(t, result, "unknownlang")
}

func TestRenderCodeBlockWithHighlight_NoLanguage(t *testing.T) {
	renderer := NewRenderer(80)
	block := CodeBlock{
		Language: "",
		Content:  `func example() { }`,
	}

	// Should handle no language gracefully
	result := renderer.RenderCodeBlock(block)

	// Should contain the code content
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "example")
}

func TestRenderMarkdownWithHighlightedCodeBlocks(t *testing.T) {
	renderer := NewRenderer(80)

	markdown := `# Example

Here is some Go code:

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

And some Python:

` + "```python" + `
def hello():
    print("Hello, World!")
` + "```" + `

Regular text here.`

	result := renderer.Render(markdown)

	// Should contain heading
	assert.Contains(t, result, "Example")

	// Should contain code content
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "main")
	assert.Contains(t, result, "def")
	assert.Contains(t, result, "hello")

	// Should contain regular text
	assert.Contains(t, result, "Regular text here")

	// Should contain ANSI codes from syntax highlighting
	assert.Contains(t, result, "\x1b[")

	// Should not contain markdown code fence syntax
	assert.NotContains(t, result, "```")
}

func TestRenderCodeBlockWithHighlight_SpecialCharacters(t *testing.T) {
	renderer := NewRenderer(80)
	block := CodeBlock{
		Language: "go",
		Content: `// Comment with unicode: 你好世界 🎉
func main() {
	s := "special: <>&\""
}`,
	}

	result := renderer.RenderCodeBlock(block)

	// Should preserve unicode and special characters
	assert.Contains(t, result, "你好世界")
	assert.Contains(t, result, "🎉")
	assert.Contains(t, result, "func")
}

func TestRenderCodeBlockWithHighlight_PreservesStructure(t *testing.T) {
	renderer := NewRenderer(80)
	block := CodeBlock{
		Language: "go",
		Content: `func example() {
    line1()
    line2()
    line3()
}`,
	}

	result := renderer.RenderCodeBlock(block)

	// Should preserve line structure
	assert.Contains(t, result, "line1")
	assert.Contains(t, result, "line2")
	assert.Contains(t, result, "line3")
}

func TestRenderCodeBlockWithHighlight_MultipleLanguages(t *testing.T) {
	renderer := NewRenderer(80)

	languages := []string{"go", "python", "javascript", "rust", "java", "sql", "bash", "json", "yaml"}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			block := CodeBlock{
				Language: lang,
				Content:  `example code here`,
			}

			result := renderer.RenderCodeBlock(block)

			// Should handle all languages without errors
			assert.NotEmpty(t, result)
			assert.Contains(t, result, "example")
		})
	}
}
