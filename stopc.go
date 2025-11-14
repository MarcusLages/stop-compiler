package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run %s <sourcefile>", os.Args[0])
		return
	}

	// Read from the input file (it will be the last argument)
	srcPath := os.Args[len(os.Args)-1]
	data, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Println("Error reading source file:", err)
		return
	}

	tokens := Lexer(string(data))
	// Peek_lexer(tokens)

	ast := Parser(tokens)
	// Peek_parser(ast)

	sa_analyzer := new_semantic_analyzer()
	sa_analyzer.analyze(ast)
	Peek_semantic(sa_analyzer, ast)
}
