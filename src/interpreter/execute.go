package interpreter

import (
	"fmt"
	"io/ioutil"
	"github.com/KennethanCeyer/up/src/interpreter/internal"
)

type Environment struct {
	store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]interface{})
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val interface{}) {
	e.store[name] = val
}

func executeNode(node interpreter.Node, env *Environment) interface{} {
	switch n := node.(type) {
	case *interpreter.ProgramNode:
		var result interface{}
		for _, function := range n.Functions {
			env.Set(function.Name, function)
			if function.Name == "main" {
				// execute the main function
				result = executeNode(function, env)
			}
		}
		return result
	case *interpreter.FuncDeclarationNode:
		return n
	case *interpreter.FunctionCallNode:
		if function, ok := env.Get(n.FunctionName); ok {
			funcObj := function.(*interpreter.FuncDeclarationNode)
			newEnv := NewEnvironment()
			newEnv.outer = env
			for i, param := range funcObj.Parameters {
				newEnv.Set(param.Name, executeNode(n.Arguments[i], env))
			}
			var result interface{}
			for _, stmt := range funcObj.Body {
				result = executeNode(stmt, newEnv)
			}
			return result
		} else {
			panic("Unknown function: " + n.FunctionName)
		}
	case *interpreter.BinOpNode:
		left := executeNode(n.Left, env)
		right := executeNode(n.Right, env)
		
		// Integer operations
		if lInt, lOk := left.(int); lOk {
			if rInt, rOk := right.(int); rOk {
				switch n.Op {
				case "+":
					return lInt + rInt
				case "-":
					return lInt - rInt
				case "*":
					return lInt * rInt
				case "/":
					if rInt == 0 {
						panic("Division by zero.")
					}
					return lInt / rInt
				case "%":
					return lInt % rInt
				default:
					panic("Unknown operator: " + n.Op)
				}
			}
		}
		
		// String operations
		if lStr, lOk := left.(string); lOk {
			if rStr, rOk := right.(string); rOk {
				switch n.Op {
				case "+":
					return lStr + rStr
				default:
					panic("Invalid operation between strings: " + n.Op)
				}
			}
		}
		
		panic("Invalid operation between different data types.")
	case *interpreter.IntNode:
		return n.Value
	case *interpreter.IdentifierNode:
		if val, ok := env.Get(n.Name); ok {
			return val
		}
		panic("Unknown identifier: " + n.Name)
	default:
		panic("Unknown node type")
	}
}

func Execute(filepath string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens, err := interpreter.Lexer(string(data))
	if err != nil {
		fmt.Println("Error in lexical analysis:", err)
		return
	}
	interpreter.VisualizeTokens(tokens)

	ast, err := interpreter.Parse(tokens)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}

	interpreter.VisualizeNode(ast)

	env := NewEnvironment()
	executeNode(ast, env)
}
