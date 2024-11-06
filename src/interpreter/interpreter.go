package interpreter

import (
	"fmt"
	"os"

	interpreter "github.com/KennethanCeyer/up/src/interpreter/internal"
)

type Options struct {
	Debug bool
}

func Execute(filepath string, options *Options) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens, err := interpreter.Lexer(string(data))
	if err != nil {
		fmt.Println("Error in lexical analysis:", err)
		return
	}

	// for logging.
	if options.Debug {
		interpreter.VisualizeTokens(tokens)
	}

	ast, err := interpreter.Parse(tokens)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}

	// for logging.
	if options.Debug {
		interpreter.VisualizeNode(ast)
	}

	env := interpreter.NewEnvironment()
	interpreter.ExecuteNode(ast, env)

	// for logging.
	if options.Debug {
		env.Visualize()
	}
}
