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

func (p *Parser) parseFloat() *FloatNode {
	token := p.consume(FLOAT)
	value, _ := strconv.ParseFloat(token.Value, 64)
	return &FloatNode{Value: value}
}

func (p *Parser) parseInt() *IntNode {
	token := p.consume(INT)
	value, _ := strconv.Atoi(token.Value)
	return &IntNode{Value: value}
}

func (p *Parser) parseString() *StringNode {
	token := p.consume(STRING)
	return &StringNode{Value: token.Value}
}

func (p *Parser) parseParameter() *ParameterNode {
	identifier := p.parseIdentifier()
	p.consume(COLON)
	typeToken := p.consume(IDENTIFIER)
	return &ParameterNode{Name: identifier.Name, Type: typeToken.Value}
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

	if p.current().Type == COLON {
		p.consume(COLON)
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
		if varType != "" {
			panic(fmt.Sprintf("Cannot specify type with compound assignment at [%d:%d]", opToken.Row, opToken.Col))
		}
		value = &BinOpNode{
			Left:  &IdentifierNode{Name: varName},
			Op:    opToken.Value[:1],
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
	left := p.parsePrimary()

	for precedence := getPrecedence(p.current().Type); precedence > 0; precedence = getPrecedence(p.current().Type) {
		left = p.parseBinOp(left, precedence)
	}

	return left
}

func (p *Parser) parseTypedDeclaration() *AssignmentNode {
	varName := p.parseIdentifier().Name
	p.consume(COLON)
	typeToken := p.consume(IDENTIFIER)
	varType := typeToken.Value

	var value Node
	if p.current().Type == ASSIGN {
		p.consume(ASSIGN)
		value = p.parseExpression()
	}

	return &AssignmentNode{VarName: varName, Type: varType, Value: value}
}

func (p *Parser) parsePrimary() Node {
	switch p.current().Type {
	case IDENTIFIER:
		if p.lookahead(1).Type == LPAREN {
			return p.parseFunctionCall()
		} else if isAssignmentOperator(p.lookahead(1).Type) || p.lookahead(1).Type == COLON {
			return p.parseAssignment()
		}
		return p.parseIdentifier()
	case INT:
		return p.parseInt()
	case FLOAT:
		return p.parseFloat()
	case STRING:
		return p.parseString()
	case LPAREN:
		p.consume(LPAREN)
		expr := p.parseExpression()
		p.consume(RPAREN)
		return expr
	case FOR:
		return p.parseForLoop()
	default:
		panic(fmt.Sprintf("Unexpected token %s at [%d:%d]", p.current().Type, p.current().Row, p.current().Col))
	}
}

func (p *Parser) parseBinOp(left Node, minPrecedence int) Node {
	for {
		opToken := p.current()
		precedence := getPrecedence(opToken.Type)

		if precedence < minPrecedence {
			break
		}

		p.pos++
		right := p.parsePrimary()

		nextPrecedence := getPrecedence(p.current().Type)
		if precedence < nextPrecedence {
			right = p.parseBinOp(right, precedence+1)
		}

		left = &BinOpNode{Left: left, Op: opToken.Value, Right: right}
	}
	return left
}

func getPrecedence(tokenType TokenType) int {
	switch tokenType {
	case MUL, DIV:
		return 2
	case ADD, SUB:
		return 1
	default:
		return 0
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
