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

type BlockNode struct {
	nodes []Node
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
	case TOKEN_ID:
		id_tok := p.eat(TOKEN_ID)

		// Differentiate ID vs ID Assignment
		if p.peek().tk_type == TOKEN_ASSIGN {
			p.eat(TOKEN_ASSIGN)
			return AssignNode{
				id:   IdNode{id_tok.val},
				expr: p.parse_next(),
			}
		} else {
			return IdNode{id_tok.val}
		}
	case TOKEN_LIT:
		lit := p.eat(TOKEN_LIT)
		return LitNode{lit.val}
	case TOKEN_DO:
		p.eat(TOKEN_DO)
		nodes := []Node{}

		// Keep parsing until there's no more tokens or you find an END token
		for p.peek().tk_type != TOKEN_END || p.more() {
			nodes = append(nodes, p.parse_next())
		}

		// If you get to the end and there's no END token, throw error
		if !p.more() {
			return ErrNode{"Missing `PARE` statement."}
		}

		p.eat(TOKEN_END)
		return BlockNode{nodes}
	case TOKEN_IF:
	case TOKEN_PRINT:
	case TOKEN_ASSIGN:
		return parser_err_node("Isolated assignment. `<-` doesn't an identifier.")
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
