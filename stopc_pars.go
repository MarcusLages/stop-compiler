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
	else_then Node // Accepts nil
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

func is_parse_error(n Node) bool {
	if _, ok := n.(ErrNode); ok {
		return true
	} else {
		return false
	}
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

// Parses only expression terms
func (p *ParserState) parse_term() Node {
	tok := p.peek()
	switch tok.tk_type {
	case TOKEN_LIT:
		lit := p.eat(TOKEN_LIT)
		return LitNode{lit.val}
	case TOKEN_ID:
		id := p.eat(TOKEN_ID)
		return IdNode{id.val}
	default:
		return parser_err_node("Not a valid expression term.")
	}
}

func (p *ParserState) parse_val_expr() Node {
	root_node := p.parse_term()

	// Parse binary operations
	for p.peek().tk_type == TOKEN_OP || p.peek().tk_type == TOKEN_CMP {
		// Left hand gets included and now the operation becomes the root now
		// Always add at the right so it all becomes left associative
		op := p.eat(p.peek().tk_type)
		root_node = BinOpNode{op.val, root_node, p.parse_term()}
	}

	return root_node
}

// Do-End block
func (p *ParserState) parse_block() Node {
	tok := p.eat(TOKEN_DO)
	if tok.is_error() {
		return parser_err_node(tok.val)
	}

	nodes := []Node{}

	// Keep parsing until there's no more tokens or you find an END token
	for p.peek().tk_type != TOKEN_END && p.more() {
		nodes = append(nodes, p.parse_next())
	}

	// If you get to the end and there's no END token, throw error
	if !p.more() {
		return ErrNode{fmt.Sprintf("Missing `%s` statement.", TOKEN_END)}
	}

	p.eat(TOKEN_END)
	return BlockNode{nodes}
}

func (p *ParserState) parse_next() Node {
	cur := p.peek()

	switch cur.tk_type {
	case TOKEN_ID:
		node := p.parse_val_expr()
		switch id_node := node.(type) {
		case IdNode:
			if p.peek().tk_type == TOKEN_ASSIGN {
				p.eat(TOKEN_ASSIGN)
				return AssignNode{
					id:   id_node,
					expr: p.parse_val_expr(),
				}
			}
		case ErrNode:
			return id_node
		default:
			return id_node
		}
	case TOKEN_DO:
		return p.parse_block()
	case TOKEN_IF:
		p.eat(TOKEN_IF)

		cond := p.parse_val_expr()
		if is_parse_error(cond) {
			return cond
		}

		then := p.parse_block()
		if is_parse_error(then) {
			return then
		}

		var else_then Node = nil
		if p.peek().tk_type == TOKEN_ELSE {
			p.eat(TOKEN_ELSE)
			else_then = p.parse_block()
		}

		return IfNode{cond, then, else_then}
	case TOKEN_PRINT:
		p.eat(TOKEN_PRINT)
		printable := p.parse_val_expr()
		if is_parse_error(printable) {
			return printable
		}
		return PrintNode{printable}
	case TOKEN_LIT:
		lit_expr := p.parse_val_expr()
		return lit_expr
	case TOKEN_ELSE:
		return parser_err_node(
			fmt.Sprintf("`%s` doesn't have a respective `%s` block.", TOKEN_ELSE, TOKEN_IF),
		)
	case TOKEN_END:
		return parser_err_node(
			fmt.Sprintf("`%s` doesn't have a respective `%s` block.", TOKEN_END, TOKEN_DO),
		)
	case TOKEN_ASSIGN:
		return parser_err_node("Isolated assignment. `<-` doesn't have an identifier.")
	default:
		return parser_err_node(
			fmt.Sprintf("Invalid isolated command `%s`.", cur.val),
		)
	}
	return parser_err_node("Unexpected command.")
}

func Parser(tokens []Token) AST {
	parser := ParserState{tokens, 0}
	nodes := []Node{}
	for parser.more() {
		nodes = append(nodes, parser.parse_next())
	}

	return AST{nodes, nodes[0]}
}

func Peek_parser_tree(node Node, indent string) {
	switch n := node.(type) {
	case IdNode:
		fmt.Println(indent + "Id: " + n.id)
	case LitNode:
		fmt.Println(indent + "Literal: " + n.val)
	case AssignNode:
		fmt.Println(indent + "Assign:")
		fmt.Println(indent + "  LHS:")
		Peek_parser_tree(n.id, indent+"    ")
		fmt.Println(indent + "  RHS:")
		Peek_parser_tree(n.expr, indent+"    ")
	case BinOpNode:
		fmt.Println(indent + "Op: " + n.op)
		fmt.Println(indent + "  Left:")
		Peek_parser_tree(n.left, indent+"    ")
		fmt.Println(indent + "  Right:")
		Peek_parser_tree(n.right, indent+"    ")
	case IfNode:
		fmt.Println(indent + "If:")
		fmt.Println(indent + "  Condition:")
		Peek_parser_tree(n.cond, indent+"    ")
		fmt.Println(indent + "  Then:")
		Peek_parser_tree(n.then, indent+"    ")
		if n.else_then != nil {
			fmt.Println(indent + "  Else:")
			Peek_parser_tree(n.else_then, indent+"    ")
		}
	case PrintNode:
		fmt.Println(indent + "Print:")
		Peek_parser_tree(n.expr, indent+"    ")
	case BlockNode:
		fmt.Println(indent + "Block:")
		for i, child := range n.nodes {
			fmt.Printf(indent+"  [%d]:\n", i)
			Peek_parser_tree(child, indent+"    ")
		}
	case ErrNode:
		fmt.Println(indent + "Err: " + n.err_msg)
	default:
		fmt.Println(indent + "Unknown node type")
	}
}

// Helper for the AST root
func Peek_parser(ast AST) {
	fmt.Println("---- Parser AST ----")
	for i, node := range ast.nodes {
		fmt.Printf("[%d]:\n", i)
		Peek_parser_tree(node, "  ")
	}
	fmt.Println("-------------------")
}
