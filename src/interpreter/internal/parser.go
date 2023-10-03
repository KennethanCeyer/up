package interpreter

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func isAssignmentOperator(t TokenType) bool {
	return t == ASSIGN || t == ADD_ASSIGN || t == SUB_ASSIGN || t == MUL_ASSIGN || t == DIV_ASSIGN
}

func (p *Parser) isTypeAssignment(offset int) bool {
	return p.lookahead(offset).Type == COLON && p.lookahead(offset+1).Type == IDENTIFIER && isAssignmentOperator(p.lookahead(offset+2).Type)
}

func (p *Parser) consume(t TokenType) Token {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == t {
		p.pos++
		return p.tokens[p.pos-1]
	}
	panic(fmt.Sprintf("Expected token %s at [%d:%d], but got %s", t, p.tokens[p.pos].Row, p.tokens[p.pos].Col, p.tokens[p.pos].Type))
}

func (p *Parser) current() Token {
	return p.tokens[p.pos]
}

func (p *Parser) lookahead(n int) Token {
	if p.pos+n < len(p.tokens) {
		return p.tokens[p.pos+n]
	}
	return Token{}
}

func (p *Parser) parseIdentifier() *IdentifierNode {
	token := p.consume(IDENTIFIER)
	return &IdentifierNode{Name: token.Value}
}

func (p *Parser) parseInt() *IntNode {
	token := p.consume(INT)
	value, _ := strconv.Atoi(token.Value)
	return &IntNode{Value: value}
}

func (p *Parser) parseString() *StringNode {
	token := p.consume(STRING)
	// Assuming the lexer provides string tokens without surrounding quotes
	return &StringNode{Value: token.Value}
}

func (p *Parser) parseParameter() *ParameterNode {
	identifier := p.parseIdentifier()
	p.consume(COLON)
	typeToken := p.consume(IDENTIFIER)
	return &ParameterNode{Name: identifier.Name, Type: typeToken.Value}
}

func (p *Parser) parseBinOp() *BinOpNode {
	left := p.parseExpression()
	opToken := p.current()
	switch opToken.Type {
	case ADD, SUB, MUL, DIV:
		p.pos++
	default:
		panic(fmt.Sprintf("Unexpected token %s for binary operation at [%d:%d]", p.current().Type, p.current().Row, p.current().Col))
	}	
	right := p.parseExpression()
	return &BinOpNode{Left: left, Op: opToken.Value, Right: right}
}

func (p *Parser) parseFunctionCall() *FunctionCallNode {
	funcName := p.parseIdentifier().Name
	p.consume(LPAREN)
	var args []Node
	if p.current().Type != RPAREN {
		args = append(args, p.parseExpression())
		for p.current().Type == COMMA {
			p.consume(COMMA)
			args = append(args, p.parseExpression())
		}
	}
	p.consume(RPAREN)
	return &FunctionCallNode{FunctionName: funcName, Arguments: args}
}

func (p *Parser) parseAssignment() *AssignmentNode {
	varName := p.parseIdentifier().Name

	var varType string
	// Check if type annotation is present
	if p.current().Type == COLON {
		p.consume(COLON) // Consume the colon
		typeToken := p.consume(IDENTIFIER)
		varType = typeToken.Value
	}

	var value Node
	switch p.current().Type {
	case ASSIGN:
		p.consume(ASSIGN)
		value = p.parseExpression()
	case ADD_ASSIGN, SUB_ASSIGN, MUL_ASSIGN, DIV_ASSIGN:
		opToken := p.current()
		p.pos++
		// Check if type is specified with compound assignments, which is not allowed.
		if varType != "" {
			panic(fmt.Sprintf("Cannot specify type with compound assignment at [%d:%d]", opToken.Row, opToken.Col))
		}
		value = &BinOpNode{
			Left:  &IdentifierNode{Name: varName},
			Op:    opToken.Value[:1], // Assuming "+=", "-=", "*=", and "/=" are the token values, we take only the first character
			Right: p.parseExpression(),
		}
	default:
		panic(fmt.Sprintf("Unexpected token %s for assignment at [%d:%d]", p.current().Type, p.current().Row, p.current().Col))
	}
	return &AssignmentNode{VarName: varName, Type: varType, Value: value}
}

func (p *Parser) parseForLoop() *ForLoopNode {
	p.consume(FOR)
	variable := p.parseIdentifier().Name
	p.consume(IN)
	p.consume(RANGE)
	p.consume(LPAREN)
	rng := p.parseExpression()
	p.consume(RPAREN)
	p.consume(LBRACE)

	var body []Node
	for p.current().Type != RBRACE && p.current().Type != EOF {
		body = append(body, p.parseExpression())
	}
	p.consume(RBRACE)
	return &ForLoopNode{Variable: variable, Range: rng, Body: body}
}

func (p *Parser) parseFunction() *FuncDeclarationNode {
	p.consume(FUNC)
	var funcName *IdentifierNode
    if p.current().Type == MAIN {
        funcName = &IdentifierNode{Name: "main"}
        p.consume(MAIN)
    } else {
        funcName = p.parseIdentifier()
    }
	p.consume(LPAREN)
	var parameters []*ParameterNode
	if p.current().Type != RPAREN {
		parameters = append(parameters, p.parseParameter())
		for p.current().Type == COMMA {
			p.consume(COMMA)
			parameters = append(parameters, p.parseParameter())
		}
	}
	p.consume(RPAREN)
	p.consume(ARROW)
	returnType := p.parseIdentifier()
	p.consume(LBRACE)
	var body []Node
	for p.current().Type != RBRACE && p.current().Type != EOF {
		if p.current().Type == RETURN {
			body = append(body, p.parseReturn())
		} else {
			body = append(body, p.parseExpression())
		}
	}
	p.consume(RBRACE)
	return &FuncDeclarationNode{Name: funcName.Name, Parameters: parameters, ReturnType: returnType.Name, Body: body}
}

func (p *Parser) parseReturn() *ReturnNode {
	p.consume(RETURN)
	value := p.parseExpression()
	return &ReturnNode{Value: value}
}

func (p *Parser) parseExpression() Node {
	switch p.current().Type {
	case IDENTIFIER:
		// for logging.
		fmt.Println(p.current(), p.lookahead(1), p.lookahead(2))

		if p.lookahead(1).Type == LPAREN { 
            return p.parseFunctionCall()
        } else if isAssignmentOperator(p.lookahead(1).Type) || p.isTypeAssignment(1) {
            return p.parseAssignment()
        }
        return p.parseIdentifier()
	case INT:
		return p.parseInt()
	case STRING:
		return p.parseString()
	case FOR:
		return p.parseForLoop()
	case ADD, SUB, MUL, DIV:
		return p.parseBinOp()
	case RETURN:
		return p.parseReturn()
	default:
		panic(fmt.Sprintf("Unexpected token %s at [%d:%d]", p.current().Type,  p.current().Row, p.current().Col))
	}
}

func (p *Parser) parseProgram() *ProgramNode {
	var functions []*FuncDeclarationNode
	for p.current().Type != EOF {
		function := p.parseFunction()
		functions = append(functions, function)
	}
	return &ProgramNode{Functions: functions}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func Parse(tokens []Token) (*ProgramNode, error) {
	parser := NewParser(tokens)
	return parser.parseProgram(), nil
}
