package main

import (
	"fmt"
	"strconv"
	"strings"
)

// This Module represents the semantic analyzer.
// Transforms an AST into a semantically correct annotated AST or records
// errors.
// It also creates symbols to check if variables were created already
// and adds types to all literals/vars
// Also, takes out opening/closing pipes (|) from strings while giving types
// to them.

type SymbType string

const (
	STRING_TYPE SymbType = "string"
	INT_TYPE    SymbType = "int"
	ERR_TYPE    SymbType = "error"
)

type SemanticAnalyzer struct {
	symbols map[string]SymbType
	errors  []error
}

func new_semantic_analyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		make(map[string]SymbType),
		[]error{},
	}
}

func semantic_err(err_msg string) error {
	return fmt.Errorf("[Semantic Error] %s", err_msg)
}

func literal_type(val string) SymbType {
	if strings.HasPrefix(val, "|") && strings.HasSuffix(val, "|") {
		return STRING_TYPE
	}
	if _, err := strconv.Atoi(val); err == nil {
		return INT_TYPE
	}
	return ERR_TYPE
}

// Checks for semantic analyzes on the nodes and records all errors
func (sa *SemanticAnalyzer) check_node(node Node) {
	switch n := node.(type) {
	case *IdNode:
		// Must have an existent symbol already
		if _, ok := sa.symbols[n.id]; !ok {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf("Undeclared variable: %s.", n.id),
			))
		}
	case *IfNode:
		sa.check_node(n.cond)
		sa.check_node(n.then)
		if n.else_then != nil {
			sa.check_node(n.else_then)
		}
	case *PrintNode:
		sa.check_node(n.expr)
	case *BlockNode:
		for _, instruction := range n.nodes {
			sa.check_node(instruction)
		}
	case *ErrNode:
		sa.errors = append(sa.errors, n.err_msg)
	case *LitNode:
		return
	default:
		sa.errors = append(sa.errors, semantic_err(
			fmt.Sprintf("Unknown node type: %v", n),
		))
	}
}

// Analyzes the AST collecting errors and creating symbols
// Type checking is done, but no type correction/
func (sa *SemanticAnalyzer) analyze(ast AST) []error {
	for node := range ast.nodes {
		sa.check_node(node)
	}
	return sa.errors
}
