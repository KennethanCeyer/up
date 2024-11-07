package up

import (
	"fmt"
	"os"

	core "github.com/KennethanCeyer/up/src/core"
)

func Execute(filepath string, options *core.Options) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens, err := core.Lexer(string(data))
	if err != nil {
		fmt.Println("Error in lexical analysis:", err)
		return
	}

	// for logging.
	if options.Debug {
		core.VisualizeTokens(tokens)
	}

	ast, err := core.Parse(tokens)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}

	// for logging.
	if options.Debug {
		core.VisualizeNode(ast)
	}

	env := core.NewEnvironment()
	core.ExecuteNode(ast, env)

	// for logging.
	if options.Debug {
		env.Visualize()
	}
}
