package main

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

func main() {
	// Initialize LLVM context
	ctx := llvm.NewContext()
	defer ctx.Dispose()

	// Create a new module
	module := ctx.NewModule("example")
	defer module.Dispose()

	// Create a function type (e.g., a function that takes two integers and returns an integer)
	intType := ctx.Int32Type()
	funcType := llvm.FunctionType(intType, []llvm.Type{intType, intType}, false)

	// Add the function to the module
	function := llvm.AddFunction(module, "add", funcType)

	// Create a new basic block and a builder
	block := ctx.AddBasicBlock(function, "entry")
	builder := ctx.NewBuilder()
	defer builder.Dispose()

	// Set insertion point to the basic block
	builder.SetInsertPointAtEnd(block)

	// Create function arguments
	arg1 := function.Param(0)
	arg2 := function.Param(1)
	result := builder.CreateAdd(arg1, arg2, "result")

	// Return the result
	builder.CreateRet(result)

	// Print out the generated LLVM IR
	module.Dump()
	fmt.Println("LLVM IR generated successfully.")
}
