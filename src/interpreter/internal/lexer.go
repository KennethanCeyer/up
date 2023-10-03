package interpreter

import (
	"strings"
)

type TokenType string

const (
	FUNC       TokenType = "FUNC"
	LPAREN                = "LPAREN"
	RPAREN                = "RPAREN"
	COLON                 = "COLON"
	COMMA                 = "COMMA"
	ARROW                 = "ARROW"
	IDENTIFIER            = "IDENTIFIER"
	INT                   = "INT"
	ADD                   = "ADD"
	ADD_ASSIGN            = "ADD_ASSIGN"
	RETURN                = "RETURN"
	ASSIGN                = "ASSIGN"
	FOR                   = "FOR"
	IN                    = "IN"
	RANGE                 = "RANGE"
	STRING                = "STRING"
	MAIN                  = "MAIN"
	PRINT                 = "PRINT"
	EOF                   = "EOF"
	ENDFOR                = "ENDFOR"
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
				case "print":
					tokens = append(tokens, Token{Type: PRINT, Value: "print", Row: row, Col: col})
				default:
					tokens = append(tokens, Token{Type: IDENTIFIER, Value: identifier, Row: row, Col: col})
				}
			case strings.HasPrefix(input[i:], "->"):
				tokens = append(tokens, Token{Type: ARROW, Value: "->", Row: row, Col: col})
				i += 2
			case input[i] == ':':
				tokens = append(tokens, Token{Type: COLON, Value: ":", Row: row, Col: col})
				i++
			case input[i] == ',':
				tokens = append(tokens, Token{Type: COMMA, Value: ",", Row: row, Col: col})
				i++
			case input[i] == '(':
				tokens = append(tokens, Token{Type: LPAREN, Value: "(", Row: row, Col: col})
				i++
			case input[i] == ')':
				tokens = append(tokens, Token{Type: RPAREN, Value: ")", Row: row, Col: col})
				i++
			case strings.HasPrefix(input[i:], "+="):
				tokens = append(tokens, Token{Type: ADD_ASSIGN, Value: "+=", Row: row, Col: col})
				i += 2
			case input[i] == '+':
				tokens = append(tokens, Token{Type: ADD, Value: "+"})
				i++
			case input[i] == '=':
				tokens = append(tokens, Token{Type: ASSIGN, Value: "=", Row: row, Col: col})
				i++
			case isDigit(input[i]):
				start := i
				for isDigit(input[i]) {
					i++
				}
				tokens = append(tokens, Token{Type: INT, Value: input[start:i], Row: row, Col: col})
			case input[i] == '"':
				start := i
				i++
				for i < len(input) && input[i] != '"' {
					i++
				}
				i++ // Move past the closing "
				tokens = append(tokens, Token{Type: STRING, Value: input[start+1 : i-1], Row: row, Col: col})
		default:
			if input[i] == '\n' {
				row++
				col = 1
			} else {
				col++
			}
			i++
		}
	}

	return tokens, nil
}
