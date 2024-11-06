package interpreter

import (
	"fmt"
	"strings"
)

type Environment struct {
	store map[string]interface{}
	outer *Environment
}

type BuiltinFunction func(args []interface{}) interface{}

func NewEnvironment() *Environment {
	s := make(map[string]interface{})
	env := &Environment{store: s, outer: nil}
	
	// add built-in functions
	env.store["print"] = BuiltinFunction(func(args []interface{}) interface{} {
        for _, arg := range args {
            fmt.Print(arg)
        }
        fmt.Println() // newline after print
        return nil
    })
	
	return env
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

func (e *Environment) Visualize() {
    // Calculate max length of variable name for nice formatting
    maxLen := 0
    for name := range e.store {
        if len(name) > maxLen {
            maxLen = len(name)
        }
    }

    // Header
    fmt.Println("+", strings.Repeat("-", maxLen+2), "+------------------+")
    fmt.Printf("| %-*s | %14s |\n", maxLen, "Variable", "Value")
    fmt.Println("+", strings.Repeat("-", maxLen+2), "+------------------+")

    // Data
    for name, value := range e.store {
        displayValue := formatValue(value)
        fmt.Printf("| %-*s | %14s |\n", maxLen, name, displayValue)
    }

    // Footer
    fmt.Println("+", strings.Repeat("-", maxLen+2), "+------------------+")

    // If the environment has an outer scope, visualize that too
    if e.outer != nil {
        fmt.Println("\nOuter Environment:")
        e.outer.Visualize()
    }
}

func formatValue(value interface{}) string {
    switch v := value.(type) {
    case *FuncDeclarationNode:
        // Return function signature instead of full body
        return fmt.Sprintf("func %s(...)", v.Name)
    case func(...interface{}) interface{}:  // For built-in functions
        return "builtin func"
    default:
        strVal := fmt.Sprintf("%v", value)
        // Truncate long values to fit in table
        if len(strVal) > 12 {
            return strVal[:9] + "..."
        }
        return strVal
    }
}

func ExecuteNode(node Node, env *Environment) interface{} {
	switch n := node.(type) {
	case *ProgramNode:
		var result interface{}
		for _, function := range n.Functions {
			env.Set(function.Name, function)
		}

		if mainFunc, ok := env.Get("main"); ok {
			if mainFuncObj, isFunc := mainFunc.(*FuncDeclarationNode); isFunc {
				newEnv := NewEnvironment()
				newEnv.outer = env
				for _, stmt := range mainFuncObj.Body {
					result = ExecuteNode(stmt, newEnv)
				}
			}
		}
		return result
	case *FuncDeclarationNode:
		return n
	case *FunctionCallNode:
		if function, ok := env.Get(n.FunctionName); ok {
			if funcObj, isUserDefined := function.(*FuncDeclarationNode); isUserDefined {
				newEnv := NewEnvironment()
				newEnv.outer = env

				if len(n.Arguments) != len(funcObj.Parameters) {
					panic(fmt.Sprintf("Expected %d arguments but got %d", len(funcObj.Parameters), len(n.Arguments)))
				}

				for i, param := range funcObj.Parameters {
					newEnv.Set(param.Name, ExecuteNode(n.Arguments[i], env))
				}

				var result interface{}
				for _, stmt := range funcObj.Body {
					result = ExecuteNode(stmt, newEnv)
				}
				return result
			} else if fn, isBuiltIn := function.(BuiltinFunction); isBuiltIn {
				argsVal := make([]interface{}, len(n.Arguments))
				for i, argNode := range n.Arguments {
					argsVal[i] = ExecuteNode(argNode, env)
				}
				return fn(argsVal)
			} else {
				panic("Function " + n.FunctionName + " is neither user-defined nor built-in!")
			}
		} else {
			panic("Function " + n.FunctionName + " not found!")
		}
	case *AssignmentNode:
		val := ExecuteNode(n.Value, env)
		env.Set(n.VarName, val)
		return val
	case *BinOpNode:
		left := ExecuteNode(n.Left, env)
		right := ExecuteNode(n.Right, env)
		
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
		
		// string operations
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
	case *IntNode:
		return n.Value
	case *IdentifierNode:
		if val, ok := env.Get(n.Name); ok {
			return val
		}
		panic("Unknown identifier: " + n.Name)
	case *ForLoopNode:
		rangeValue := ExecuteNode(n.Range, env)
		rangeInt, ok := rangeValue.(int)
		if !ok {
			panic(fmt.Sprintf("Expected integer range, but got: %T", rangeValue))
		}

		var result interface{}
		for i := 0; i < rangeInt; i++ {
			env.Set(n.Variable, i)
			for _, bodyNode := range n.Body {
				result = ExecuteNode(bodyNode, env)
			}
		}
		return result
	case *ReturnNode:
		return ExecuteNode(n.Value, env)
	default:
		panic("Unknown node type")
	}
}
