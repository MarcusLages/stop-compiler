package main

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
	symbols map[string]string
	errors  []string
}

func new_semantic_analyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		make(map[string]string),
		[]string{},
	}
}

// Analyzes the AST collecting errors and creating symbols
// Type checking is done, but no type correction/
func (sa *SemanticAnalyzer) analyze(ast AST) []string {

}
