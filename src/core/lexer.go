package up

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	FUNC        TokenType = "FUNC"
    LBRACE      TokenType = "LBRACE"
    RBRACE      TokenType = "RBRACE"
	LPAREN      TokenType = "LPAREN"
	RPAREN      TokenType = "RPAREN"
	COLON       TokenType = "COLON"
	COMMA       TokenType = "COMMA"
	ARROW       TokenType = "ARROW"
	IDENTIFIER  TokenType = "IDENTIFIER"
	FLOAT       TokenType = "FLOAT"
	INT         TokenType = "INT"
	ADD         TokenType = "ADD"
	SUB         TokenType = "SUB"
	MUL         TokenType = "MUL"
	DIV         TokenType = "DIV"
	ADD_ASSIGN  TokenType = "ADD_ASSIGN"
	SUB_ASSIGN  TokenType = "SUB_ASSIGN"
	MUL_ASSIGN  TokenType = "MUL_ASSIGN"
	DIV_ASSIGN  TokenType = "DIV_ASSIGN"
	RETURN      TokenType = "RETURN"
	ASSIGN      TokenType = "ASSIGN"
	FOR         TokenType = "FOR"
	IN          TokenType = "IN"
	RANGE       TokenType = "RANGE"
	STRING      TokenType = "STRING"
	MAIN        TokenType = "MAIN"
	EOF         TokenType = "EOF"
	ENDFUNC     TokenType = "ENDFUNC"
	ENDFOR      TokenType = "ENDFOR"
)

type Token struct {
	Type  TokenType
	Value string
	Row   int
	Col   int
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func Lexer(input string) ([]Token, error) {
	var tokens []Token

	i := 0
	row, col := 1, 1
	for i < len(input) {
		switch {
		case isAlpha(input[i]) || input[i] == '_':
			start := i
			for isAlpha(input[i]) || isDigit(input[i]) || input[i] == '_' {
				i++
			}
			identifier := input[start:i]
			switch identifier {
				case "func":
					tokens = append(tokens, Token{Type: FUNC, Value: "func", Row: row, Col: col})
				case "return":
					tokens = append(tokens, Token{Type: RETURN, Value: "return", Row: row, Col: col})
				case "for":
					tokens = append(tokens, Token{Type: FOR, Value: "for", Row: row, Col: col})
				case "in":
					tokens = append(tokens, Token{Type: IN, Value: "in", Row: row, Col: col})
				case "range":
					tokens = append(tokens, Token{Type: RANGE, Value: "range", Row: row, Col: col})
				case "main":
					tokens = append(tokens, Token{Type: MAIN, Value: "main", Row: row, Col: col})
				default:
					tokens = append(tokens, Token{Type: IDENTIFIER, Value: identifier, Row: row, Col: col})
				}
			case input[i] == '/' && i+1 < len(input) && input[i+1] == '/':
				i += 2
				for i < len(input) && input[i] != '\n' {
					i++
				}
			case input[i] == ' ' || input[i] == '\t':
				i++
				col++
			case input[i] == '\n':
				i++
				row++
				col = 1
			case strings.HasPrefix(input[i:], "->"):
				tokens = append(tokens, Token{Type: ARROW, Value: "->", Row: row, Col: col})
				i += 2
			case input[i] == ':':
				tokens = append(tokens, Token{Type: COLON, Value: ":", Row: row, Col: col})
				i++
				col++
			case input[i] == ',':
				tokens = append(tokens, Token{Type: COMMA, Value: ",", Row: row, Col: col})
				i++
				col++
			case input[i] == '{':
				tokens = append(tokens, Token{Type: LBRACE, Value: "{", Row: row, Col: col})
				i++
				col++
			case input[i] == '}':
				tokens = append(tokens, Token{Type: RBRACE, Value: "}", Row: row, Col: col})
				i++
				col++
			case input[i] == '(':
				tokens = append(tokens, Token{Type: LPAREN, Value: "(", Row: row, Col: col})
				i++
				col++
			case input[i] == ')':
				tokens = append(tokens, Token{Type: RPAREN, Value: ")", Row: row, Col: col})
				i++
				col++
			case strings.HasPrefix(input[i:], "+="):
				tokens = append(tokens, Token{Type: ADD_ASSIGN, Value: "+=", Row: row, Col: col})
				i += 2
				col += 2
			case strings.HasPrefix(input[i:], "-="):
				tokens = append(tokens, Token{Type: SUB_ASSIGN, Value: "-=", Row: row, Col: col})
				i += 2
				col += 2
			case strings.HasPrefix(input[i:], "*="):
				tokens = append(tokens, Token{Type: MUL_ASSIGN, Value: "*=", Row: row, Col: col})
				i += 2
				col += 2
			case strings.HasPrefix(input[i:], "/="):
				tokens = append(tokens, Token{Type: DIV_ASSIGN, Value: "/=", Row: row, Col: col})
				i += 2
				col += 2
			case input[i] == '+':
				tokens = append(tokens, Token{Type: ADD, Value: "+"})
				i++
				col++
			case input[i] == '-':
				tokens = append(tokens, Token{Type: SUB, Value: "-"})
				i++
				col++
			case input[i] == '*':
				tokens = append(tokens, Token{Type: MUL, Value: "*"})
				i++
				col++
			case input[i] == '/':
				tokens = append(tokens, Token{Type: DIV, Value: "/"})
				i++
				col++
			case input[i] == '=':
				tokens = append(tokens, Token{Type: ASSIGN, Value: "=", Row: row, Col: col})
				i++
				col++
			case isDigit(input[i]):
				start := i
				for isDigit(input[i]) {
					i++
				}
				tokens = append(tokens, Token{Type: INT, Value: input[start:i], Row: row, Col: col})
				col += (i - start)
			case input[i] == '"':
				start := i
				i++
				for i < len(input) && input[i] != '"' {
					i++
				}
				if i < len(input) {
					i++ // Move past the closing "
				}
				tokens = append(tokens, Token{Type: STRING, Value: input[start+1 : i-1], Row: row, Col: col})
				col += (i - start)
		default:
			return nil, fmt.Errorf("unexpected character '%c' at %d:%d", input[i], row, col)
		}
	}

	tokens = append(tokens, Token{Type: EOF, Value: "", Row: row, Col: col})

	return tokens, nil
}
