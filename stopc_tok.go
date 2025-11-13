package main

import (
	"fmt"
	"unicode"
)

// This Module represents the lexical analyzer.
// Contains tokenizer + lexer all together (more of a lexer).

// All possible types of token
type TokenType string

const (
	// Data types
	TOKEN_ID  TokenType = "ID" // Identifiers, such as var names
	TOKEN_INT TokenType = "INT"
	TOKEN_STR TokenType = "STRING"

	// Especial symbols
	TOKEN_ASSIGN TokenType = "<-"
	TOKEN_IF     TokenType = "SE"
	TOKEN_ELSE   TokenType = "SENAO"
	TOKEN_PRINT  TokenType = "ESCREVA"
	TOKEN_DO     TokenType = "VA"   // Opening brackets/do
	TOKEN_END    TokenType = "PARE" // Ending brackets/end
	TOKEN_EOF    TokenType = "EOF"  // Explicit end of file char

	// Operators
	TOKEN_OP TokenType = "OP" // Generic operator
	TOKEN_EQ TokenType = "="
	TOKEN_LT TokenType = "<"
	TOKEN_GT TokenType = ">"
)

type Token struct {
	tk_type TokenType
	val     string
}

func Lexer(input string) []Token {
	tokens := []Token{}
	i := 0

	for i < len(input) {
		ch := rune(input[i]) // Rune is basically byte -> unicode char

		// Skip spaces
		if unicode.IsSpace(ch) {
			i++
			continue

			// Check for digits
		} else if unicode.IsDigit(ch) {
			start := i
			for i < len(input) && unicode.IsDigit(rune(input[i])) {
				i++
			}
			// Not interested on the value, just its classification, so string is ok
			tokens = append(tokens, Token{TOKEN_INT, input[start:i]})
			continue

			// Check for string commands
		} else if unicode.IsLetter(ch) {
			start := i
			// Accepts chars or numbers after the first char
			for i < len(input) &&
				(unicode.IsLetter(rune(input[i])) || unicode.IsLetter(rune(input[i]))) {
				i++
			}

			word := input[start:i]
			switch word {
			case "se":
				tokens = append(tokens, Token{TOKEN_IF, word})
			case "senao":
				tokens = append(tokens, Token{TOKEN_ELSE, word})
			case "va":
				tokens = append(tokens, Token{TOKEN_DO, word})
			case "pare":
				tokens = append(tokens, Token{TOKEN_END, word})
			case "escreva":
				tokens = append(tokens, Token{TOKEN_PRINT, word})
			// If it's a word, that is not a keyword, it's an Identifier
			default:
				tokens = append(tokens, Token{TOKEN_ID, word})
			}
			continue

			// Take care of a string literal, differentiating its start from
			// '<' and "<-"
		} else if ch == '<' && i+1 < len(input) && input[i+1] != '-' {
			i++
			start := i
			for i < len(input) && rune(input[i]) != '>' {
				i++
			}
			tokens = append(tokens, Token{TOKEN_STR, input[start:i]})
			i++ // Skip '>'
			continue
		}

		// Special symbols
		switch ch {
		case '<':
			// Differentiate '<' from "<-"
			if i+1 < len(input) && rune(input[i+1]) == '-' {
				tokens = append(tokens, Token{TOKEN_ASSIGN, "<-"})
				i += 2
			} else {
				tokens = append(tokens, Token{TOKEN_LT, "<"})
				i++
			}
		case '=':
			tokens = append(tokens, Token{TOKEN_EQ, "="})
			i++
		case '-':
			tokens = append(tokens, Token{TOKEN_OP, "-"})
			i++
		case '+':
			tokens = append(tokens, Token{TOKEN_OP, "+"})
			i++
		// Ignore any weird case
		default:
			i++
		}
	}

	// Show that it's the end of the program (parsing sequence)
	tokens = append(tokens, Token{TOKEN_EOF, ""})
	return tokens
}

func Peek_tokens(tokens []Token) {
	fmt.Println("---- Tokenizer/Lexer ----")
	for i, tk := range tokens {
		fmt.Printf("%d. %s(%s)\n", i, tk.tk_type, tk.val)
	}
	fmt.Println("---------------")
}
