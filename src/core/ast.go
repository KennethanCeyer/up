package up

import (
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

// Expressions
type FloatNode struct {
	Value float64
}

func (n *FloatNode) String() string {
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}

type IntNode struct {
	Value int
}

func (n *IntNode) String() string {
	return strconv.Itoa(n.Value)
}

type StringNode struct {
	Value string
}

func (n *StringNode) String() string {
	return `"` + n.Value + `"`
}

type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) String() string {
	return n.Name
}

type BinOpNode struct {
	Left, Right Node
	Op          string
}

func (n *BinOpNode) String() string {
	return "(" + n.Left.String() + " " + n.Op + " " + n.Right.String() + ")"
}

type FunctionCallNode struct {
	FunctionName string
	Arguments    []Node
}

func (n *FunctionCallNode) String() string {
	args := []string{}
	for _, arg := range n.Arguments {
		args = append(args, arg.String())
	}
	return n.FunctionName + "(" + strings.Join(args, ", ") + ")"
}

type AssignmentNode struct {
	VarName string
	Type    string
	Value   Node
}

func (n *AssignmentNode) String() string {
	return n.VarName + ": " + n.Type + " = " + n.Value.String()
}

// Statements
type ReturnNode struct {
	Value Node
}

func (n *ReturnNode) String() string {
	return "return " + n.Value.String()
}

type ForLoopNode struct {
	Variable     string
	Range        Node
	Body         []Node
}

func (n *ForLoopNode) String() string {
	bodyStrs := []string{}
	for _, stmt := range n.Body {
		bodyStrs = append(bodyStrs, stmt.String())
	}
	return "for " + n.Variable + " in range(" + n.Range.String() + ") {\n\t" + strings.Join(bodyStrs, "\n\t") + "\n}"
}

// Function related
type ParameterNode struct {
	Name string
	Type string
}

func (n *ParameterNode) String() string {
	return n.Name + ": " + n.Type
}

type FuncDeclarationNode struct {
	Name       string
	Parameters []*ParameterNode
	ReturnType string
	Body       []Node
}

func (n *FuncDeclarationNode) String() string {
	paramStrs := []string{}
	for _, param := range n.Parameters {
		paramStrs = append(paramStrs, param.String())
	}
	bodyStrs := []string{}
	for _, stmt := range n.Body {
		bodyStrs = append(bodyStrs, stmt.String())
	}
	return "func " + n.Name + "(" + strings.Join(paramStrs, ", ") + ") -> " + n.ReturnType + " {\n\t" + strings.Join(bodyStrs, "\n\t") + "\n}"
}

type ProgramNode struct {
	Functions []*FuncDeclarationNode
}

func (n *ProgramNode) String() string {
	funcStrs := []string{}
	for _, fn := range n.Functions {
		funcStrs = append(funcStrs, fn.String())
	}
	return strings.Join(funcStrs, "\n\n")
}
