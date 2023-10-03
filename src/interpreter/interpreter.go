package interpreter

import (
	"fmt"
	"io/ioutil"
	"github.com/KennethanCeyer/up/src/interpreter/internal"
)

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

	// for logging.
	interpreter.VisualizeTokens(tokens)

	ast, err := interpreter.Parse(tokens)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}

	// for logging.
	interpreter.VisualizeNode(ast)

	env := interpreter.NewEnvironment()
	interpreter.ExecuteNode(ast, env)

	// for logging.
	env.Visualize()
}
