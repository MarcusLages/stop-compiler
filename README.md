# STOP Compiler
Toy compiler used as demonstration for a presentation on compilers.
Main purpose of the compiler is to compile the `stop.stp` file into a `.c` file.
The different parts of the compiler are divided into 4 files + a `stop.go` used for general execution and I/O.
Parts not used or skipped for time consuming purposes:
- Tokenization (straight into the lexer)
- IR representation
- Code Optimization

You can find an the target file that we want to compile for the `stop` language as `stop.stp`.

## Lexical Analysis (`stop_tok.go`)
Straight into tokens (no lexemes).
These are the tokens available.
```go
TOKEN_ID     = "ID"
TOKEN_INT    = "INT"
TOKEN_STRING = "STRING"
TOKEN_ASSIGN = "<-"
TOKEN_IF     = "SE"
TOKEN_ELSE   = "SENAO"
TOKEN_PRINT  = "ESCREVA"
TOKEN_OP     = "OP"
TOKEN_EQ     = "="
TOKEN_OPEN   = "/"
TOKEN_CLOSE  = "\\"
TOKEN_EOF    = "EOF"
```

## Synctatic Analysis (`stop_pars.go`)
AST Nodes:
```go
type Node interface{}

type AssignNode struct {
	Name string
	Expr Node
}
type BinOpNode struct {
	Left  Node
	Op    string
	Right Node
}
type NumNode struct {
	Val string
}
type IfNode struct {
	Cond Node
	Then Node
	Else Node
}
type PrintNode struct {
	Expr Node
}
type StringNode struct {
	Val string
}
```

## Semantic Analysis (`stop_sem.go`)
- Basic tree traversal
- Does not produce an annotated AST, only checks for semantics in AST.

## Code Generation (`stop_gen.go`)
- Code generation as string
- Done in one go
- No output to file yet and no execution