package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run %s <sourcefile>", os.Args[0])
		return
	}

	// Read from the input file (it will be the last argument)
	src_path := os.Args[len(os.Args)-1]
	data, err := os.ReadFile(src_path)
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
	// Peek_semantic(sa_analyzer, ast)

	code, ok := generate(ast, sa_analyzer)
	if ok {
		// fmt.Println(code)
		out_path := strings.TrimSuffix(src_path, ".stp") + ".c"
		err := os.WriteFile(out_path, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing output file: ", err)
			return
		}

		fmt.Println("Compiled to", out_path)
	}
}
