// ABOUTME: Comprehensive tests for syntax highlighting using Chroma v2 library.
// Tests cover language detection, error handling, performance, and integration.
package render

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHighlightCode_BasicHighlighting tests basic syntax highlighting for common languages
func TestHighlightCode_Go(t *testing.T) {
	code := `func main() {
	fmt.Println("Hello, World!")
}`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "main")
	// Should contain ANSI escape codes
	assert.Contains(t, result, "\x1b[")
}

func TestHighlightCode_Python(t *testing.T) {
	code := `def hello():
    print("Hello, World!")
    return True`
	result, err := HighlightCode(code, "python")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "def")
	assert.Contains(t, result, "hello")
	// Should contain ANSI escape codes
	assert.Contains(t, result, "\x1b[")
}

func TestHighlightCode_JavaScript(t *testing.T) {
	code := `function hello() {
    const msg = "Hello, World!";
    console.log(msg);
}`
	result, err := HighlightCode(code, "javascript")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "function")
	assert.Contains(t, result, "const")
	// Should contain ANSI escape codes
	assert.Contains(t, result, "\x1b[")
}

func TestHighlightCode_PreservesStructure(t *testing.T) {
	code := `line 1
line 2
line 3`
	result, err := HighlightCode(code, "text")
	assert.NoError(t, err)
	assert.Contains(t, result, "line 1")
	assert.Contains(t, result, "line 2")
	assert.Contains(t, result, "line 3")
}

// TestHighlightCode_LanguageDetection tests language detection and fallbacks
func TestHighlightCode_ExplicitLanguageHint(t *testing.T) {
	code := `package main`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.Contains(t, result, "package")
}

func TestHighlightCode_CaseInsensitiveLanguage(t *testing.T) {
	code := `SELECT * FROM users;`

	result1, err1 := HighlightCode(code, "sql")
	assert.NoError(t, err1)

	result2, err2 := HighlightCode(code, "SQL")
	assert.NoError(t, err2)

	// Both should work (case-insensitive)
	assert.NotEmpty(t, result1)
	assert.NotEmpty(t, result2)
}

func TestHighlightCode_UnknownLanguage(t *testing.T) {
	code := `some random text`
	result, err := HighlightCode(code, "unknownlang12345")
	assert.NoError(t, err)
	// Should fall back gracefully (return the code, possibly unhighlighted)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "some random text")
}

func TestHighlightCode_EmptyLanguageHint(t *testing.T) {
	code := `func main() { }`
	result, err := HighlightCode(code, "")
	assert.NoError(t, err)
	// Should auto-detect or fall back to plaintext
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "func")
}

// TestHighlightCode_ErrorHandling tests robustness and error scenarios
func TestHighlightCode_EmptyCode(t *testing.T) {
	result, err := HighlightCode("", "go")
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestHighlightCode_InvalidCode(t *testing.T) {
	// Syntactically invalid code should not crash
	code := `func {{{ invalid syntax ]]]`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	// Should return something (even if just the original code)
	assert.NotEmpty(t, result)
}

func TestHighlightCode_SpecialCharacters(t *testing.T) {
	code := `// Comment with unicode: 你好世界 🎉
func main() {
	s := "special chars: <>&\""
}`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.Contains(t, result, "你好世界")
	assert.Contains(t, result, "🎉")
	assert.Contains(t, result, "main")
}

func TestHighlightCode_VeryLargeCode(t *testing.T) {
	// Generate a large code block (5000 lines)
	var builder strings.Builder
	for i := 0; i < 5000; i++ {
		builder.WriteString("func example")
		builder.WriteString(string(rune('0' + (i % 10))))
		builder.WriteString("() { return ")
		builder.WriteString(string(rune('0' + (i % 10))))
		builder.WriteString(" }\n")
	}
	code := builder.String()

	start := time.Now()
	result, err := HighlightCode(code, "go")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Should complete in reasonable time (<1s for 5000 lines)
	assert.Less(t, elapsed, 1*time.Second, "Large code block took too long to highlight")
}

// TestHighlightCode_SpecificLanguages tests various language-specific highlighting
func TestHighlightCode_Rust(t *testing.T) {
	code := `fn main() {
    let x: i32 = 42;
    println!("x = {}", x);
}`
	result, err := HighlightCode(code, "rust")
	assert.NoError(t, err)
	assert.Contains(t, result, "fn")
	assert.Contains(t, result, "let")
}

func TestHighlightCode_TypeScript(t *testing.T) {
	code := `interface User {
    name: string;
    age: number;
}
const user: User = { name: "Alice", age: 30 };`
	result, err := HighlightCode(code, "typescript")
	assert.NoError(t, err)
	assert.Contains(t, result, "interface")
	assert.Contains(t, result, "const")
}

func TestHighlightCode_JSON(t *testing.T) {
	code := `{
    "name": "test",
    "count": 42,
    "enabled": true
}`
	result, err := HighlightCode(code, "json")
	assert.NoError(t, err)
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "42")
	assert.Contains(t, result, "true")
}

func TestHighlightCode_YAML(t *testing.T) {
	code := `name: test
config:
  enabled: true
  count: 42`
	result, err := HighlightCode(code, "yaml")
	assert.NoError(t, err)
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "config")
}

func TestHighlightCode_Bash(t *testing.T) {
	code := `#!/bin/bash
echo "Hello"
for i in {1..5}; do
    echo $i
done`
	result, err := HighlightCode(code, "bash")
	assert.NoError(t, err)
	assert.Contains(t, result, "echo")
	assert.Contains(t, result, "for")
}

func TestHighlightCode_SQL(t *testing.T) {
	code := `SELECT u.name, o.total
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.total > 100;`
	result, err := HighlightCode(code, "sql")
	assert.NoError(t, err)
	assert.Contains(t, result, "SELECT")
	assert.Contains(t, result, "FROM")
	assert.Contains(t, result, "JOIN")
}

// TestHighlightCode_EdgeCases tests various edge cases
func TestHighlightCode_InlineComments(t *testing.T) {
	code := `func main() { // inline comment
    // full line comment
    println("test") // trailing comment
}`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.Contains(t, result, "inline comment")
	assert.Contains(t, result, "full line comment")
	assert.Contains(t, result, "trailing comment")
}

func TestHighlightCode_MultilineStrings(t *testing.T) {
	code := "s := `multi\nline\nstring`"
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.Contains(t, result, "multi")
	assert.Contains(t, result, "line")
	assert.Contains(t, result, "string")
}

func TestHighlightCode_NestedStructures(t *testing.T) {
	code := `type Config struct {
    DB struct {
        Host string
        Port int
    }
}`
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	assert.Contains(t, result, "type")
	assert.Contains(t, result, "struct")
	assert.Contains(t, result, "Host")
}

func TestHighlightCode_WhitespaceOnly(t *testing.T) {
	code := "   \n\t\n   "
	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)
	// Should handle gracefully
	assert.NotEmpty(t, result)
}

// TestHighlightCode_Performance tests highlighting performance
func TestHighlightCode_TypicalCodeBlock(t *testing.T) {
	// Typical code block (~100 lines)
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteString("func example() {\n")
		builder.WriteString("    x := 42\n")
		builder.WriteString("    return x\n")
		builder.WriteString("}\n")
	}
	code := builder.String()

	start := time.Now()
	result, err := HighlightCode(code, "go")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Should be very fast (<50ms for 400 lines)
	assert.Less(t, elapsed, 50*time.Millisecond, "Typical code block took too long: %v", elapsed)
}

func TestHighlightCode_MediumCodeBlock(t *testing.T) {
	// Medium code block (~1000 lines)
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		builder.WriteString("func example() { return 42 }\n")
	}
	code := builder.String()

	start := time.Now()
	result, err := HighlightCode(code, "go")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	// Should complete reasonably fast (<200ms for 1000 lines)
	assert.Less(t, elapsed, 200*time.Millisecond, "Medium code block took too long: %v", elapsed)
}

// TestHighlightQuick tests the convenience function
func TestHighlightQuick(t *testing.T) {
	code := `func main() { }`
	result := HighlightQuick(code, "go")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "func")
}

func TestHighlightQuick_Error(t *testing.T) {
	// Even with errors, should return something
	code := `invalid code`
	result := HighlightQuick(code, "unknownlang")
	assert.NotEmpty(t, result)
}

// TestGetAvailableLanguages tests language listing
func TestGetAvailableLanguages(t *testing.T) {
	languages := GetAvailableLanguages()

	assert.NotEmpty(t, languages)
	assert.Greater(t, len(languages), 100, "Should support many languages")

	// Check for common languages
	containsGo := false
	containsPython := false
	containsJavaScript := false

	for _, lang := range languages {
		lang = strings.ToLower(lang)
		if lang == "go" {
			containsGo = true
		}
		if lang == "python" || lang == "py" {
			containsPython = true
		}
		if lang == "javascript" || lang == "js" {
			containsJavaScript = true
		}
	}

	assert.True(t, containsGo, "Should support Go")
	assert.True(t, containsPython, "Should support Python")
	assert.True(t, containsJavaScript, "Should support JavaScript")
}

// TestHighlightCode_ANSICodes tests that output contains ANSI formatting
func TestHighlightCode_ContainsANSICodes(t *testing.T) {
	code := `func main() { println("test") }`
	result, err := HighlightCode(code, "go")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Should contain ANSI escape sequences for colors
	assert.Contains(t, result, "\x1b[", "Should contain ANSI escape codes")

	// Should contain some color codes (terminal256 format)
	// Example: \x1b[38;5;NNNm for foreground colors
	assert.Regexp(t, `\x1b\[[0-9;]+m`, result, "Should contain ANSI color codes")
}

func TestHighlightCode_DifferentStyles(t *testing.T) {
	// Test that different token types get different colors
	code := `// comment
package main
import "fmt"
func main() {
    x := 42
    s := "string"
}`

	result, err := HighlightCode(code, "go")
	assert.NoError(t, err)

	// Should contain multiple different ANSI codes (different colors for different token types)
	assert.Contains(t, result, "\x1b[")

	// Count unique ANSI codes - should be multiple (comments, keywords, strings, etc.)
	ansiCodes := strings.Count(result, "\x1b[")
	assert.Greater(t, ansiCodes, 5, "Should have multiple different colors for different token types")
}

// Benchmark tests
func BenchmarkHighlightCode_Small(b *testing.B) {
	code := `func main() {
    fmt.Println("Hello, World!")
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HighlightCode(code, "go")
	}
}

func BenchmarkHighlightCode_Medium(b *testing.B) {
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteString("func example() { return 42 }\n")
	}
	code := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HighlightCode(code, "go")
	}
}

func BenchmarkHighlightCode_Large(b *testing.B) {
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		builder.WriteString("func example() { return 42 }\n")
	}
	code := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HighlightCode(code, "go")
	}
}
