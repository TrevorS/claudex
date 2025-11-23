// ABOUTME: JSONL parser for streaming line-by-line JSON parsing of message files
package jsonl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Parser reads and parses JSONL (JSON Lines) format messages from an io.Reader
// Each line should be a valid JSON object representing a Message
type Parser struct {
	scanner *bufio.Scanner
	line    int // Current line number for error reporting
}

// New creates a new JSONL parser from an io.Reader
func New(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		line:    0,
	}
}

// Next reads and returns the next Message from the JSONL stream
// Returns io.EOF when there are no more messages
// Returns an error if JSON parsing fails
func (p *Parser) Next() (*Message, error) {
	for p.scanner.Scan() {
		p.line++
		lineText := p.scanner.Text()

		// Skip empty lines and lines that are only whitespace
		if strings.TrimSpace(lineText) == "" {
			continue
		}

		// Unmarshal the JSON line into a Message
		var msg Message
		if err := json.Unmarshal([]byte(lineText), &msg); err != nil {
			return nil, fmt.Errorf("line %d: invalid JSON: %w", p.line, err)
		}

		return &msg, nil
	}

	// Check for any scanner errors (e.g., read errors)
	if err := p.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	// No more lines; return EOF
	return nil, io.EOF
}
