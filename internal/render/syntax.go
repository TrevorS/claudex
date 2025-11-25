// ABOUTME: This file provides syntax highlighting for code blocks using the Chroma library.
// It tokenizes and colorizes source code in various languages for terminal display.

package render

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// HighlightCode applies syntax highlighting to a code snippet.
//
// Parameters:
//   - code: The source code to highlight
//   - language: The programming language (e.g., "go", "python", "javascript")
//
// Returns:
//   - The highlighted code as a string (with ANSI color codes for terminal display)
//   - An error if highlighting fails (gracefully falls back to unhighlighted code)
func HighlightCode(code, language string) (string, error) {
	// Handle empty code
	if len(code) == 0 {
		return "", nil
	}

	// Get lexer for the specified language
	var lexer chroma.Lexer
	if language != "" {
		lexer = lexers.Get(language)
	}

	// If no lexer found (unknown language or empty hint), try to auto-detect
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}

	// If still no lexer, fallback to plaintext
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Ensure lexer is not nil
	if lexer == nil {
		// Ultimate fallback: return code unchanged
		return code, nil
	}

	// Coalesce the lexer to ensure proper tokenization
	lexer = chroma.Coalesce(lexer)

	// Get formatter (terminal256 for broad compatibility)
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		// Fallback to terminal16 if 256 not available
		formatter = formatters.Get("terminal16")
	}
	if formatter == nil {
		// Ultimate fallback: return code unchanged
		return code, nil
	}

	// Get style (monokai for nice colors)
	style := styles.Get("monokai")
	if style == nil {
		// Fallback to vim if monokai not available
		style = styles.Get("vim")
	}
	if style == nil {
		// Fallback to default style
		style = styles.Fallback
	}

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		// Graceful degradation: return unhighlighted code
		return code, nil
	}

	// Format to terminal with colors
	var buf strings.Builder
	if err := formatter.Format(&buf, style, iterator); err != nil {
		// Graceful degradation: return unhighlighted code
		return code, nil
	}

	return buf.String(), nil
}

// HighlightQuick is a convenience function that highlights code without error handling.
// If highlighting fails for any reason, it returns the original code.
func HighlightQuick(code, language string) string {
	result, _ := HighlightCode(code, language)
	if result == "" && code != "" {
		// If highlighting returned empty but we had code, return original
		return code
	}
	return result
}

// GetAvailableLanguages returns a list of all supported language names and aliases.
func GetAvailableLanguages() []string {
	return lexers.Names(true) // true = include aliases
}
