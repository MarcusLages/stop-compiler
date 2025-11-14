package main

import "fmt"

// This Module represents the semantic analyzer.
// Transforms an AST into a semantically correct annotated AST.
// It also creates symbols to check if variables were created already
// and adds types to all literals/vars
// Also, takes out opening/closing pipes (|) from strings while giving types
// to them.

type SymbType string

const (
	STRING_TYPE SymbType = "string"
	INT_TYPE    SymbType = "int"
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
	return fmt.Errorf("[Semantic Error]: %s", err_msg)
}

func (sa *SemanticAnalyzer) check_node(node Node) {
}

// Analyzes the AST collecting errors and creating symbols
// Type checking is done, but no type correction/
func (sa *SemanticAnalyzer) analyze(ast AST) []error {
	for node := range ast.nodes {
		sa.check_node(node)
	}
	return sa.errors
}
