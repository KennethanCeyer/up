package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

type TokenType string

const (
	FUN    TokenType = "FUN"
	IDENTIFIER       = "IDENTIFIER"
	LBRACE           = "LBRACE"
	RBRACE           = "RBRACE"
	LPAREN           = "LPAREN"
	RPAREN           = "RPAREN"
	INT              = "INT"
	ADD              = "ADD"
	COMMA            = "COMMA"
)

type Token struct {
	Type  TokenType
	Value string
}

func Lexer(input string) []Token {
	var tokens []Token

	components := strings.Fields(input)
	for _, component := range components {
		switch component {
		case "fun":
			tokens = append(tokens, Token{Type: FUN, Value: component})
		case "{":
			tokens = append(tokens, Token{Type: LBRACE, Value: component})
		case "}":
			tokens = append(tokens, Token{Type: RBRACE, Value: component})
		case "(":
			tokens = append(tokens, Token{Type: LPAREN, Value: component})
		case ")":
			tokens = append(tokens, Token{Type: RPAREN, Value: component})
		case "+":
			tokens = append(tokens, Token{Type: ADD, Value: component})
		case ",":
			tokens = append(tokens, Token{Type: COMMA, Value: component})
		default:
			if _, err := strconv.Atoi(component); err == nil {
				tokens = append(tokens, Token{Type: INT, Value: component})
			} else {
				tokens = append(tokens, Token{Type: IDENTIFIER, Value: component})
			}
		}
	}

	return tokens
}
