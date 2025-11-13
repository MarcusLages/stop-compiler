package main

// Interface used to represent the AST tree itself
type AST struct {
	start_node Node
	cur_node   Node
}

// Interface used to represent an AST Node
type Node interface{}

// Possible Node types
// Usually, it would be num, we know that there's only int numbers in this
type IntNode struct {
	val string
}

type StringNode struct {
	val string
}

type AssignNode struct {
	id   string
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

func Parser(tokens []Token) AST {
	return AST{}
}
