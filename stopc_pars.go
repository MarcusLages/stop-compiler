package main

import "fmt"

// Interface used to represent an AST Node
type Node interface{}

// Possible Node types
type IdNode struct {
	id string
}

type LitNode struct {
	val string
}

type AssignNode struct {
	id   IdNode
	expr Node
}

type IfNode struct {
	cond      Node
	then      Node
	else_then Node
}

type PrintNode struct {
	expr Node
}

type BinOpNode struct {
	op    string
	left  Node
	right Node
}

type ErrNode struct {
	err_msg string
}

// Interface used to represent the AST tree itself
type AST struct {
	nodes    []Node
	cur_node Node
}

// Used to keep track of the parser in the peek/eat parser style
// Since the parser has to be very deterministic, you may check the current
// token, but you can only consume it if it knows exactly what it's expecting.
type ParserState struct {
	tokens []Token
	pos    int
}

func parser_err_node(err_msg string) Node {
	return ErrNode{fmt.Sprintf("[Parsing Error] %s", err_msg)}
}

// Checks the current token to be parsed
func (p *ParserState) peek() Token {
	return p.tokens[p.pos]
}

// Consumes the current token if it's of the expected type
// Or else, returns an error. You must know the type of the token to be
// expected before calling eat.
func (p *ParserState) eat(t TokenType) Token {
	cur := p.peek()
	if cur.tk_type == t {
		p.pos++
		return cur
	}
	// Wrong expected token
	return err_token(fmt.Sprintf("Expected: %d, Got: %v", t, cur.tk_type))
}

// True if there's still tokens to parse
func (p *ParserState) more() bool {
	return p.peek().tk_type != TOKEN_EOF
}

func (p *ParserState) parse_next() Node {
	cur := p.peek()

	switch cur.tk_type {
	default:
		return parser_err_node("Unexpected command.")
	}
}

func Parser(tokens []Token) AST {
	parser := ParserState{tokens, 0}
	nodes := []Node{}
	for parser.more() {
		nodes = append(nodes, parser.parse_next())
	}

	return AST{}
}
