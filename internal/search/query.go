// ABOUTME: query.go implements the query parser that builds an Abstract Syntax Tree (AST)
// from tokenized search queries. Supports operator precedence: NOT > AND > OR.
package search

import (
	"fmt"
)

// Node represents a node in the query Abstract Syntax Tree.
type Node interface {
	// Match evaluates this node against a conversation.
	// This will be implemented by the filter layer.
	node()
}

// TermNode represents a simple search term.
type TermNode struct {
	Term string
}

func (n *TermNode) node() {}

// FieldNode represents a field-specific search (e.g., title:value).
type FieldNode struct {
	Field string
	Value string
	Op    string // "=", ">", "<", ">=", "<="
}

func (n *FieldNode) node() {}

// AndNode represents a boolean AND operation.
type AndNode struct {
	Left  Node
	Right Node
}

func (n *AndNode) node() {}

// OrNode represents a boolean OR operation.
type OrNode struct {
	Left  Node
	Right Node
}

func (n *OrNode) node() {}

// NotNode represents a boolean NOT operation.
type NotNode struct {
	Node Node
}

func (n *NotNode) node() {}

// Parser parses search queries into an AST.
type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
}

// NewParser creates a new parser for the given query string.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.nextToken() // Initialize current
	p.nextToken() // Initialize peek
	return p
}

// Parse parses the query and returns the root AST node.
func (p *Parser) Parse() (Node, error) {
	if p.current.Type == TokenEOF {
		return nil, fmt.Errorf("empty query")
	}
	node, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	// Ensure we've consumed all input
	if p.current.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token after expression: %v", p.current)
	}
	return node, nil
}

// parseOr handles OR expressions (lowest precedence).
func (p *Parser) parseOr() (Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current.Type == TokenOperator && p.current.Value == "OR" {
		p.nextToken() // consume OR
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &OrNode{Left: left, Right: right}
	}

	return left, nil
}

// parseAnd handles AND expressions (medium precedence).
func (p *Parser) parseAnd() (Node, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	// Implicit AND: if we see a term/field/paren without an operator, treat it as AND
	for {
		if p.current.Type == TokenOperator && p.current.Value == "AND" {
			p.nextToken() // consume AND
		} else if p.current.Type == TokenOperator && p.current.Value == "OR" {
			// Stop here, OR has lower precedence
			break
		} else if p.current.Type == TokenEOF || p.current.Type == TokenRParen {
			// End of expression or group
			break
		} else if p.current.Type == TokenWord || p.current.Type == TokenString ||
			p.current.Type == TokenField || p.current.Type == TokenLParen ||
			(p.current.Type == TokenOperator && p.current.Value == "NOT") {
			// Implicit AND - continue parsing
		} else {
			break
		}

		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &AndNode{Left: left, Right: right}
	}

	return left, nil
}

// parseNot handles NOT expressions (high precedence).
func (p *Parser) parseNot() (Node, error) {
	if p.current.Type == TokenOperator && p.current.Value == "NOT" {
		p.nextToken()             // consume NOT
		node, err := p.parseNot() // NOT is right-associative
		if err != nil {
			return nil, err
		}
		return &NotNode{Node: node}, nil
	}

	return p.parsePrimary()
}

// parsePrimary handles primary expressions (terms, fields, parentheses).
func (p *Parser) parsePrimary() (Node, error) {
	switch p.current.Type {
	case TokenLParen:
		p.nextToken() // consume (
		if p.current.Type == TokenRParen {
			return nil, fmt.Errorf("empty parentheses")
		}
		node, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.current.Type != TokenRParen {
			return nil, fmt.Errorf("expected ')', got %v", p.current)
		}
		p.nextToken() // consume )
		return node, nil

	case TokenField:
		return p.parseField()

	case TokenWord:
		term := p.current.Value
		p.nextToken()
		return &TermNode{Term: term}, nil

	case TokenString:
		term := p.current.Value
		p.nextToken()
		return &TermNode{Term: term}, nil

	case TokenEOF:
		return nil, fmt.Errorf("unexpected end of input")

	case TokenRParen:
		return nil, fmt.Errorf("unexpected ')'")

	default:
		return nil, fmt.Errorf("unexpected token: %v", p.current)
	}
}

// parseField handles field:value expressions.
func (p *Parser) parseField() (Node, error) {
	field := p.current.Value
	p.nextToken() // consume field name

	if p.current.Type != TokenColon {
		return nil, fmt.Errorf("expected ':' after field name")
	}
	p.nextToken() // consume :

	// Check for comparison operator
	op := "="
	if p.current.Type == TokenOperator {
		switch p.current.Value {
		case ">", "<", ">=", "<=":
			op = p.current.Value
			p.nextToken()
		}
	}

	// Get the value
	var value string
	switch p.current.Type {
	case TokenWord:
		value = p.current.Value
		p.nextToken()
	case TokenString:
		value = p.current.Value
		p.nextToken()
	default:
		return nil, fmt.Errorf("expected value after field:, got %v", p.current)
	}

	return &FieldNode{Field: field, Value: value, Op: op}, nil
}

// nextToken advances the parser to the next token.
func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}
