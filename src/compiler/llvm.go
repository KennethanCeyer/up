package up

import (
	"fmt"
	"os"

	core "github.com/KennethanCeyer/up/src/core"
	llvm "tinygo.org/x/go-llvm"
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
	if options.Debug {
		core.VisualizeTokens(tokens)
	}

	ast, err := core.Parse(tokens)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}
	if options.Debug {
		core.VisualizeNode(ast)
	}

	ctx := llvm.NewContext()
	defer ctx.Dispose()

	mod := ctx.NewModule("main")
	builder := ctx.NewBuilder()
	defer builder.Dispose()

	varMap := make(map[string]llvm.Value)

	addPrintFunction(mod, ctx)

	generateLLVMIR(ast, varMap, mod, builder, ctx, options.Debug)

	if options.Debug {
		mod.Dump()
	}
}

func addPrintFunction(mod llvm.Module, ctx llvm.Context) {
	int8PtrType := llvm.PointerType(ctx.Int8Type(), 0)
	printType := llvm.FunctionType(ctx.VoidType(), []llvm.Type{int8PtrType}, true)
	llvm.AddFunction(mod, "print", printType)
}

func generateLLVMIR(node core.Node, varMap map[string]llvm.Value, mod llvm.Module, builder llvm.Builder, ctx llvm.Context, debug bool) llvm.Value {
	var result llvm.Value

	switch n := node.(type) {
	case *core.ProgramNode:
		for _, function := range n.Functions {
			result = generateLLVMIR(function, varMap, mod, builder, ctx, debug)
		}
	case *core.FuncDeclarationNode:
		paramTypes := make([]llvm.Type, len(n.Parameters))
		for i := range n.Parameters {
			paramTypes[i] = ctx.IntType(32)
		}
		funcType := llvm.FunctionType(ctx.IntType(32), paramTypes, false)
		function := llvm.AddFunction(mod, n.Name, funcType)
		block := llvm.AddBasicBlock(function, "entry")
		builder.SetInsertPointAtEnd(block)

		for _, param := range n.Parameters {
			varMap[param.Name] = llvm.Undef(ctx.IntType(32))
		}

		for _, stmt := range n.Body {
			result = generateLLVMIR(stmt, varMap, mod, builder, ctx, debug)
		}
		builder.CreateRet(llvm.ConstInt(ctx.IntType(32), 0, false))

	case *core.FunctionCallNode:
		if n.FunctionName == "print" {
			args := make([]llvm.Value, len(n.Arguments))
			for i, arg := range n.Arguments {
				args[i] = generateLLVMIR(arg, varMap, mod, builder, ctx, debug)
				if args[i].IsNil() {
					fmt.Printf("Error: argument %d for function 'print' is invalid\n", i)
					return llvm.Value{}
				}
			}
			printFn := mod.NamedFunction("print")
			if printFn.IsNil() {
				fmt.Println("Error: 'print' function not found")
				return llvm.Value{}
			}
			builder.CreateCall(printFn.Type(), printFn, args, "")
			return llvm.Value{}
		} else {
			args := make([]llvm.Value, len(n.Arguments))
			for i, arg := range n.Arguments {
				args[i] = generateLLVMIR(arg, varMap, mod, builder, ctx, debug)
				if args[i].IsNil() {
					fmt.Printf("Error: argument %d for function '%s' is invalid\n", i, n.FunctionName)
					return llvm.Value{}
				}
			}
			function := mod.NamedFunction(n.FunctionName)
			if function.IsNil() {
				fmt.Printf("Error: function '%s' not found in module\n", n.FunctionName)
				return llvm.Value{}
			}
			result = builder.CreateCall(function.Type(), function, args, "")
		}

	case *core.AssignmentNode:
		val := generateLLVMIR(n.Value, varMap, mod, builder, ctx, debug)
		varMap[n.VarName] = val
		result = val

	case *core.BinOpNode:
		left := generateLLVMIR(n.Left, varMap, mod, builder, ctx, debug)
		right := generateLLVMIR(n.Right, varMap, mod, builder, ctx, debug)
		if left.IsNil() || right.IsNil() {
			fmt.Println("Error: Invalid operands for binary operation")
			return llvm.Value{}
		}
		switch n.Op {
		case "+":
			result = builder.CreateAdd(left, right, "")
		case "-":
			result = builder.CreateSub(left, right, "")
		case "*":
			result = builder.CreateMul(left, right, "")
		case "/":
			result = builder.CreateSDiv(left, right, "")
		}

	case *core.ForLoopNode:
		initialValue := llvm.ConstInt(ctx.IntType(32), 0, false)
		varMap[n.Variable] = initialValue
		loopVarLLVM := varMap[n.Variable]

		loopCond := llvm.AddBasicBlock(builder.GetInsertBlock().Parent(), "loop_cond")
		loopBody := llvm.AddBasicBlock(builder.GetInsertBlock().Parent(), "loop_body")
		loopEnd := llvm.AddBasicBlock(builder.GetInsertBlock().Parent(), "loop_end")

		builder.CreateBr(loopCond)
		builder.SetInsertPointAtEnd(loopCond)

		rangeVal := generateLLVMIR(n.Range, varMap, mod, builder, ctx, debug)
		loopCondition := builder.CreateICmp(llvm.IntSLT, loopVarLLVM, rangeVal, "loop_cond")
		builder.CreateCondBr(loopCondition, loopBody, loopEnd)

		builder.SetInsertPointAtEnd(loopBody)
		for _, bodyNode := range n.Body {
			generateLLVMIR(bodyNode, varMap, mod, builder, ctx, debug)
		}
		builder.CreateBr(loopCond)
		builder.SetInsertPointAtEnd(loopEnd)

	case *core.ReturnNode:
		result = generateLLVMIR(n.Value, varMap, mod, builder, ctx, debug)
		builder.CreateRet(result)

	case *core.FloatNode:
		result = llvm.ConstFloat(ctx.FloatType(), n.Value)

	case *core.IntNode:
		result = llvm.ConstInt(ctx.IntType(32), uint64(n.Value), false)

	case *core.IdentifierNode:
		value, exists := varMap[n.Name]
		if !exists {
			fmt.Printf("Error: identifier '%s' not found in varMap\n", n.Name)
			return llvm.Value{}
		}
		result = value

	default:
		fmt.Println("Unknown node type encountered")
	}

	return result
}
