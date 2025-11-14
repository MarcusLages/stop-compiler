# STOP Compiler
Toy compiler used for educational purposes to talk about and present the skeleton of compilers.
Main purpose of the compiler is to compile the `stop.stp` file into a `.c` file.
The different parts of the compiler are divided into 4 files + a `stop.go` used for general execution and I/O.
Parts not used or skipped for time consuming purposes:
- Tokenization (straight into the lexer)
- IR representation
- Code Optimization

You can find an the target file that we want to compile for the `stop` language as `stop.stp`.

## Example Code
```
x <- 2 + 7
y <- |hello\n|
escreva y
se x < 0 va
    escreva(|stop|) 
pare senao va
    escreva(|go|)
pare
```

## Lexical Analysis (`stop_tok.go`)
Straight into tokens (no lexemes).
These are the tokens available.
```go
TOKEN_ID  TokenType = "ID"
TOKEN_LIT TokenType = "LIT"
TOKEN_ASSIGN TokenType = "<-"
TOKEN_IF     TokenType = "SE"
TOKEN_ELSE   TokenType = "SENAO"
TOKEN_PRINT  TokenType = "ESCREVA"
TOKEN_DO     TokenType = "VA"
TOKEN_END    TokenType = "PARE"
TOKEN_EOF    TokenType = "EOF"
TOKEN_OP  TokenType = "OP"
TOKEN_CMP TokenType = "CMP"
TOKEN_ERR TokenType = "ERROR"
```

## Synctatic Analysis (`stop_pars.go`)
- Parses the tokens into AST

## Semantic Analysis (`stop_sem.go`)
- Basic tree traversal
- Does not produce an annotated AST, only checks for semantics in AST.

## Code Generation (`stop_gen.go`)
- Code generation as string
- Done in one go
- No output to file yet and no execution

## Limitations
- Too many to count

## Sources
- [The Super Tiny Compiler](https://github.com/jamiebuilds/the-super-tiny-compiler) by [Jamie Kyle](https://github.com/jamiebuilds)
- [The Super Tiny Compiler's Walkthrough](https://citw.dev/tutorial/create-your-own-compiler) by (@yairhaimo)[https://twitter.com/yairhaimo]
- [Writing a Compiler in Go](https://compilerbook.com)(book) by [Thorsten Ball](https://thorstenball.com)
