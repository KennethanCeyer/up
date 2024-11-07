package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	compiler "github.com/KennethanCeyer/up/src/compiler"
	core "github.com/KennethanCeyer/up/src/core"
	interpreter "github.com/KennethanCeyer/up/src/interpreter"
)

var options core.Options

func parseOptions() {
	flag.BoolVar(&options.Debug, "debug", true, "")
	flag.BoolVar(&options.Compile, "compile", false, "")
	flag.Parse()
}

func getAttr(obj interface{}, fieldName string) (reflect.Value, error) {
    pointToStruct := reflect.ValueOf(obj)
    curStruct := pointToStruct.Elem()
    if curStruct.Kind() != reflect.Struct {
        return reflect.Value{}, fmt.Errorf("obj is not a struct")
    }
    curField := curStruct.FieldByName(fieldName)
    if !curField.IsValid() {
        return reflect.Value{}, fmt.Errorf("field not found: %s", fieldName)
    }
    return curField, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename.up>")
		for _, arg := range os.Args[1:] {
			fmt.Println(getAttr(options, arg))
		}
		return
	}
	parseOptions()
	cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current directory:", err)
        return
    }
	filename := os.Args[1]
	absolutePath := filepath.Join(cwd, filename)
	fmt.Println(options)
	if options.Compile {
		compiler.Execute(absolutePath, &options)
	} else {
		interpreter.Execute(absolutePath, &options)
	}
}
