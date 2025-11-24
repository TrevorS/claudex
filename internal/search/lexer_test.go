// ABOUTME: lexer_test.go contains comprehensive tests for the search query lexer,
// including tokenization of words, fields, operators, strings, and special characters.
package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_SimpleWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "single word",
			input: "authentication",
			expected: []Token{
				{Type: TokenWord, Value: "authentication"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple words",
			input: "authentication error handling",
			expected: []Token{
				{Type: TokenWord, Value: "authentication"},
				{Type: TokenWord, Value: "error"},
				{Type: TokenWord, Value: "handling"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "words with extra whitespace",
			input: "  hello   world  ",
			expected: []Token{
				{Type: TokenWord, Value: "hello"},
				{Type: TokenWord, Value: "world"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_FieldValuePairs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple field:value",
			input: "title:authentication",
			expected: []Token{
				{Type: TokenField, Value: "title"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "authentication"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "field with quoted value",
			input: `title:"my chat"`,
			expected: []Token{
				{Type: TokenField, Value: "title"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: "my chat"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple field:value pairs",
			input: "model:sonnet tokens:5000",
			expected: []Token{
				{Type: TokenField, Value: "model"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "sonnet"},
				{Type: TokenField, Value: "tokens"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "5000"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "field with comparison operator",
			input: "tokens:>5000",
			expected: []Token{
				{Type: TokenField, Value: "tokens"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenWord, Value: "5000"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "AND operator uppercase",
			input: "ai AND ml",
			expected: []Token{
				{Type: TokenWord, Value: "ai"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenWord, Value: "ml"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "AND operator lowercase",
			input: "ai and ml",
			expected: []Token{
				{Type: TokenWord, Value: "ai"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenWord, Value: "ml"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "OR operator",
			input: "python OR rust",
			expected: []Token{
				{Type: TokenWord, Value: "python"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenWord, Value: "rust"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "NOT operator",
			input: "NOT archived",
			expected: []Token{
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenWord, Value: "archived"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "mixed operators",
			input: "ai AND ml OR NOT archived",
			expected: []Token{
				{Type: TokenWord, Value: "ai"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenWord, Value: "ml"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenWord, Value: "archived"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "comparison operators",
			input: "> < >= <=",
			expected: []Token{
				{Type: TokenOperator, Value: ">"},
				{Type: TokenOperator, Value: "<"},
				{Type: TokenOperator, Value: ">="},
				{Type: TokenOperator, Value: "<="},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_QuotedStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple quoted string",
			input: `"hello world"`,
			expected: []Token{
				{Type: TokenString, Value: "hello world"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "quoted string with escapes",
			input: `"hello \"world\""`,
			expected: []Token{
				{Type: TokenString, Value: `hello "world"`},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple quoted strings",
			input: `"first string" "second string"`,
			expected: []Token{
				{Type: TokenString, Value: "first string"},
				{Type: TokenString, Value: "second string"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "quoted string in field value",
			input: `title:"exact phrase"`,
			expected: []Token{
				{Type: TokenField, Value: "title"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: "exact phrase"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "empty quoted string",
			input: `""`,
			expected: []Token{
				{Type: TokenString, Value: ""},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_Parentheses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple parentheses",
			input: "(hello)",
			expected: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenWord, Value: "hello"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "nested expression",
			input: "(ai OR ml) AND python",
			expected: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenWord, Value: "ai"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenWord, Value: "ml"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenWord, Value: "python"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple groups",
			input: "(ai OR ml) AND (python OR rust)",
			expected: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenWord, Value: "ai"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenWord, Value: "ml"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenLParen, Value: "("},
				{Type: TokenWord, Value: "python"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenWord, Value: "rust"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_ComplexQueries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "complex query with all features",
			input: `(title:"my chat" OR content:authentication) AND model:sonnet AND NOT tokens:>5000`,
			expected: []Token{
				{Type: TokenLParen, Value: "("},
				{Type: TokenField, Value: "title"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenString, Value: "my chat"},
				{Type: TokenOperator, Value: "OR"},
				{Type: TokenField, Value: "content"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "authentication"},
				{Type: TokenRParen, Value: ")"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenField, Value: "model"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "sonnet"},
				{Type: TokenOperator, Value: "AND"},
				{Type: TokenOperator, Value: "NOT"},
				{Type: TokenField, Value: "tokens"},
				{Type: TokenColon, Value: ":"},
				{Type: TokenOperator, Value: ">"},
				{Type: TokenWord, Value: "5000"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexer_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "empty input",
			input: "",
			expected: []Token{
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "only whitespace",
			input: "   \t\n  ",
			expected: []Token{
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "unclosed quoted string",
			input: `"hello`,
			expected: []Token{
				{Type: TokenString, Value: "hello"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "special characters in words",
			input: "hello-world foo_bar",
			expected: []Token{
				{Type: TokenWord, Value: "hello-world"},
				{Type: TokenWord, Value: "foo_bar"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "colon without field name",
			input: ":value",
			expected: []Token{
				{Type: TokenColon, Value: ":"},
				{Type: TokenWord, Value: "value"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token
			for {
				tok := lexer.NextToken()
				tokens = append(tokens, tok)
				if tok.Type == TokenEOF {
					break
				}
			}
			assert.Equal(t, tt.expected, tokens)
		})
	}
}
