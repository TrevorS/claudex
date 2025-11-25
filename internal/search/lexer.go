// ABOUTME: lexer.go implements a tokenizer for search query strings, breaking input
// into tokens (words, fields, operators, strings, parentheses) for parsing.
package search

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenWord     TokenType = iota // Regular word
	TokenField                     // Field name (before colon)
	TokenColon                     // :
	TokenOperator                  // AND, OR, NOT, >, <, >=, <=
	TokenString                    // Quoted string
	TokenLParen                    // (
	TokenRParen                    // )
	TokenEOF                       // End of input
)

// Token represents a single lexical token.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer tokenizes search query strings.
type Lexer struct {
	input string
	pos   int
	ch    byte
}

// NewLexer creates a new lexer for the given input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos > len(l.input) {
		return Token{Type: TokenEOF, Value: ""}
	}

	switch l.ch {
	case 0:
		return Token{Type: TokenEOF, Value: ""}
	case '(':
		tok := Token{Type: TokenLParen, Value: string(l.ch)}
		l.readChar()
		return tok
	case ')':
		tok := Token{Type: TokenRParen, Value: string(l.ch)}
		l.readChar()
		return tok
	case ':':
		tok := Token{Type: TokenColon, Value: string(l.ch)}
		l.readChar()
		return tok
	case '"':
		return Token{Type: TokenString, Value: l.readString()}
	case '>':
		return l.readComparisonOperator()
	case '<':
		return l.readComparisonOperator()
	default:
		if isLetterOrDigit(l.ch) || l.ch == '-' || l.ch == '_' {
			word := l.readWord()
			// Check if it's a boolean operator
			upper := strings.ToUpper(word)
			if upper == "AND" || upper == "OR" || upper == "NOT" {
				return Token{Type: TokenOperator, Value: upper}
			}
			// Check if next char is colon (making this a field name)
			if l.ch == ':' {
				return Token{Type: TokenField, Value: word}
			}
			return Token{Type: TokenWord, Value: word}
		}
		// Unknown character, skip it
		l.readChar()
		return l.NextToken()
	}
}

// readChar advances the lexer position and updates current character.
func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
	l.pos++
}

// peekChar looks at the next character without advancing position.
func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readWord reads a word (alphanumeric + hyphens + underscores).
func (l *Lexer) readWord() string {
	start := l.pos - 1
	for isLetterOrDigit(l.ch) || l.ch == '-' || l.ch == '_' {
		l.readChar()
	}
	return l.input[start : l.pos-1]
}

// readString reads a quoted string, handling escape sequences.
func (l *Lexer) readString() string {
	var result strings.Builder
	l.readChar() // skip opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			// Handle escape sequences
			switch l.ch {
			case '"', '\\':
				result.WriteByte(l.ch)
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			default:
				// Unknown escape, just include the character
				result.WriteByte(l.ch)
			}
		} else {
			result.WriteByte(l.ch)
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // skip closing quote
	}

	return result.String()
}

// readComparisonOperator reads >, <, >=, or <=.
func (l *Lexer) readComparisonOperator() Token {
	ch := l.ch
	l.readChar()

	// Check for >= or <=
	if l.ch == '=' {
		op := string(ch) + string(l.ch)
		l.readChar()
		return Token{Type: TokenOperator, Value: op}
	}

	return Token{Type: TokenOperator, Value: string(ch)}
}

// isLetterOrDigit checks if a character is alphanumeric.
func isLetterOrDigit(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch))
}
