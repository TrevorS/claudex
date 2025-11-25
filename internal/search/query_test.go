// ABOUTME: query_test.go contains comprehensive tests for the query parser,
// including AST construction, operator precedence, and error handling.
package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_SingleTerm(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:     "simple word",
			input:    "authentication",
			expected: &TermNode{Term: "authentication"},
		},
		{
			name:     "quoted string",
			input:    `"hello world"`,
			expected: &TermNode{Term: "hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_FieldValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "simple field:value",
			input: "title:authentication",
			expected: &FieldNode{
				Field: "title",
				Value: "authentication",
				Op:    "=",
			},
		},
		{
			name:  "field with quoted value",
			input: `title:"my chat"`,
			expected: &FieldNode{
				Field: "title",
				Value: "my chat",
				Op:    "=",
			},
		},
		{
			name:  "numeric field with greater than",
			input: "tokens:>5000",
			expected: &FieldNode{
				Field: "tokens",
				Value: "5000",
				Op:    ">",
			},
		},
		{
			name:  "numeric field with less than",
			input: "tokens:<1000",
			expected: &FieldNode{
				Field: "tokens",
				Value: "1000",
				Op:    "<",
			},
		},
		{
			name:  "numeric field with greater or equal",
			input: "tokens:>=5000",
			expected: &FieldNode{
				Field: "tokens",
				Value: "5000",
				Op:    ">=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_ANDOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "explicit AND",
			input: "ai AND ml",
			expected: &AndNode{
				Left:  &TermNode{Term: "ai"},
				Right: &TermNode{Term: "ml"},
			},
		},
		{
			name:  "implicit AND (space)",
			input: "ai ml",
			expected: &AndNode{
				Left:  &TermNode{Term: "ai"},
				Right: &TermNode{Term: "ml"},
			},
		},
		{
			name:  "multiple AND (left associative)",
			input: "ai AND ml AND python",
			expected: &AndNode{
				Left: &AndNode{
					Left:  &TermNode{Term: "ai"},
					Right: &TermNode{Term: "ml"},
				},
				Right: &TermNode{Term: "python"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_OROperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "simple OR",
			input: "python OR rust",
			expected: &OrNode{
				Left:  &TermNode{Term: "python"},
				Right: &TermNode{Term: "rust"},
			},
		},
		{
			name:  "multiple OR (left associative)",
			input: "python OR rust OR go",
			expected: &OrNode{
				Left: &OrNode{
					Left:  &TermNode{Term: "python"},
					Right: &TermNode{Term: "rust"},
				},
				Right: &TermNode{Term: "go"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_NOTOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "simple NOT",
			input: "NOT archived",
			expected: &NotNode{
				Node: &TermNode{Term: "archived"},
			},
		},
		{
			name:  "NOT with field",
			input: "NOT title:old",
			expected: &NotNode{
				Node: &FieldNode{
					Field: "title",
					Value: "old",
					Op:    "=",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_Parentheses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:     "simple grouping",
			input:    "(hello)",
			expected: &TermNode{Term: "hello"},
		},
		{
			name:  "grouping with OR",
			input: "(ai OR ml)",
			expected: &OrNode{
				Left:  &TermNode{Term: "ai"},
				Right: &TermNode{Term: "ml"},
			},
		},
		{
			name:  "grouping changes precedence",
			input: "(ai OR ml) AND python",
			expected: &AndNode{
				Left: &OrNode{
					Left:  &TermNode{Term: "ai"},
					Right: &TermNode{Term: "ml"},
				},
				Right: &TermNode{Term: "python"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_OperatorPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "AND has higher precedence than OR",
			input: "ai OR ml AND python",
			expected: &OrNode{
				Left: &TermNode{Term: "ai"},
				Right: &AndNode{
					Left:  &TermNode{Term: "ml"},
					Right: &TermNode{Term: "python"},
				},
			},
		},
		{
			name:  "NOT has highest precedence",
			input: "ai AND NOT ml",
			expected: &AndNode{
				Left: &TermNode{Term: "ai"},
				Right: &NotNode{
					Node: &TermNode{Term: "ml"},
				},
			},
		},
		{
			name:  "complex precedence",
			input: "ai OR ml AND NOT python",
			expected: &OrNode{
				Left: &TermNode{Term: "ai"},
				Right: &AndNode{
					Left: &TermNode{Term: "ml"},
					Right: &NotNode{
						Node: &TermNode{Term: "python"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Node
	}{
		{
			name:  "mixed operators with parentheses",
			input: "(title:chat OR content:auth) AND model:sonnet",
			expected: &AndNode{
				Left: &OrNode{
					Left: &FieldNode{
						Field: "title",
						Value: "chat",
						Op:    "=",
					},
					Right: &FieldNode{
						Field: "content",
						Value: "auth",
						Op:    "=",
					},
				},
				Right: &FieldNode{
					Field: "model",
					Value: "sonnet",
					Op:    "=",
				},
			},
		},
		{
			name:  "complex query with all features",
			input: `(title:"my chat" OR content:authentication) AND NOT tokens:>5000`,
			expected: &AndNode{
				Left: &OrNode{
					Left: &FieldNode{
						Field: "title",
						Value: "my chat",
						Op:    "=",
					},
					Right: &FieldNode{
						Field: "content",
						Value: "authentication",
						Op:    "=",
					},
				},
				Right: &NotNode{
					Node: &FieldNode{
						Field: "tokens",
						Value: "5000",
						Op:    ">",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			node, err := parser.Parse()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, node)
		})
	}
}

func TestParser_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "unclosed parenthesis",
			input:   "(hello",
			wantErr: true,
		},
		{
			name:    "extra closing parenthesis",
			input:   "hello)",
			wantErr: true,
		},
		{
			name:    "operator without right operand",
			input:   "hello AND",
			wantErr: true,
		},
		{
			name:    "NOT without operand",
			input:   "NOT",
			wantErr: true,
		},
		{
			name:    "field without value",
			input:   "title:",
			wantErr: true,
		},
		{
			name:    "empty parentheses",
			input:   "()",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
