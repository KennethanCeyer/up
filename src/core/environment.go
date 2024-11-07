package up

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
