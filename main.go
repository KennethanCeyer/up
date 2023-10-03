package main

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/KennethanCeyer/up/src/interpreter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename.up>")
		return
	}
	cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current directory:", err)
        return
    }
	filename := os.Args[1]
	absolutePath := filepath.Join(cwd, filename)
	interpreter.Execute(absolutePath)
}
