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
	NO_TYPE     SymbType = "none"
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

func literal_type(lit LitNode) (SymbType, bool) {
	if strings.HasPrefix(lit.val, "|") && strings.HasSuffix(lit.val, "|") {
		return STRING_TYPE, true
	}
	if _, err := strconv.Atoi(lit.val); err == nil {
		return INT_TYPE, true
	}
	return ERR_TYPE, false
}

func (sa *SemanticAnalyzer) node_type(node Node) (SymbType, bool) {
	switch n := node.(type) {
	case *IdNode:
		t, ok := sa.symbols[n.id]
		if !ok {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf("Undeclared variable: %s.", n.id),
			))
			return ERR_TYPE, false
		}
		return t, true
	case *LitNode:
		return literal_type(*n)
	case *BinOpNode:
		lh_type, ok_l := sa.node_type(n.left)
		rh_type, ok_r := sa.node_type(n.right)
		if !ok_l || !ok_r ||
			(lh_type != rh_type) {
			return ERR_TYPE, false
		}
		return lh_type, true
	case *AssignNode:
		rh_type, ok := sa.node_type(n.expr)
		if !ok {
			return ERR_TYPE, ok
		}
		if sa_id_type, ok := sa.symbols[n.id.id]; ok &&
			sa_id_type != rh_type {
			return ERR_TYPE, ok
		}
		return rh_type, true
	case *BlockNode:
		return NO_TYPE, true
	default:
		return ERR_TYPE, false
	}
}

// Checks for semantics in one node and records its errors
func (sa *SemanticAnalyzer) check_node(node Node) {
	switch n := node.(type) {
	case *IdNode:
		// Must have an existent symbol already
		if _, ok := sa.symbols[n.id]; !ok {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf("Undeclared variable: %s.", n.id),
			))
		}
	case *AssignNode:
		rh_type, ok := sa.node_type(n.expr)
		if !ok {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf(
					"Error in assignment operation '%s' <- %s",
					n.id.id, rh_type,
				),
			))
			return
		}
		if sa_id_type, ok := sa.symbols[n.id.id]; ok &&
			sa_id_type != rh_type {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf(
					"Type mismatch between variable and expression '%s' (%s) <- %s",
					n.id.id, sa_id_type, rh_type,
				),
			))
			return
		}
		sa.symbols[n.id.id] = rh_type
	case *BinOpNode:
		lh_type, ok_l := sa.node_type(n.left)
		rh_type, ok_r := sa.node_type(n.right)

		if !ok_l || !ok_r {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf(
					"Error in binary operation typing '%s': %s %s %s",
					n.op, lh_type, rh_type,
				),
			))
			return
		}

		if lh_type != rh_type {
			sa.errors = append(sa.errors, semantic_err(
				fmt.Sprintf(
					"Type mismatch in binary operation '%s': %s %s %s",
					n.op, lh_type, rh_type,
				),
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
// Type checking is done, but no type correction/casting
func (sa *SemanticAnalyzer) analyze(ast AST) []error {
	for _, node := range ast.nodes {
		sa.check_node(node)
	}
	return sa.errors
}
