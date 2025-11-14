package main

import (
	"fmt"
	"os"
)

const (
	header string = `
	#include <stdio.h>
	#include <stdlib.h>

	int main() {
	`
	footer string = `
		return 0;
	}
	`
)

func get_format_type(t SymbType) string {
	switch t {
	case STRING_TYPE:
		return `"%s"`
	case INT_TYPE:
		return `"%d"`
	default:
		return ""
	}
}

func gen_node(node Node, sa *SemanticAnalyzer) string {
	switch n := node.(type) {
	case *IdNode:
		return n.id
	case *LitNode:
		t, _ := literal_type(*n)
		if t == STRING_TYPE {
			// Remove the pipes (|) and add double quote (")
			str := n.val[1 : len(n.val)-1]
			return `"` + str + `"`
		}
		return n.val
	case *AssignNode:
		return fmt.Sprintf("%s = %s;\n", n.id.id, gen_node(n.expr, sa))
	case *BinOpNode:
		return fmt.Sprintf("(%s %s %s)", gen_node(n.left, sa), n.op, gen_node(n.right, sa))
	case *PrintNode:
		expr_type, _ := sa.node_type(n.expr)
		format := get_format_type(expr_type)
		return fmt.Sprintf("printf(%s, %s);", format, gen_node(n.expr, sa))
	case *IfNode:
		var else_str string
		if n.else_then != nil {
			else_str = fmt.Sprintf(" else %s", gen_node(n.else_then, sa))
		}
		return fmt.Sprintf("if (%s) %s %s\n",
			gen_node(n.cond, sa),
			gen_node(n.then, sa),
			else_str,
		)
	case *BlockNode:
		output := "{\n"
		for _, inst := range n.nodes {
			output += "\t" + gen_node(inst, sa) + "\n"
		}
		output += "}\n"
		return output
	default:
		return ""
	}
}

func generate(ast AST, sa *SemanticAnalyzer) (string, bool) {
	if len(sa.errors) > 0 {
		fmt.Fprintln(os.Stderr, "Errors found:")
		for _, err := range sa.errors {
			fmt.Fprintln(os.Stderr, " -", err)
		}
		return "", false
	}

	code := header
	for _, n := range ast.nodes {
		code += gen_node(n, sa)
	}
	code += footer
	return code, true
}
