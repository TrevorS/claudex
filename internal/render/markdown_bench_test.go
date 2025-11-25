// ABOUTME: markdown_bench_test.go contains benchmarks for the markdown renderer
// to ensure performance targets are met (<10ms for typical 1000-char messages).
package render

import (
	"strings"
	"testing"
)

// BenchmarkRenderTypicalMessage benchmarks a typical 1000-character message.
func BenchmarkRenderTypicalMessage(b *testing.B) {
	renderer := NewRenderer(80)
	input := strings.Repeat("This is a test message with **bold** and *italic* and `code`. ", 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.Render(input)
	}
}

// BenchmarkRenderComplexMarkdown benchmarks complex markdown with code blocks.
func BenchmarkRenderComplexMarkdown(b *testing.B) {
	renderer := NewRenderer(80)
	input := `# Title

This is a paragraph with **bold** and *italic* text.

## Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```" + `

## Features

- Feature 1 with **bold**
- Feature 2 with *italic*
- Feature 3 with ` + "`code`" + `

> This is a quote

Check out [the docs](https://example.com) for more information.
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.Render(input)
	}
}

// BenchmarkRenderPlainText benchmarks plain text rendering.
func BenchmarkRenderPlainText(b *testing.B) {
	renderer := NewRenderer(80)
	input := strings.Repeat("This is plain text without any formatting. ", 25)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.Render(input)
	}
}

// BenchmarkRenderCodeBlocksOnly benchmarks code block rendering.
func BenchmarkRenderCodeBlocksOnly(b *testing.B) {
	renderer := NewRenderer(80)
	input := "```go\nfunc main() {\n    fmt.Println(\"hello world\")\n    for i := 0; i < 10; i++ {\n        process(i)\n    }\n}\n```"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.Render(input)
	}
}

// BenchmarkParseCodeBlocks benchmarks code block extraction.
func BenchmarkParseCodeBlocks(b *testing.B) {
	renderer := NewRenderer(80)
	input := "```go\nfunc main() {}\n```\n\nText\n\n```python\ndef foo():\n    pass\n```"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.ParseCodeBlocks(input)
	}
}

// BenchmarkRenderInlineFormatting benchmarks inline formatting only.
func BenchmarkRenderInlineFormatting(b *testing.B) {
	renderer := NewRenderer(80)
	input := "This is a test with **bold** and *italic* and `code` and [link](url)."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = renderer.RenderInlineFormatting(input)
	}
}
